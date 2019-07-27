package grpc

import (
	"context"
	"errors"
	"time"

	kubemq "github.com/kubemq-io/kubemq-go"

	"github.com/kubemq-io/kubetools/transport"

	"github.com/google/uuid"
	"github.com/kubemq-io/kubetools/transport/option"
)

// Client
type Client struct {
	client *kubemq.Client
}

// New - create new client
func New(ctx context.Context, opts *option.Options) (*Client, error) {
	c := &Client{
		client: nil,
	}
	var err error
	if opts.IsSecured {
		c.client, err = kubemq.NewClient(ctx,
			kubemq.WithAddress(opts.Host, opts.Port),
			kubemq.WithClientId(uuid.New().String()),
			kubemq.WithTransportType(kubemq.TransportTypeGRPC),
			kubemq.WithCredentials(opts.CertFile, ""))
	} else {
		c.client, err = kubemq.NewClient(ctx,
			kubemq.WithAddress(opts.Host, opts.Port),
			kubemq.WithClientId(uuid.New().String()),
			kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	}

	if err != nil {
		return nil, err
	}
	return c, nil
}

// Close - close client connection
func (c *Client) Close() error {
	return c.client.Close()
}

// SendEvent - send event message
func (c *Client) SendEvent(ctx context.Context, channel string, m *transport.Message) error {
	m.SetSendTime()
	return c.client.E().
		SetId(m.Id).
		SetChannel(channel).
		SetBody(m.Marshal()).
		Send(ctx)

}

// SendEventStore - send event store message
func (c *Client) SendEventStore(ctx context.Context, channel string, m *transport.Message) error {
	m.SetSendTime()
	result, err := c.client.ES().
		SetId(m.Id).
		SetChannel(channel).
		SetBody(m.Marshal()).
		Send(ctx)
	if err != nil {
		return err
	}
	if result.Sent {
		return nil
	}
	return result.Err
}

// Send Command - send command message
func (c *Client) SendCommand(ctx context.Context, channel string, m *transport.Message, timeout time.Duration) error {
	m.SetSendTime()
	result, err := c.client.C().
		SetId(m.Id).
		SetChannel(channel).
		SetBody(m.Marshal()).
		SetTimeout(timeout).
		Send(ctx)
	if err != nil {
		return err
	}
	if result.Executed {
		return nil
	}
	return errors.New(result.Error)
}

// Send Query - send query message
func (c *Client) SendQuery(ctx context.Context, channel string, m *transport.Message, timeout time.Duration) (*transport.Message, error) {
	m.SetSendTime()
	result, err := c.client.Q().
		SetId(m.Id).
		SetChannel(channel).
		SetBody(m.Marshal()).
		SetTimeout(timeout).
		Send(ctx)
	if err != nil {
		return nil, err
	}
	if result.Executed {
		m, err := transport.Unmarshal(result.Body)
		if err != nil {
			return nil, err
		}
		m.SetReceiveResponseTime()
		err = m.Validate()
		if err != nil {
			return nil, err
		}
		return m, nil
	}
	return nil, errors.New(result.Error)
}

// ReceiveEvent - start receiving event messages
func (c *Client) ReceiveEvent(ctx context.Context, channel string, group string, rxCh chan *transport.Message, errCh chan error) error {

	eventsCh, err := c.client.SubscribeToEvents(ctx, channel, group, errCh)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case event, more := <-eventsCh:
				if !more {
					return
				}
				m, err := transport.Unmarshal(event.Body)
				if err != nil {
					errCh <- err
					return
				}
				m.SetReceiveTime()
				err = m.Validate()
				if err != nil {
					errCh <- err
					return
				}
				rxCh <- m

			case <-ctx.Done():
				return
			}

		}
	}()
	return nil
}

// ReceiveEventStore - start receiving events store messages
func (c *Client) ReceiveEventStore(ctx context.Context, channel string, group string, rxCh chan *transport.Message, errCh chan error) error {

	eventsStoreCh, err := c.client.SubscribeToEventsStore(ctx, channel, group, errCh, kubemq.StartFromNewEvents())
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case event, more := <-eventsStoreCh:
				if !more {
					return
				}
				m, err := transport.Unmarshal(event.Body)
				if err != nil {
					errCh <- err
					return
				}
				m.SetReceiveTime()
				err = m.Validate()
				if err != nil {
					errCh <- err
					return
				}
				rxCh <- m

			case <-ctx.Done():
				return
			}
		}

	}()
	return nil
}

// ReceiveCommand - start receiving command messages
func (c *Client) ReceiveCommand(ctx context.Context, channel string, group string, rxCh chan *transport.Message, errCh chan error) error {
	commandChannel, err := c.client.SubscribeToCommands(ctx, channel, group, errCh)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case command, more := <-commandChannel:
				if !more {
					return
				}
				m, err := transport.Unmarshal(command.Body)
				if err != nil {
					errCh <- err
					return
				}
				m.SetReceiveTime()
				err = m.Validate()
				if err != nil {
					errCh <- err
					return
				}
				err = c.client.R().SetRequestId(command.Id).
					SetExecutedAt(time.Now()).
					SetResponseTo(command.ResponseTo).
					Send(ctx)
				if err != nil {
					errCh <- err
					return
				}
				rxCh <- m

			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

// ReceiveQuery - start receiving query messages
func (c *Client) ReceiveQuery(ctx context.Context, channel string, group string, rxCh chan *transport.Message, errCh chan error) error {
	queryChannel, err := c.client.SubscribeToQueries(ctx, channel, group, errCh)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case query, more := <-queryChannel:
				if !more {
					return
				}
				m, err := transport.Unmarshal(query.Body)
				if err != nil {
					errCh <- err
					return
				}
				m.SetReceiveTime()
				err = m.Validate()
				if err != nil {
					errCh <- err
					return
				}
				m.SetSendResponseTime()
				err = c.client.R().SetRequestId(query.Id).
					SetExecutedAt(time.Now()).
					SetResponseTo(query.ResponseTo).
					SetBody(m.Marshal()).
					Send(ctx)
				if err != nil {
					errCh <- err
					return
				}
				rxCh <- m

			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
