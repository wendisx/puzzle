package proxy

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type Proxy interface {
	State() bool
	Scheme() string
	DirectAddr() string
	RedirectAddr() string
}

func (p *proxy) State() bool {
	return p.stateFunc()
}

func (p *proxy) Scheme() string {
	return p.downStream.Scheme
}

func (p *proxy) DirectAddr() string {
	return p.downStream.Addr
}

func (p *proxy) RedirectAddr() string {
	return p.redirect
}

func (p *proxy) setStateFunc(fn func() bool) {
	p.stateFunc = fn
}

func (p *proxy) setDirectAddr(addr string) {
	ip, port, err := net.SplitHostPort(addr)
	if err != nil {
		panic(fmt.Sprintf("[proxy] - bad downstream addr: %s", addr))
	}
	p.downStream.Addr = addr
	p.downStream.IP = ip
	p.downStream.Port = port
}

func (p *proxy) setRedirectAddr(addr string) {
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		panic(fmt.Sprintf("[proxy] - bad redirect addr: %s", addr))
	}
	p.redirect = addr
}

func (p *proxy) setScheme(scheme string) {
	p.downStream.Scheme = scheme
}

// proxy std downstream
type downStream struct {
	Scheme string
	IP     string
	Port   string
	Addr   string
}

// Standard embed agent structure
type proxy struct {
	downStream
	stateFunc func() bool
	redirect  string
}

// new default empty proxy struct
func newDefaultProxy() proxy {
	return proxy{}
}

var (
	globalHttpClient *http.Client
)

const (
	DefaultDerect             = "127.0.0.1:80"
	DefaultTimeout            = 10 * time.Second
	DefaultMaxIdleConn        = 100
	DefaultMaxConnPerHost     = 20
	DefaultMaxIdleConnPerHost = 10
	DefaultMaxResHeaderSize   = 1 << 20

	SchemeHttp  = "http://"
	SchemeHttps = "https://"
)

func init() {
	globalHttpClient = http.DefaultClient
}
