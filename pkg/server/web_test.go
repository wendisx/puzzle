package server

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
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
