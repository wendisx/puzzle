package router

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/palette"
)

const (
	_default_echo_gateway = ""
)

var (
	_default_echopack *EchoPack
)

type (
	// EchoPack with Pack
	EchoPack struct {
		P Pack
		G *echo.Group
	}
	// EchoRoute Receive, process, forward, route EchoPack.
	EchoRoute struct {
		ep     EchoPack
		path   string            // for self path
		qroute []Route[EchoPack] // next route queue
		qpeer  []Peer[EchoPack]  // next peer queue
	}
	// EchoPeer hold all endpoints and call Parse to mount the whole endpoints.
	EchoPeer struct {
		qendpoint []Endpoint[echo.HandlerFunc, echo.MiddlewareFunc]
	}
	// EchoGateway always the first route starting routing.
	EchoGateway struct {
		EchoRoute // start from prefix = ""
	}
)

func NewEchoPack(p Pack, g *echo.Group) EchoPack {
	return EchoPack{
		P: p,
		G: g,
	}
}

func NewEchoRoute(path string) *EchoRoute {
	if path == "" {
		clog.Warn("<pkg.router.echo> Empty path route exists.")
	}
	return &EchoRoute{
		path: path,
	}
}

// DefaultEchoPack return the single instance of EchoPack
func DefaultEchoPack(e *echo.Echo) EchoPack {
	if _default_echopack == nil {
		_default_echopack = &EchoPack{
			P: Pack{
				Prefix: _default_echo_gateway,
			},
			G: e.Group(_default_echo_gateway),
		}
	}
	return *_default_echopack
}

// NewEchoGateway return a new EchoRoot to start all
func NewEchoGateway(e *echo.Echo) Route[EchoPack] {
	er := &EchoGateway{
		*NewEchoRoute(_default_echo_gateway),
	}
	er.Inbound(DefaultEchoPack(e))
	return er
}

func (r *EchoRoute) Active() bool {
	return true
}

func (r *EchoRoute) Path() string {
	return r.path
}

// Handle should be overridden and called in Inbound.
func (r *EchoRoute) Handle(p EchoPack) {
	clog.Warn("<pkg.router.echo> Call default Handle function, this method should be explicitly overridden.")
}

// Inbound update EchoPack context and wait to send them out.
func (r *EchoRoute) Inbound(p EchoPack) {
	if !r.Active() {
		return
	}
	// new outbound EchoPack to next route or just peer
	prefix := p.P.Prefix + r.path
	clog.Info("<pkg.router.echo> router cover " + palette.SkyBlue(prefix))
	r.ep = EchoPack{
		P: Pack{
			Prefix: prefix,
		},
		G: p.G.Group(r.path),
	}
}

func (r *EchoRoute) ToRoute(rr Route[EchoPack]) {
	r.qroute = append(r.qroute, rr)
}

func (r *EchoRoute) ToPeer(p Peer[EchoPack]) {
	r.qpeer = append(r.qpeer, p)
}

// Outbound send EchoPack to subnode.
func (r *EchoRoute) Outbound() {
	if !r.Active() {
		return
	}
	if r.ep.G == nil {
		clog.Panic("<pkg.router.echo> Outbound routing link terminate.")
	}
	for i := range r.qroute {
		// if override Active() to false, then just skip all.
		if !strings.Contains(r.ep.P.Prefix, r.qroute[i].Path()) && r.qroute[i].Active() {
			r.qroute[i].Inbound(r.ep)
			r.qroute[i].Outbound()
		}
	}
	for i := range r.qpeer {
		r.qpeer[i].Parse(r.ep)
	}
}

func (p *EchoPeer) ToEndpoint(ep Endpoint[echo.HandlerFunc, echo.MiddlewareFunc]) {
	if p.qendpoint == nil {
		p.qendpoint = make([]Endpoint[echo.HandlerFunc, echo.MiddlewareFunc], 0)
	}
	p.qendpoint = append(p.qendpoint, ep)
}

// Parse parse all endpoints to the instance of global echo
func (p EchoPeer) Parse(pp EchoPack) {
	if len(p.qendpoint) == 0 {
		clog.Warn("Call default Mount function, this method should be explicitly overridden and is eventually called.")
	}
	for _, ep := range p.qendpoint {
		if ep.PreHandlers != nil {
			pp.G.Add(ep.Method, ep.Path, ep.Handler, ep.PreHandlers...)
			clog.Info(fmt.Sprintf("endpoint %s#%s with %s pre handlers", palette.SkyBlue(ep.Method), palette.SkyBlue(pp.P.Prefix+ep.Path), palette.Green(len(ep.PreHandlers))))
		} else {
			pp.G.Add(ep.Method, ep.Path, ep.Handler)
			clog.Info(fmt.Sprintf("endpoint %s#%s with no pre handlers", palette.SkyBlue(ep.Method), palette.SkyBlue(pp.P.Prefix+ep.Path)))
		}
	}
}
