package web

import (
	"context"
	"embed"
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

var (
	StaticAssets embed.FS
)

type Server struct {
	echoWebServer *echo.Echo
	context       context.Context
	cancelFunc    context.CancelFunc
	cfg           *config.Config
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init(ctx context.Context, cfg *config.Config) error {
	s.cfg = cfg
	s.context, s.cancelFunc = context.WithCancel(ctx)
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	fs := echo.MustSubFS(StaticAssets, "assets")
	e.StaticFS("/", fs)
	s.echoWebServer = e
	return nil

}
func (s *Server) Start() error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.echoWebServer.Start(fmt.Sprintf("0.0.0.0:%d", s.cfg.WebPort))
	}()

	select {
	case err := <-errCh:
		return err
	case <-s.context.Done():
		return nil
	case <-time.After(1 * time.Second):
		return nil
	}

}

func (s *Server) Stop() {
	_ = s.echoWebServer.Shutdown(context.Background())
	s.cancelFunc()
}
