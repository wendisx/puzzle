package server

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	_ "github.com/wendisx/puzzle/docs"
	"github.com/wendisx/puzzle/pkg/router"
)

type (
	ServerPeer struct {
		router.EchoPeer
	}
)

func (sp *ServerPeer) Parse(p router.EchoPack) {
	sp.ToEndpoint(router.Endpoint[echo.HandlerFunc, echo.MiddlewareFunc]{
		Method: http.MethodGet,
		Path:   "/ping",
		Handler: func(c echo.Context) error {
			res := NewEchoResponder(c)
			return res.ResponseHttp(HttpEvent{
				Type:     0x02,
				HttpCode: http.StatusOK,
				Data:     "pong",
			})
		},
		PreHandlers: []echo.MiddlewareFunc{},
	})
	sp.EchoPeer.Parse(p)
}

// test basic echo server [passed]
func Test_echo(t *testing.T) {
	s := InitEchoServer()
	rr := s.MountRoute()
	sp := &ServerPeer{
		router.EchoPeer{},
	}
	rr.ToPeer(sp)
	rr.Outbound()
	s.Start()
}

// test swagger server [passed]
func Test_swagger(t *testing.T) {
	// test echo swagger []
	s := InitEchoServer()
	rr := s.MountSwagRoute()
	rr.Outbound()
	s.Start()
}

// test check server [passed]
func Test_check(t *testing.T) {
	// test echo swagger [passed]
	s := InitEchoServer()
	rr := s.MountCheckRoute()
	rr.Outbound()
	s.Start()
}

// test many routes []
func Test_many_route(t *testing.T) {
	// test echo swagger [passed]
	s := InitEchoServer()
	_ = s.MountCheckRoute()
	rr := s.MountSwagRoute()
	rr.Outbound()
	s.Start()
}
