package proxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	defaultDirectAddr   = "127.0.0.1:3333"
	defaultRedirectAddr = "127.0.0.1:80"

	UserAgent = "go-puzzle/1.0"

	ParamHeader = "header"
	ParamBody   = "body"
	ParamPath   = "path"
	ParamQuery  = "query"

	HeaderCookie        = "Cookie"
	HeaderContentType   = "Content-Type"
	HeaderUserAgent     = "User-Agent"
	HeaderAuthorization = "Authorization"
	HeaderXForwardFor   = "X-Forward-For"
	HeaderXRequestId    = "X-Request-ID"

	BodyJson = "application/json"
)

var (
	ErrInvalidPattern = errors.New("invalid uri pattern")
	ErrUnmatchPath    = errors.New("no match path param")
	ErrProxyReq       = errors.New("empty proxy request")
	ErrMissingQuery   = errors.New("missing query param")
	ErrEmptyResponse  = errors.New("parse empty response")
	ErrContextKey     = errors.New("unmatched http proxy context key")

	coverHeader = []string{HeaderContentType, HeaderAuthorization}
)

type HttpProxy interface {
	Proxy
	ContextKey() string
	GenRequest(c echo.Context, peer HttpPeer) error
	SetupProxy(c echo.Context)
	GoProxy() error
	GetParamSet() ParamSet
	GetRequest() *http.Request
	GetResponse() *http.Response
	OriginBody() ([]byte, error)
	Bind(s interface{}) error
	Transparent(c echo.Context, peer HttpPeer) ([]byte, error)
	TransparentWithJson(c echo.Context, peer HttpPeer, s interface{}) error
}

// Any parameter is ultimately represented as key=value
type param struct {
	Key, Value string
}

// Used to pass parameters
type ParamSet map[string][]param

func (ps ParamSet) SetParam(group string, k, v string) {
	if params, found := ps[group]; found {
		for i := range params {
			if params[i].Key == k {
				params[i].Value = v
				return
			}
		}
	}
	ps[group] = append(ps[group], param{Key: k, Value: v})
}

func (ps ParamSet) AddParam(group string, k, v string) {
	ps[group] = append(ps[group], param{Key: k, Value: v})
}

type httpProxy struct {
	echo.Logger
	proxy
	tls        bool
	contextKey string
	paramSet   ParamSet
	client     *http.Client
	req        *http.Request
	res        *http.Response
}

type HttpPeer struct {
	uri      string
	path     string
	query    string
	Direct   bool
	Path     string
	QuerySet []string
}

func (hp HttpPeer) Uri() string {
	if hp.Direct {
		return hp.uri
	}
	if len(hp.QuerySet) == 0 {
		if hp.query == "" {
			return hp.path
		} else {
			return hp.path + "?" + hp.query
		}
	}
	return hp.path + "?" + strings.Join(hp.QuerySet, "&")
}

type HttpProxyOptioner struct{}

type HttpProxyOption func(p *httpProxy)

func (o HttpProxyOptioner) WithLogger(logger echo.Logger) HttpProxyOption {
	return func(p *httpProxy) {
		p.Logger = logger
	}
}

func (o HttpProxyOptioner) WithTls(enable bool) HttpProxyOption {
	return func(p *httpProxy) {
		p.tls = enable
	}
}

func (o HttpProxyOptioner) WithKey(contextKey string) HttpProxyOption {
	return func(p *httpProxy) {
		p.contextKey = contextKey
	}
}

func (o HttpProxyOptioner) WithAddr(direct, redirect string) HttpProxyOption {
	return func(p *httpProxy) {
		p.proxy.setDirectAddr(direct)
		p.proxy.setRedirectAddr(redirect)
	}
}

func (o HttpProxyOptioner) WithScheme(scheme string) HttpProxyOption {
	return func(p *httpProxy) {
		p.proxy.setScheme(scheme)
	}
}

func NewHttpProxy(opts ...HttpProxyOption) HttpProxy {
	httpProxy := &httpProxy{}
	for _, fn := range opts {
		fn(httpProxy)
	}
	return httpProxy
}

func (p *httpProxy) ContextKey() string {
	return p.contextKey
}

func (p *httpProxy) SetupProxy(c echo.Context) {
	// TODO: setup proxy config
}

func (p *httpProxy) checkUriPath(path string, chars []string) (string, bool) {
	path = strings.ReplaceAll(path, " ", "")
	if !strings.HasPrefix(path, "/") {
		return "", false
	}
	cleanPath := filepath.Clean(path)
	for _, c := range chars {
		if strings.Contains(cleanPath, c) {
			return "", false
		}
	}
	return cleanPath, true
}

func (p *httpProxy) GenRequest(c echo.Context, peer HttpPeer) error {
	uri := c.Request().RequestURI
	rawBody := make(map[string]interface{})
	if c.Request().Header.Get(HeaderContentType) == BodyJson {
		buf := new(bytes.Buffer)
		for _, param := range p.paramSet[ParamBody] {
			rawBody[param.Key] = param.Value
		}
		if err := json.NewEncoder(buf).Encode(rawBody); err != nil {
			return err
		}
		c.Request().Body = io.NopCloser(buf)
	}
	if peer.Direct {
		peer.uri = c.Request().RequestURI
	} else {
		// path:query...
		var ok bool
		peer.path, ok = p.checkUriPath(peer.Path, []string{"!", "#"})
		if !ok {
			panic(fmt.Sprintf("[proxy::http] - invalid http proxy uri path: %s", peer.Path))
		}
		pathSet := strings.SplitN(peer.path, "/", -1)
		for i := range pathSet {
			if pathSet[i] == "" || (pathSet[i][0] != ':' && pathSet[i][0] != '*') {
				continue
			}
			if pathSet[i][0] == '*' && i != len(pathSet)-1 {
				panic(fmt.Sprintf("[proxy::http] - unexpected http proxy uri path param: %s", pathSet[i][1:]))
			}
			// get path value from origin path params
			v := c.Param(pathSet[i][1:])
			match := false
			if v != "" {
				pathSet[i] = v
				if v[0] == '/' {
					// remove the first '/' for join later
					pathSet[i] = v[1:]
				}
				match = true
			} else {
				// get the first path param value from ParamSet
				if params, found := p.paramSet[ParamPath]; found {
					for j := range params {
						if params[j].Key == pathSet[i][1:] {
							pathSet[i] = params[j].Value
							match = true
							break
						}
					}
				}
			}
			if !match {
				return ErrUnmatchPath
			}
		}
		peer.path = strings.Join(pathSet, "/")
		// deal with query param
		if len(peer.QuerySet) == 0 {
			peer.query = c.Request().URL.RawQuery
		} else {
			peer.QuerySet = slices.DeleteFunc(peer.QuerySet, func(s string) bool { return strings.Replace(s, " ", "", -1) == "" })
			for i := range peer.QuerySet {
				match := false
				requested := false
				if peer.QuerySet[i][0] == '!' {
					requested = true
					peer.QuerySet[i] = peer.QuerySet[i][1:]
				}
				// get querySet value from origin querySet params
				v := c.QueryParam(peer.QuerySet[i])
				if v != "" {
					peer.QuerySet[i] += fmt.Sprintf("=%s", v)
					match = true
				} else {
					// get querySet value from raw request body, only Content-Type = application/json
					var value string
					if len(rawBody) > 0 {
						for k, v := range rawBody {
							if k == peer.QuerySet[i] {
								if s, ok := v.(string); ok {
									value = "=" + s
									match = true
									break
								}
							}
						}
					}
					// get querySet from ParamSet, cover query param value from body
					if params, found := p.paramSet[ParamQuery]; found {
						for j := range params {
							if params[j].Key == peer.QuerySet[i] {
								value = "=" + params[j].Value
								match = true
								break
							}
						}
					}
					peer.QuerySet[i] += value
				}
				if !match {
					if requested {
						return ErrMissingQuery
					}
					peer.QuerySet[i] = ""
				}
			}
			peer.QuerySet = slices.DeleteFunc(peer.QuerySet, func(s string) bool { return s == "" })
		}
	}
	// /a/b/c?d=xxx&e=xxx
	uri = peer.Uri()
	req, err := http.NewRequestWithContext(c.Request().Context(), c.Request().Method, p.proxy.Scheme()+p.proxy.DirectAddr()+uri, c.Request().Body)
	if err != nil {
		return err
	}
	p.req = req
	// setup request header
	p.setupHeader()
	return nil
}

// setup all headers
func (p *httpProxy) setupHeader() {
	if p.req == nil {
		panic("[proxy::http] - setup header fail for request is nil")
	}
	for _, param := range p.paramSet[ParamHeader] {
		if param.Key[0] == '!' {
			p.req.Header.Set(param.Key[1:], param.Value)
		} else {
			p.req.Header.Add(param.Key, param.Value)
		}
	}
}

func (p *httpProxy) GoProxy() error {
	if p.req == nil {
		return ErrProxyReq
	}
	var err error
	if !p.proxy.State() {
		p.req, err = http.NewRequestWithContext(p.req.Context(), p.req.Method, p.proxy.Scheme()+p.proxy.RedirectAddr()+p.req.RequestURI, p.req.Body)
		if err != nil {
			return err
		}
	}
	p.res, err = p.client.Do(p.req)
	if err != nil {
		return err
	}
	return nil
}

// status code filter
func (p *httpProxy) statusCodeFilter(left, right int) bool {
	return p.res.StatusCode >= left && p.res.StatusCode < right
}

// get origin body -> []byte
func (p *httpProxy) OriginBody() ([]byte, error) {
	if p.res == nil {
		return []byte(nil), ErrEmptyResponse
	}
	var buf []byte
	_, err := p.res.Body.Read(buf)
	if err != nil {
		return []byte(nil), err
	}
	return buf, nil
}

// get structed body -> json
// s should be pointer
func (p *httpProxy) Bind(s interface{}) error {
	if p.res == nil {
		return ErrEmptyResponse
	}
	buf, err := p.OriginBody()
	if err != nil {
		return err
	}
	if err := json.NewDecoder(bytes.NewBuffer(buf)).Decode(s); err != nil {
		return err
	}
	return nil
}

func (p *httpProxy) GetParamSet() ParamSet {
	return p.paramSet
}

func (p *httpProxy) GetRequest() *http.Request {
	return p.req
}

func (p *httpProxy) GetResponse() *http.Response {
	return p.res
}

// transparent proxy -> []byte
func (p *httpProxy) Transparent(c echo.Context, peer HttpPeer) ([]byte, error) {
	err := p.GenRequest(c, peer)
	if err != nil {
		return []byte(nil), err
	}
	err = p.GoProxy()
	if err != nil {
		return []byte(nil), err
	}
	return p.OriginBody()
}

// transparent proxy -> json
func (p *httpProxy) TransparentWithJson(c echo.Context, peer HttpPeer, s interface{}) error {
	err := p.GenRequest(c, peer)
	if err != nil {
		return err
	}
	err = p.GoProxy()
	if err != nil {
		return err
	}
	return p.Bind(s)
}

func UseHttpProxy(contextKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			hpo := HttpProxyOptioner{}
			hp := NewHttpProxy(
				hpo.WithAddr(defaultDirectAddr, defaultRedirectAddr),
				hpo.WithKey(contextKey),
				hpo.WithScheme(c.Request().URL.Scheme),
				hpo.WithTls(c.IsTLS()),
				hpo.WithLogger(c.Logger()),
			)
			paramSet := hp.GetParamSet()
			// collect header
			for k, vs := range c.Request().Header {
				key := k
				if slices.Contains(coverHeader, k) {
					// TODO: so why not use SetParam(...)?
					key = "!" + k
				}
				for i := range vs {
					// no cover param, just append
					paramSet.AddParam(ParamHeader, key, vs[i])
				}
			}
			// can use AddParam(...)
			c.Echo().IPExtractor = echo.ExtractIPDirect()
			paramSet.SetParam(ParamHeader, HeaderUserAgent, UserAgent)
			paramSet.AddParam(ParamHeader, HeaderXForwardFor, c.RealIP())
			c.Set(hp.ContextKey(), hp)
			return next(c)
		}
	}
}

func GetHttpProxy(c echo.Context, contextKey string) (HttpProxy, error) {
	if hp, ok := c.Get(contextKey).(HttpProxy); ok {
		return hp, nil
	} else {
		return nil, ErrContextKey
	}
}
