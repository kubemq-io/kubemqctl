package transport

import (
	"context"
	"time"
)

type TransportType int

const (
	TransportTypeGRPC TransportType = iota
	TransportTypeREST
)

type Transport interface {
	SendEvent(ctx context.Context, channel string, m *Message) error
	SendEventStore(ctx context.Context, channel string, m *Message) error
	SendCommand(ctx context.Context, channel string, m *Message, timeout time.Duration) error
	SendQuery(ctx context.Context, channel string, m *Message, timeout time.Duration) (*Message, error)
	ReceiveEvent(ctx context.Context, channel string, group string, rxCh chan *Message, errCh chan error) error
	ReceiveEventStore(ctx context.Context, channel string, group string, rxCh chan *Message, errCh chan error) error
	ReceiveCommand(ctx context.Context, channel string, group string, rxCh chan *Message, errCh chan error) error
	ReceiveQuery(ctx context.Context, channel string, group string, rxCh chan *Message, errCh chan error) error
	Close() error
}
