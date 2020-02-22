package kubemq

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"io/ioutil"
)

func GetKubemqClient(ctx context.Context, transport string, cfg *config.Config) (*kubemq.Client, error) {
	clientId := uuid.New().String()
	if cfg.ClientId != "" {
		clientId = cfg.ClientId
	}
	var opts []kubemq.Option
	opts = append(opts, kubemq.WithClientId(clientId))
	if cfg.AuthTokenFile != "" {
		data, err := ioutil.ReadFile(cfg.AuthTokenFile)
		if err != nil {
			return nil, fmt.Errorf("error loading Authentication token file: %s", err.Error())
		}
		opts = append(opts, kubemq.WithAuthToken(string(data)))
	}
	if cfg.CertFile != "" {
		opts = append(opts, kubemq.WithCredentials(cfg.CertFile, ""))
	}
	switch transport {
	case "grpc":
		opts = append(opts, kubemq.WithAddress(cfg.GetGRPCHostPort()),
			kubemq.WithTransportType(kubemq.TransportTypeGRPC))

		return kubemq.NewClient(ctx, opts...)

	case "rest":
		opts = append(opts, kubemq.WithUri(cfg.GetRestHttpURI()),
			kubemq.WithTransportType(kubemq.TransportTypeRest))
		return kubemq.NewClient(ctx, opts...)
	}
	return nil, errors.New("invalid transport type")
}
