package server

import (
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/errors"
	"github.com/wendisx/puzzle/pkg/router"
)

var (
	_default_skipper = func(c echo.Context) bool {
		return false
	}
	_default_error_handler = errors.EchoErrHandler

	_echo_swagger_path = "/swagger/*"
	_echo_check_path   = "/ping"
)

type (
	// Functional echo server configuration.
	EchoServerOption func(es *EchoServer)
	// web server echo encapsulation
	EchoServer struct {
		webServer[*echo.Echo]
		// More Echo Context...
		gateway router.Route[router.EchoPack]
	}
)

func InitEchoServer() *EchoServer {
	e := echo.New()
	// some default config for echo instance
	e.HTTPErrorHandler = _default_error_handler
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper:     _default_skipper,
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogLatency:  true,
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			// from 127.0.0.1 req GET(200,10)#/healthy
			message := fmt.Sprintf("from %s req %s(%d,%fs)#%s", v.RemoteIP, v.Method, v.Status, v.Latency.Seconds(), v.URI)
			if v.Error == nil {
				clog.Info(message)
			} else {
				clog.Error(message)
			}
			return nil
		},
	}))
	es := &EchoServer{
		webServer: webServer[*echo.Echo]{
			h:    e,
			quit: make(chan os.Signal, 1),
			exit: make(chan struct{}),
			s: &http.Server{
				Addr:    _default_addr,
				Handler: e,
			},
		},
	}
	return es
}

// MountRoute return the default gateway to mount the specified echo instance.
// The Gateway Routing from prefix==""
func (es *EchoServer) MountRoute() router.Route[router.EchoPack] {
	if es.gateway == nil {
		es.gateway = router.NewEchoGateway(es.h)
		clog.Info("mount route for echo server.")
	}
	return es.gateway
}

// MountCheckRoute return the route after mounting check route.
// Nothing to do if flag check is false.
func (es *EchoServer) MountCheckRoute() router.Route[router.EchoPack] {
	gateway := es.MountRoute()
	checkPeer := router.EchoPeer{}
	checkPeer.ToEndpoint(router.Endpoint[echo.HandlerFunc, echo.MiddlewareFunc]{
		Method: http.MethodGet,
		Path:   _echo_check_path,
		Handler: func(c echo.Context) error {
			return c.String(http.StatusOK, "pong")
		},
		PreHandlers: nil,
	})
	gateway.ToPeer(checkPeer)
	clog.Info("mount check peer for echo server.")
	return gateway
}

// MountSwagRoute return the route after mounting swagger route.
// Nothing to do if flag swag is false.
func (es *EchoServer) MountSwagRoute() router.Route[router.EchoPack] {
	gateway := es.MountRoute()
	swagPeer := router.EchoPeer{}
	swagPeer.ToEndpoint(router.Endpoint[echo.HandlerFunc, echo.MiddlewareFunc]{
		Method:      http.MethodGet,
		Path:        _echo_swagger_path,
		Handler:     echoSwagger.WrapHandler,
		PreHandlers: nil,
	})
	gateway.ToPeer(swagPeer)
	clog.Info("mount swag peer for echo server.")
	return gateway
}

func (es *EchoServer) SetupEchoServer(opts ...EchoServerOption) {
	for _, fn := range opts {
		fn(es)
	}
}

func (es *EchoServer) Start() {
	es.gateway.Outbound()
	es.startServer()
}

func (es *EchoServer) Stop() {
	es.quit <- syscall.SIGQUIT
}

func (es *EchoServer) WithCheckRoute(check bool) {
	if check {
		_ = es.MountCheckRoute()
	}
}

func (es *EchoServer) WithSwagRoute(swag bool) {
	if swag {
		_ = es.MountSwagRoute()
	}
}
