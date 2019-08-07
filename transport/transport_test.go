package transport

import (
	"context"
	"testing"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/kubemq-io/kubetools/transport/grpc"
	"github.com/kubemq-io/kubetools/transport/option"
	"github.com/kubemq-io/kubetools/transport/rest"

	"github.com/google/uuid"

	"github.com/stretchr/testify/require"
)

var (
	grpcHost  = "localhost"
	grpcPort  = 50000
	restHost  = "localhost"
	restPort  = 9090
	isSecured = false
	certFile  = ""
	timeout   = time.Duration(5 * time.Second)
)

func getGrpcClient() *grpc.Client {
	c, _ := grpc.New(context.Background(), option.NewOptions(option.ConnectionTypeGrpc, grpcHost, grpcPort, isSecured, certFile))
	return c
}
func getRestClient() *rest.Client {
	c, _ := rest.New(context.Background(), option.NewOptions(option.ConnectionTypeRest, restHost, restPort, isSecured, ""))
	return c
}

func TestClient_Event(t *testing.T) {

	tests := []struct {
		name           string
		clientProducer Transport
		clientConsumer Transport
	}{
		{
			name:           "send_event_grpc",
			clientProducer: getGrpcClient(),
			clientConsumer: getGrpcClient(),
		},

		{
			name:           "send_event_rest",
			clientProducer: getRestClient(),
			clientConsumer: getRestClient(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			m := NewMessage(100)
			errCh := make(chan error, 1)
			msgCh := make(chan *Message, 1)
			channel := uuid.New().String()
			err = tt.clientConsumer.ReceiveEvent(ctx, channel, "", msgCh, errCh)
			require.NoError(t, err)
			time.Sleep(100 * time.Millisecond)
			err = tt.clientProducer.SendEvent(ctx, channel, m)
			require.NoError(t, err)
			select {
			case err := <-errCh:
				require.NoError(t, err)
			case msgIn := <-msgCh:
				log.Print(msgIn.Latency())
			case <-ctx.Done():
				require.NoError(t, ctx.Err())
			}
			cancel()
			err = tt.clientProducer.Close()
			require.NoError(t, err)
			err = tt.clientConsumer.Close()
			require.NoError(t, err)

		})
	}
}

func TestClient_EventStore(t *testing.T) {

	tests := []struct {
		name           string
		clientProducer Transport
		clientConsumer Transport
	}{
		{
			name:           "send_event_store_grpc",
			clientProducer: getGrpcClient(),
			clientConsumer: getGrpcClient(),
		},
		{
			name:           "send_event_store_rest",
			clientProducer: getRestClient(),
			clientConsumer: getRestClient(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			m := NewMessage(100)
			errCh := make(chan error, 1)
			msgCh := make(chan *Message, 1)
			channel := uuid.New().String()
			err = tt.clientConsumer.ReceiveEventStore(ctx, channel, "", msgCh, errCh)
			require.NoError(t, err)
			time.Sleep(100 * time.Millisecond)
			err = tt.clientProducer.SendEventStore(ctx, channel, m)
			require.NoError(t, err)
			select {
			case err := <-errCh:
				require.NoError(t, err)
			case msgIn := <-msgCh:
				log.Print(msgIn.Latency())
			case <-ctx.Done():
				require.NoError(t, ctx.Err())
			}
			cancel()
			err = tt.clientProducer.Close()
			require.NoError(t, err)
			err = tt.clientConsumer.Close()
			require.NoError(t, err)

		})
	}
}

func TestClient_Command(t *testing.T) {

	tests := []struct {
		name           string
		clientProducer Transport
		clientConsumer Transport
		timeout        time.Duration
	}{
		{
			name:           "send_command_grpc",
			clientProducer: getGrpcClient(),
			clientConsumer: getGrpcClient(),
			timeout:        time.Duration(100 * time.Millisecond),
		},
		{
			name:           "send_command_rest",
			clientProducer: getRestClient(),
			clientConsumer: getRestClient(),
			timeout:        time.Duration(100 * time.Millisecond),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			m := NewMessage(100)
			errCh := make(chan error, 1)
			msgCh := make(chan *Message, 1)
			channel := uuid.New().String()
			err = tt.clientConsumer.ReceiveCommand(ctx, channel, "", msgCh, errCh)
			require.NoError(t, err)
			time.Sleep(100 * time.Millisecond)
			err = tt.clientProducer.SendCommand(ctx, channel, m, tt.timeout)
			require.NoError(t, err)
			select {
			case err := <-errCh:
				require.NoError(t, err)
			case msgIn := <-msgCh:
				log.Error(msgIn.Latency())
			case <-ctx.Done():
				require.NoError(t, ctx.Err())
			}
			cancel()
			err = tt.clientProducer.Close()
			require.NoError(t, err)
			err = tt.clientConsumer.Close()
			require.NoError(t, err)

		})
	}
}

func TestClient_Query(t *testing.T) {

	tests := []struct {
		name           string
		clientProducer Transport
		clientConsumer Transport
		timeout        time.Duration
	}{
		{
			name:           "send_query_grpc",
			clientProducer: getGrpcClient(),
			clientConsumer: getGrpcClient(),
			timeout:        time.Duration(100 * time.Millisecond),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			m := NewMessage(100)
			errCh := make(chan error, 1)
			msgCh := make(chan *Message, 1)
			channel := uuid.New().String()
			err = tt.clientConsumer.ReceiveQuery(ctx, channel, "", msgCh, errCh)
			require.NoError(t, err)
			time.Sleep(100 * time.Millisecond)
			resp, err := tt.clientProducer.SendQuery(ctx, channel, m, tt.timeout)
			require.NoError(t, err)
			select {
			case err := <-errCh:
				require.NoError(t, err)
			case msgIn := <-msgCh:
				log.Print(msgIn.Latency())
			case <-ctx.Done():
				require.NoError(t, ctx.Err())
			}
			require.NotNil(t, resp)
			log.Print(resp.Latency())
			cancel()
			err = tt.clientProducer.Close()
			require.NoError(t, err)
			err = tt.clientConsumer.Close()
			require.NoError(t, err)

		})
	}
}
