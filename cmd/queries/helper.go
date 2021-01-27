package queries

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
	Body       string            `json:"body"`
	Executed   string            `json:"executed,omitempty"`
	ExecutedAt string            `json:"executed_at,omitempty"`
	Error      string            `json:"error,omitempty"`
	CacheHit   string            `json:"cache_hit,omitempty"`
	payload    []byte
}

func newObjectWithQueryReceive(query *kubemq.QueryReceive) *object {
	obj := &object{
		Id:         query.Id,
		Channel:    query.Channel,
		ClientId:   query.ClientId,
		Metadata:   query.Metadata,
		Tags:       query.Tags,
		Body:       "",
		Executed:   "",
		ExecutedAt: "",
		Error:      "",
		CacheHit:   "",
		payload:    query.Body,
	}

	sDec, err := b64.StdEncoding.DecodeString(string(query.Body))
	if err != nil {
		obj.Body = string(query.Body)
	} else {
		obj.Body = string(sDec)
	}
	return obj
}
func newObjectWithQueryResponse(response *kubemq.QueryResponse) *object {
	obj := &object{
		Id:         response.QueryId,
		Channel:    "",
		ClientId:   response.ResponseClientId,
		Metadata:   response.Metadata,
		Tags:       response.Tags,
		Body:       "",
		Executed:   strconv.FormatBool(response.Executed),
		ExecutedAt: response.ExecutedAt.Format("2006-01-02 15:04:05.999"),
		Error:      response.Error,
		CacheHit:   strconv.FormatBool(response.CacheHit),
		payload:    response.Body,
	}

	sDec, err := b64.StdEncoding.DecodeString(string(response.Body))
	if err != nil {
		obj.Body = string(response.Body)
	} else {
		obj.Body = string(sDec)
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

func printQueryReceive(query *kubemq.QueryReceive) {
	fmt.Println(newObjectWithQueryReceive(query))
}

func printQueryResponse(response *kubemq.QueryResponse) {
	fmt.Println(newObjectWithQueryResponse(response))
}
