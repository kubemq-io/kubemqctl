package kubemq

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	//"fmt"
	"github.com/google/uuid"
	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	//"io/ioutil"
)

var authKey = `eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJDbGllbnRJRCI6ImNsaWVudF9pZCIsIlN0YW5kYXJkQ2xhaW1zIjp7ImV4cCI6MTYwNTQzMjgyMX19.lPZpTxTBWnNa0RUuMqJuJR1CIKQXFUssg5SRv_0QmeM-SsVt8qFH42gdM5-t4Z4NgoQSDrxuFy81eoD1ZemiTilvKlpqlW3YPxVMclIfOdH7UueB_3lM5ZWXFWwXnzkjM7ngm-NMm1RoDKFhA7AXzxHFDegXp4wn51jxfRYX3j-HHwLRl-nHZQITjFYZrdHHhKraeYAj_X9KP8vOyOsTFEN-oiLNSj2Kr4iBnUDHmapaCqH683o2xPb39RFaQC8eZbO1FsCGbvVyo4ToLMzrv3FIfarMG778Jwwm8EaBaajqR9kn96jx7gxWPCpi-ECsz9KqhyVju86jXol-UmAPwA`

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

		key := string(data)
		key = strings.Replace(key, "\r\n", "", -1)

		opts = append(opts, kubemq.WithAuthToken(key))
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
