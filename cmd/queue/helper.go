package queue

import (
	b64 "encoding/base64"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	kubemq "github.com/kubemq-io/kubemq-go"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type queueMessageObject struct {
	Id        string            `json:"id"`
	Channel   string            `json:"channel"`
	ClientId  string            `json:"client_id"`
	Timestamp string            `json:"timestamp"`
	Sequence  uint64            `json:"sequence"`
	Metadata  string            `json:"metadata"`
	Tags      map[string]string `json:"tags,omitempty"`
	Body      string            `json:"body"`
	payload   []byte
}

func newQueueMessageObject(msg *kubemq.QueueMessage) *queueMessageObject {
	obj := &queueMessageObject{
		Id:        msg.MessageID,
		Channel:   msg.Channel,
		ClientId:  msg.ClientID,
		Timestamp: "",
		Sequence:  0,
		Metadata:  msg.Metadata,
		Tags:      msg.Tags,
		Body:      "",
		payload:   msg.Body,
	}
	if msg.Attributes != nil {
		obj.Timestamp = time.Unix(0, msg.Attributes.Timestamp).Format("2006-01-02 15:04:05.999")
		obj.Sequence = msg.Attributes.Sequence
	}
	sDec, err := b64.StdEncoding.DecodeString(string(msg.Body))
	if err != nil {
		obj.Body = string(msg.Body)
	} else {
		obj.Body = string(sDec)
	}
	return obj
}

func (o *queueMessageObject) String() string {
	data, _ := json.MarshalIndent(o, "", " ")
	return string(data)
}

func printItems(items []*kubemq.QueueMessage) {
	for _, item := range items {
		fmt.Println(newQueueMessageObject(item))
	}
}
