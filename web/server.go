package web

import (
	"context"
	"embed"
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/web/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

var (
	StaticAssetsWeb     embed.FS
	StaticAssetsBuilder embed.FS
)

type Server struct {
	echoWebServer *echo.Echo
	context       context.Context
	cancelFunc    context.CancelFunc
	apiService    *api.Service
	cfg           *config.Config
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init(ctx context.Context, cfg *config.Config) error {
	s.cfg = cfg
	s.context, s.cancelFunc = context.WithCancel(ctx)
	s.apiService = api.NewApiService()
	if err := s.apiService.Init(s.context); err != nil {
		return err
	}
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	webFs := echo.MustSubFS(StaticAssetsWeb, "assets-web")
	builderFs := echo.MustSubFS(StaticAssetsBuilder, "assets-builder")
	e.StaticFS("/", webFs)
	e.StaticFS("/builder", builderFs)
	e.POST("/api/build-request", s.handleProcessBuild)
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
func (s *Server) handleProcessBuild(c echo.Context) error {
	return s.apiService.BuildRequest(c)
}
func (s *Server) Stop() {
	_ = s.echoWebServer.Shutdown(context.Background())
	s.cancelFunc()
}
