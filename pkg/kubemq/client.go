package kubemq

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemqctl/pkg/config"
)

func GetKubeMQClient(ctx context.Context, transport string, cfg *config.Config) (*kubemq.Client, error) {
	switch transport {
	case "grpc":

		return kubemq.NewClient(ctx,
			kubemq.WithAddress(cfg.GetGRPCHostPort()),
			kubemq.WithClientId(uuid.New().String()),
			kubemq.WithTransportType(kubemq.TransportTypeGRPC))

	case "rest":
		return kubemq.NewClient(ctx,
			kubemq.WithUri(cfg.GetRestHttpURI()),
			kubemq.WithClientId(uuid.New().String()),
			kubemq.WithTransportType(kubemq.TransportTypeRest))

	}
	return nil, errors.New("invalid transport type")
}
