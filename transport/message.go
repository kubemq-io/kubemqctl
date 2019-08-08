package transport

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/kubemq-io/kubemq-go/pb"
	"time"
)

// Event
type Event = pb.Event

// EventReceive
type EventReceive = pb.EventReceive

// Request
type Request = pb.Request

// Response
type Response = pb.Response

// Empty
type Empty = pb.Empty

// RequestTypeData
type RequestTypeData = pb.Request_RequestType

// Message - base message
type Message struct {
	Id                  string   `json:"id,omitempty"`
	SendTime            int64    `json:"send_time,omitempty"`
	ReceiveTime         int64    `json:"receive_time,omitempty"`
	SendResponseTime    int64    `json:"send_response_time,omitempty"`
	ReceiveResponseTime int64    `json:"receive_response_time,omitempty"`
	Payload             []byte   `json:"payload,omitempty"`
	CRC                 [32]byte `json:"crc,omitempty"`
}

// NewMessage - create new message
func NewMessage(n int) *Message {
	m := &Message{
		Id:                  uuid.New().String(),
		SendTime:            0,
		ReceiveTime:         0,
		SendResponseTime:    0,
		ReceiveResponseTime: 0,
		Payload:             nil,
		CRC:                 [32]byte{},
	}
	m.Payload, m.CRC = randStringBytes(n)
	return m
}

// SendLatency
func (m *Message) SendLatency() time.Duration {
	rx := time.Unix(0, m.ReceiveTime)
	tx := time.Unix(0, m.SendTime)
	return rx.Sub(tx)
}

// ResponseLatency
func (m *Message) ResponseLatency() time.Duration {
	rx := time.Unix(0, m.ReceiveResponseTime)
	tx := time.Unix(0, m.SendResponseTime)
	return rx.Sub(tx)
}

// Latency
func (m *Message) Latency() time.Duration {
	var rx time.Time
	if m.ReceiveResponseTime != 0 {
		rx = time.Unix(0, m.ReceiveResponseTime)
	} else {
		rx = time.Unix(0, m.ReceiveTime)
	}

	tx := time.Unix(0, m.SendTime)
	return rx.Sub(tx)
}

// Validate
func (m *Message) Validate() error {
	return nil
}

// SetPayload
func (m *Message) SetPayload(n int) {
	m.Payload, m.CRC = randStringBytes(n)
}

// SetSendTime
func (m *Message) SetSendTime() {
	m.SendTime = time.Now().UnixNano()
}

// SetReceiveTime
func (m *Message) SetReceiveTime() {
	m.ReceiveTime = time.Now().UnixNano()
}

// SetSendResponseTime
func (m *Message) SetSendResponseTime() {
	m.SendResponseTime = time.Now().UnixNano()
}

// SetReceiveResponseTime
func (m *Message) SetReceiveResponseTime() {
	m.ReceiveResponseTime = time.Now().UnixNano()
}

// Marshal
func (m *Message) Marshal() []byte {
	b, _ := json.Marshal(m)
	return b
}

// Unmarshal
func Unmarshal(data []byte) (*Message, error) {
	m := &Message{}
	var err error
	if data != nil {
		err = json.Unmarshal(data, m)
		return m, err
	}
	return nil, errors.New("empty body received")

}
