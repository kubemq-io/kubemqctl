package targets

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-io/kubemq-go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Request struct {
	Metadata Metadata `json:"metadata"`
	Data     []byte   `json:"data"`
}

func NewRequest() *Request {
	return &Request{
		Metadata: NewMetadata(),
		Data:     nil,
	}
}

func (r *Request) SetMetadata(value Metadata) *Request {
	r.Metadata = value
	return r
}
func (r *Request) SetMetadataKeyValue(key, value string) *Request {
	r.Metadata.Set(key, value)
	return r
}

func (r *Request) SetData(value []byte) *Request {
	r.Data = value
	return r
}
func (r *Request) Size() float64 {
	return float64(len(r.Data))
}

func ParseRequest(body []byte) (*Request, error) {
	if body == nil {
		return nil, fmt.Errorf("empty request")
	}
	req := &Request{}
	err := json.Unmarshal(body, req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *Request) MarshalBinary() []byte {
	data, _ := json.Marshal(r)
	return data
}
func (r *Request) ToEvent() *kubemq.Event {
	return kubemq.NewEvent().
		SetBody(r.MarshalBinary())
}
func (r *Request) ToEventStore() *kubemq.EventStore {
	return kubemq.NewEventStore().
		SetBody(r.MarshalBinary())
}
func (r *Request) ToCommand() *kubemq.Command {
	return kubemq.NewCommand().
		SetBody(r.MarshalBinary())
}
func (r *Request) ToQuery() *kubemq.Query {
	return kubemq.NewQuery().
		SetBody(r.MarshalBinary())
}
func (r *Request) ToQueueMessage() *kubemq.QueueMessage {
	return kubemq.NewQueueMessage().
		SetBody(r.MarshalBinary())
}
func (r *Request) String() string {
	str, err := json.MarshalToString(r)
	if err != nil {
		return ""
	}
	return str
}
