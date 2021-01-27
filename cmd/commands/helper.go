package commands

import (
	b64 "encoding/base64"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	kubemq "github.com/kubemq-io/kubemq-go"
	"strconv"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type object struct {
	Id         string            `json:"id"`
	Channel    string            `json:"channel,omitempty"`
	ClientId   string            `json:"client_id,omitempty"`
	Metadata   string            `json:"metadata,omitempty"`
	Tags       map[string]string `json:"tags,omitempty"`
	Body       string            `json:"body,omitempty"`
	Executed   string            `json:"executed,omitempty"`
	ExecutedAt string            `json:"executed_at,omitempty"`
	Error      string            `json:"error,omitempty"`
	payload    []byte
}

func newObjectWithCommandReceive(cmd *kubemq.CommandReceive) *object {
	obj := &object{
		Id:         cmd.Id,
		Channel:    cmd.Channel,
		ClientId:   cmd.ClientId,
		Metadata:   cmd.Metadata,
		Tags:       cmd.Tags,
		Body:       "",
		Executed:   "",
		ExecutedAt: "",
		Error:      "",
		payload:    cmd.Body,
	}

	sDec, err := b64.StdEncoding.DecodeString(string(cmd.Body))
	if err != nil {
		obj.Body = string(cmd.Body)
	} else {
		obj.Body = string(sDec)
	}
	return obj
}
func newObjectWithCommandResponse(response *kubemq.CommandResponse) *object {
	obj := &object{
		Id:         response.CommandId,
		ClientId:   response.ResponseClientId,
		Tags:       response.Tags,
		Executed:   strconv.FormatBool(response.Executed),
		ExecutedAt: response.ExecutedAt.Format("2006-01-02 15:04:05.999"),
		Error:      response.Error,
	}
	if !response.Executed {
		obj.ExecutedAt = ""
	}
	return obj
}

func (o *object) String() string {
	data, _ := json.MarshalIndent(o, "", " ")
	return string(data)
}

func printCommandReceive(command *kubemq.CommandReceive) {
	fmt.Println(newObjectWithCommandReceive(command))
}

func printCommandResponse(response *kubemq.CommandResponse) {
	fmt.Println(newObjectWithCommandResponse(response))
}
