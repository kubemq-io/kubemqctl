package api

import (
	"context"
	"github.com/kubemq-io/kubemqctl/web/api/builder"
	"github.com/labstack/echo/v4"
)

type Service struct {
	builderService *builder.Service
}

func NewApiService() *Service {
	return &Service{}
}

func (s *Service) Init(ctx context.Context) error {
	s.builderService = builder.NewBuilderService(ctx)
	return nil
}

func (s *Service) BuildRequest(c echo.Context) error {
	resp := NewResponse(c)
	req := &builder.BuilderRequest{}
	err := c.Bind(req)
	if err != nil {
		return c.String(400, err.Error())
	}
	builderResponse, err := s.builderService.ProcessRequest(req)
	if err != nil {
		return resp.SetHttpCode(400).SetError(err).Send()
	}
	return resp.SetResponseBody(builderResponse).Send()
}
