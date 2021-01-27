package events_store

import (
	b64 "encoding/base64"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	kubemq "github.com/kubemq-io/kubemq-go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type object struct {
	Id        string            `json:"id"`
	Channel   string            `json:"channel,omitempty"`
	ClientId  string            `json:"client_id,omitempty"`
	Metadata  string            `json:"metadata,omitempty"`
	Timestamp string            `json:"timestamp,omitempty"`
	Sequence  uint64            `json:"sequence,omitempty"`
	Tags      map[string]string `json:"tags,omitempty"`
	Body      string            `json:"body,omitempty"`
	payload   []byte
}

func newObjectWithEventReceive(event *kubemq.EventStoreReceive) *object {
	obj := &object{
		Id:        event.Id,
		Channel:   event.Channel,
		ClientId:  event.ClientId,
		Metadata:  event.Metadata,
		Timestamp: event.Timestamp.Format("2006-01-02 15:04:05.999"),
		Sequence:  event.Sequence,
		Tags:      event.Tags,
		Body:      "",
		payload:   event.Body,
	}
	sDec, err := b64.StdEncoding.DecodeString(string(event.Body))
	if err != nil {
		obj.Body = string(event.Body)
	} else {
		obj.Body = string(sDec)
	}
	return obj
}
func newObjectWithEventStore(event *kubemq.EventStore) *object {
	obj := &object{
		Id:       event.Id,
		Channel:  event.Channel,
		ClientId: event.ClientId,
		Metadata: event.Metadata,
		Tags:     event.Tags,
		Body:     "",
		payload:  event.Body,
	}

	sDec, err := b64.StdEncoding.DecodeString(string(event.Body))
	if err != nil {
		obj.Body = string(event.Body)
	} else {
		obj.Body = string(sDec)
	}
	return obj
}

func (o *object) String() string {
	data, _ := json.MarshalIndent(o, "", " ")
	return string(data)
}

func printEventReceive(event *kubemq.EventStoreReceive) {
	fmt.Println(newObjectWithEventReceive(event))
}
func printEventStore(event *kubemq.EventStore) {
	fmt.Println(newObjectWithEventStore(event))
}
