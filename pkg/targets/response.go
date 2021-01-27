package targets

import (
	"fmt"
	"github.com/kubemq-io/kubemq-go"
)

type Response struct {
	Metadata Metadata `json:"metadata"`
	Data     []byte   `json:"data"`
	IsError  bool     `json:"is_error"`
	Error    string   `json:"error"`
}

func NewResponse() *Response {
	return &Response{
		Metadata: NewMetadata(),
		Data:     nil,
	}
}

func (r *Response) SetMetadata(value Metadata) *Response {
	r.Metadata = value
	return r
}
func (r *Response) SetMetadataKeyValue(key, value string) *Response {
	r.Metadata.Set(key, value)
	return r
}

func (r *Response) SetData(value []byte) *Response {
	r.Data = value
	return r
}
func (r *Response) SetError(err error) *Response {
	r.IsError = true
	r.Error = err.Error()
	return r
}
func (r *Response) MarshalBinary() []byte {
	data, _ := json.Marshal(r)
	return data
}
func (r *Response) ToEvent() *kubemq.Event {
	return kubemq.NewEvent().
		SetBody(r.MarshalBinary())

}
func (r *Response) ToEventStore() *kubemq.EventStore {
	return kubemq.NewEventStore().
		SetBody(r.MarshalBinary())
}

func (r *Response) ToCommand() *kubemq.Command {
	return kubemq.NewCommand().
		SetMetadata(r.Metadata.String()).
		SetBody(r.Data)
}

func (r *Response) ToQuery() *kubemq.Query {
	return kubemq.NewQuery().
		SetMetadata(r.Metadata.String()).
		SetBody(r.Data)
}

func (r *Response) ToQueueMessage() *kubemq.QueueMessage {
	return kubemq.NewQueueMessage().
		SetBody(r.MarshalBinary())
}
func (r *Response) ToResponse() *kubemq.Response {
	return kubemq.NewResponse().
		SetMetadata(r.Metadata.String()).
		SetBody(r.Data)
}
func (r *Response) Size() float64 {
	return float64(len(r.Data))
}
func (r *Response) String() string {
	str, err := json.MarshalToString(r)
	if err != nil {
		return ""
	}
	return str
}
func parseResponse(meta string, body []byte, errText string) (*Response, error) {
	res := NewResponse()
	parsedMeta, err := UnmarshallMetadata(meta)
	if err != nil {
		return nil, fmt.Errorf("error parsing response metadata, %w", err)
	}
	if errText != "" {
		parsedMeta.Set("error", errText)
	}
	return res.
			SetMetadata(parsedMeta).
			SetData(body),

		nil
}
func ParseResponseFromEvent(event *kubemq.Event) (*Response, error) {
	return parseResponse(event.Metadata, event.Body, "")
}
func ParseResponseFromEventReceive(event *kubemq.EventStoreReceive) (*Response, error) {
	return parseResponse(event.Metadata, event.Body, "")
}
func ParseResponseFromCommandResponse(resp *kubemq.CommandResponse) (*Response, error) {
	return parseResponse("", nil, resp.Error)
}
func ParseResponseFromQueryResponse(resp *kubemq.QueryResponse) (*Response, error) {
	return parseResponse(resp.Metadata, resp.Body, resp.Error)
}
func ParseResponseFromQueueMessage(resp *kubemq.QueueMessage) (*Response, error) {
	return parseResponse(resp.Metadata, resp.Body, "")
}

func ParseResponse(body []byte) (*Response, error) {
	if body == nil {
		return nil, fmt.Errorf("empty response")
	}
	req := &Response{}
	err := json.Unmarshal(body, req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
