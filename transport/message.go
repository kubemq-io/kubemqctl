package transport

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kubemq-io/kubemq-go/pb"
)

//
type Event = pb.Event
type EventReceive = pb.EventReceive
type Request = pb.Request
type Response = pb.Response
type Empty = pb.Empty
type RequestTypeData = pb.Request_RequestType
type Message struct {
	Id                  string   `json:"id,omitempty"`
	SendTime            int64    `json:"send_time,omitempty"`
	ReceiveTime         int64    `json:"receive_time,omitempty"`
	SendResponseTime    int64    `json:"send_response_time,omitempty"`
	ReceiveResponseTime int64    `json:"receive_response_time,omitempty"`
	Payload             []byte   `json:"payload,omitempty"`
	CRC                 [32]byte `json:"crc,omitempty"`
}

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

func (m *Message) SendLatency() time.Duration {
	rx := time.Unix(0, m.ReceiveTime)
	tx := time.Unix(0, m.SendTime)
	return rx.Sub(tx)
}

func (m *Message) ResponseLatency() time.Duration {
	rx := time.Unix(0, m.ReceiveResponseTime)
	tx := time.Unix(0, m.SendResponseTime)
	return rx.Sub(tx)
}

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

func (m *Message) Validate() error {
	//if sha256.Sum256(m.Payload) != m.CRC {
	//	return errors.New("message corrupted, invalid crc ")
	//}
	return nil
}

func (m *Message) SetPayload(n int) {
	m.Payload, m.CRC = randStringBytes(n)
}

func (m *Message) SetSendTime() {
	m.SendTime = time.Now().UnixNano()
}

func (m *Message) SetReceiveTime() {
	m.ReceiveTime = time.Now().UnixNano()
}

func (m *Message) SetSendResponseTime() {
	m.SendResponseTime = time.Now().UnixNano()
}

func (m *Message) SetReceiveResponseTime() {
	m.ReceiveResponseTime = time.Now().UnixNano()
}

func (m *Message) Marshal() []byte {
	b, _ := json.Marshal(m)
	return b
}

func Unmarshal(data []byte) (*Message, error) {
	m := &Message{}
	var err error
	if data != nil {
		err = json.Unmarshal(data, m)
		return m, err
	}
	return nil, errors.New("empty body received")

}


//func UnmarshalEventProto(data []byte) (*Message, error) {
//
//	pbEvent := &kubepb.Event{}
//	err := pbEvent.Unmarshal(data)
//	if err != nil {
//		return nil, err
//	}
//	pbEvent.String()
//	m := &Message{}
//	err = json.Unmarshal(data, pbEvent.Body)
//	return m, err
//}
