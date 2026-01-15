package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/palette"
	"github.com/wendisx/puzzle/pkg/router"
)

const (
	_default_quit_delay = 1 * time.Millisecond
	_default_addr       = "127.0.0.1:3333"
)

type (
	// WebServer is the basic abstraction for web server.
	// P represents the route pack.
	WebServer[P any] interface {
		MountRoute() router.Route[P] // return the basic root router
		Start()                      // starting server
		Stop()                       // stopping server
	}
	webServer[H http.Handler] struct {
		h    H
		s    *http.Server
		quit chan os.Signal
		exit chan struct{}
	}
)

func (ws *webServer[H]) startServer() {
	go ws.stopServer(_default_quit_delay)
	clog.Info(fmt.Sprintf("web server listen %s", palette.Put(palette.RGB_SKYBLUE, palette.RGB_DEFAULT).Add(color.Underline).Sprint(ws.s.Addr)))
	if err := ws.s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		clog.Error(fmt.Sprintf("web server start fail for %s", err.Error()))
		close(ws.exit)
	}
	<-ws.exit
	clog.Info("web server exit")
}

func (ws *webServer[H]) stopServer(delay time.Duration) {
	signal.Notify(ws.quit, syscall.SIGTERM, os.Interrupt, syscall.SIGQUIT)
	<-ws.quit
	ctx, cancle := context.WithTimeout(context.Background(), delay)
	defer cancle()
	if err := ws.s.Shutdown(ctx); err != nil {
		clog.Error(fmt.Sprintf("web server shutdown fail after %s", palette.Red(delay)))
		os.Exit(1)
	}
	clog.Info(fmt.Sprintf("web server shutdown success after %s", palette.Green(delay)))
	close(ws.exit)
}
