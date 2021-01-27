package events

import (
	b64 "encoding/base64"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	kubemq "github.com/kubemq-io/kubemq-go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type object struct {
	Id       string            `json:"id"`
	Channel  string            `json:"channel,omitempty"`
	ClientId string            `json:"client_id,omitempty"`
	Metadata string            `json:"metadata,omitempty"`
	Tags     map[string]string `json:"tags,omitempty"`
	Body     string            `json:"body,omitempty"`
	payload  []byte
}

func newObjectWithEvent(event *kubemq.Event) *object {
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

func printEvent(event *kubemq.Event) {
	fmt.Println(newObjectWithEvent(event))
}
