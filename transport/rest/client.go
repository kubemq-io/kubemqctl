package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/kubemq-io/kubetools/transport"

	"github.com/kubemq-io/kubetools/transport/option"

	"github.com/google/uuid"

	"github.com/go-resty/resty"
	"github.com/gorilla/websocket"
)

type Response struct {
	Node        string          `json:"node"`
	Error       bool            `json:"error"`
	ErrorString string          `json:"error_string"`
	Data        json.RawMessage `json:"data"`
}

func newWebsocketConn(ctx context.Context, uri string, readCh chan string, ready chan struct{}, errCh chan error) (*websocket.Conn, error) {
	var c *websocket.Conn
	conn, res, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		buf := make([]byte, 1024)
		if res != nil {
			n, _ := res.Body.Read(buf)
			return nil, errors.New(string(buf[:n]))
		}
		return nil, err

	} else {
		c = conn
	}
	ready <- struct{}{}
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				errCh <- err
				return
			} else {
				readCh <- string(message)
			}
		}
	}()
	return c, nil
}

type Client struct {
	id          string
	restAddress string
	wsAddress   string
	wsConn      *websocket.Conn
}

func New(ctx context.Context, opts *option.Options) (*Client, error) {
	c := &Client{
		id:          uuid.New().String(),
		restAddress: opts.Uri(),
		wsAddress:   opts.WebsocketUri(),
		wsConn:      nil,
	}

	return c, nil
}

func (c *Client) SendEvent(ctx context.Context, channel string, m *transport.Message) error {
	resp := &Response{}
	m.SetSendTime()
	event := transport.Event{
		EventID:  m.Id,
		ClientID: c.id,
		Channel:  channel,
		Metadata: "",
		Body:     m.Marshal(),
		Store:    false,
	}
	uri := fmt.Sprintf("%s/send/event", c.restAddress)
	r, err := resty.R().SetBody(event).SetResult(resp).SetError(resp).Post(uri)
	if err != nil {
		return err
	}
	if !r.IsSuccess() || resp.Error {
		return errors.New(resp.ErrorString)
	}
	return nil
}

func (c *Client) SendEventStore(ctx context.Context, channel string, m *transport.Message) error {
	resp := &Response{}
	m.SetSendTime()
	event := transport.Event{
		EventID:  m.Id,
		ClientID: c.id,
		Channel:  channel,
		Metadata: "",
		Body:     m.Marshal(),
		Store:    true,
	}
	uri := fmt.Sprintf("%s/send/event", c.restAddress)
	r, err := resty.R().SetBody(event).SetResult(resp).SetError(resp).Post(uri)
	if err != nil {
		return err
	}
	if !r.IsSuccess() || resp.Error {
		return errors.New(resp.ErrorString)
	}
	return nil
}

func (c *Client) SendCommand(ctx context.Context, channel string, m *transport.Message, timeout time.Duration) error {
	resp := &Response{}
	m.SetSendTime()
	request := transport.Request{
		RequestID:       m.Id,
		RequestTypeData: 1,
		ClientID:        c.id,
		Channel:         channel,
		Metadata:        "",
		Body:            m.Marshal(),
		Timeout:         int32(timeout.Seconds() * 1e3),
	}
	uri := fmt.Sprintf("%s/send/request", c.restAddress)
	r, err := resty.R().SetBody(request).SetResult(resp).SetError(resp).Post(uri)
	if err != nil {
		return err
	}
	if !r.IsSuccess() || resp.Error {
		return errors.New(resp.ErrorString)
	}
	cmdResult := &transport.Response{}
	err = json.Unmarshal(resp.Data, cmdResult)
	if err != nil {
		return err
	}
	if !cmdResult.Executed {
		return errors.New(cmdResult.Error)
	}

	return nil
}

func (c *Client) SendQuery(ctx context.Context, channel string, m *transport.Message, timeout time.Duration) (*transport.Message, error) {
	resp := &Response{}
	m.SetSendTime()
	request := transport.Request{
		RequestID:       m.Id,
		RequestTypeData: 2,
		ClientID:        c.id,
		Channel:         channel,
		Metadata:        "",
		Body:            m.Marshal(),
		Timeout:         int32(timeout.Seconds() * 1e3),
	}
	uri := fmt.Sprintf("%s/send/request", c.restAddress)
	r, err := resty.R().SetBody(request).SetResult(resp).SetError(resp).Post(uri)
	if err != nil {
		return nil, err
	}
	if !r.IsSuccess() || resp.Error {
		fmt.Println(r.StatusCode())
		return nil, errors.New(resp.ErrorString)
	}
	queryResult := &transport.Response{}
	err = json.Unmarshal(resp.Data, queryResult)
	if err != nil {
		return nil, err
	}
	if !queryResult.Executed {
		return nil, errors.New(queryResult.Error)
	}
	queryMessage, err := transport.Unmarshal(queryResult.Body)
	if err != nil {
		return nil, err
	}
	queryMessage.SetReceiveResponseTime()
	return queryMessage, err
}

func (c *Client) ReceiveEvent(ctx context.Context, channel string, group string, rxCh chan *transport.Message, errCh chan error) error {
	uri := fmt.Sprintf("%s/subscribe/events?&client_id=%s&channel=%s&group=%s&subscribe_type=%s", c.wsAddress, c.id, channel, group, "events")
	rxChan := make(chan string)
	ready := make(chan struct{}, 1)
	wsErrCh := make(chan error, 1)
	conn, err := newWebsocketConn(ctx, uri, rxChan, ready, wsErrCh)
	if err != nil {
		return err
	}
	c.wsConn = conn
	<-ready
	go func() {
		for {
			select {
			case pbMsg := <-rxChan:
				event := &transport.EventReceive{}
				err := json.Unmarshal([]byte(pbMsg), event)
				if err != nil {
					errCh <- err
					return
				}
				m, err := transport.Unmarshal(event.Body)
				if err != nil {
					errCh <- err
					return
				}
				m.SetReceiveTime()
				rxCh <- m

			case err := <-wsErrCh:
				errCh <- err
			case <-ctx.Done():

				return
			}

		}
	}()
	return nil
}

func (c *Client) ReceiveEventStore(ctx context.Context, channel string, group string, rxCh chan *transport.Message, errCh chan error) error {
	uri := fmt.Sprintf("%s/subscribe/events?&client_id=%s&channel=%s&group=%s&subscribe_type=%s&events_store_type_data=%d&events_store_type_value=%d", c.wsAddress, c.id, channel, group, "events_store", 2, 0)
	rxChan := make(chan string)
	ready := make(chan struct{}, 1)
	wsErrCh := make(chan error, 1)
	conn, err := newWebsocketConn(ctx, uri, rxChan, ready, wsErrCh)
	if err != nil {
		return err
	}
	c.wsConn = conn
	<-ready
	go func() {
		for {
			select {
			case pbMsg := <-rxChan:
				event := &transport.EventReceive{}
				err := json.Unmarshal([]byte(pbMsg), event)
				if err != nil {
					errCh <- err
					return
				}
				m, err := transport.Unmarshal(event.Body)
				if err != nil {
					errCh <- err
					return
				}
				m.SetReceiveTime()
				rxCh <- m

			case err := <-wsErrCh:
				errCh <- err
			case <-ctx.Done():

				return
			}

		}
	}()
	return nil

}

func (c *Client) ReceiveCommand(ctx context.Context, channel string, group string, rxCh chan *transport.Message, errCh chan error) error {
	uri := fmt.Sprintf("%s/subscribe/requests?&client_id=%s&channel=%s&group=%s&subscribe_type=%s", c.wsAddress, c.id, channel, group, "commands")
	rxChan := make(chan string)
	ready := make(chan struct{}, 1)
	wsErrCh := make(chan error, 1)
	conn, err := newWebsocketConn(ctx, uri, rxChan, ready, wsErrCh)
	if err != nil {
		return err
	}
	c.wsConn = conn
	<-ready
	go func() {
		for {
			select {
			case pbMsg := <-rxChan:
				request := &transport.Request{}
				err := json.Unmarshal([]byte(pbMsg), request)
				if err != nil {
					errCh <- err
					return
				}
				m, err := transport.Unmarshal(request.Body)
				if err != nil {
					errCh <- err
					return
				}
				m.SetReceiveTime()
				rxCh <- m
				m.SetSendResponseTime()
				cmdResponse := &transport.Response{
					ClientID:     c.id,
					RequestID:    request.RequestID,
					ReplyChannel: request.ReplyChannel,
					Metadata:     "",
					Body:         nil,
					Timestamp:    time.Now().UnixNano(),
					Executed:     true,
					Error:        "",
				}
				cmdResponseResult := &Response{}
				cmdResponseUri := fmt.Sprintf("%s/send/response", c.restAddress)
				_, err = resty.R().SetBody(cmdResponse).SetResult(cmdResponseResult).SetError(cmdResponseResult).Post(cmdResponseUri)
				if err != nil {
					errCh <- err
					return
				}
				if cmdResponseResult.Error {
					errCh <- errors.New(cmdResponseResult.ErrorString)
					return
				}

			case err := <-wsErrCh:
				errCh <- err
			case <-ctx.Done():
				return
			}

		}
	}()
	return nil

}

func (c *Client) ReceiveQuery(ctx context.Context, channel string, group string, rxCh chan *transport.Message, errCh chan error) error {
	uri := fmt.Sprintf("%s/subscribe/requests?&client_id=%s&channel=%s&group=%s&subscribe_type=%s", c.wsAddress, c.id, channel, group, "queries")
	rxChan := make(chan string)
	ready := make(chan struct{}, 1)
	wsErrCh := make(chan error, 1)
	conn, err := newWebsocketConn(ctx, uri, rxChan, ready, wsErrCh)
	if err != nil {
		return err
	}
	c.wsConn = conn
	<-ready
	go func() {
		for {
			select {
			case pbMsg := <-rxChan:
				request := &transport.Request{}
				err := json.Unmarshal([]byte(pbMsg), request)
				if err != nil {
					errCh <- err
					return
				}
				m, err := transport.Unmarshal(request.Body)
				if err != nil {
					errCh <- err
					return
				}
				m.SetReceiveTime()
				rxCh <- m
				m.SetSendResponseTime()
				queryResponse := &transport.Response{
					ClientID:     c.id,
					RequestID:    request.RequestID,
					ReplyChannel: request.ReplyChannel,
					Metadata:     "",
					Body:         m.Marshal(),
					CacheHit:     false,
					Timestamp:    time.Now().UnixNano(),
					Executed:     true,
					Error:        "",
				}
				queryResponseResult := &Response{}
				queryResponseUri := fmt.Sprintf("%s/send/response", c.restAddress)
				_, err = resty.R().SetBody(queryResponse).SetResult(queryResponseResult).SetError(queryResponseResult).Post(queryResponseUri)
				if err != nil {
					errCh <- err
					return
				}
				if queryResponseResult.Error {
					errCh <- errors.New(queryResponseResult.ErrorString)
					return
				}

			case err := <-wsErrCh:
				errCh <- err
			case <-ctx.Done():
				return
			}

		}
	}()
	return nil

}

func (c *Client) Close() error {
	if c.wsConn != nil {
		return c.wsConn.Close()
	}
	return nil
}
