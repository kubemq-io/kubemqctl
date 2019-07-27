package kubetest

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kubemq-io/kubetools/transport/option"
)

//ScenarioType - represents scenario type
type ScenarioType int

const (
	ScenarioTypeUndefined ScenarioType = iota
	ScenarioTypeEvents
	ScenarioTypeEventsStore
	ScenarioTypeCommands
	ScenarioTypeQueries
)

// ScenarioTypesNames
var ScenarioTypesNames = map[string]ScenarioType{
	"undefined":    ScenarioTypeUndefined,
	"events":       ScenarioTypeEvents,
	"events_store": ScenarioTypeEventsStore,
	"commands":     ScenarioTypeCommands,
	"queries":      ScenarioTypeQueries,
}

// ScenarioTypesStrings
var ScenarioTypesStrings = map[ScenarioType]string{
	ScenarioTypeUndefined:   "undefined",
	ScenarioTypeEvents:      "events",
	ScenarioTypeEventsStore: "events_store",
	ScenarioTypeCommands:    "commands",
	ScenarioTypeQueries:     "queries",
}

//Scenario - represent scenario struct
type Scenario struct {
	Name           string       `json:"name"`
	Type           ScenarioType `json:"type"`
	Producers      int          `json:"producers"`
	Consumers      int          `json:"consumers"`
	MessageSize    int          `json:"message_size"`
	MessagesToSend int          `json:"messages_to_send"`
	Timeout        int          `json:"timeout"`
}

// Execute
func Execute(ctx context.Context, sc *Scenario, conn *option.Options) {

	if sc.Type == 0 || sc.Type > 4 {
		log.Println("invalid scenario type")
		return
	}

	igConfig := &InstanceGroupConfig{
		ScenarioType:      sc.Type,
		Producers:         sc.Producers,
		Consumers:         sc.Consumers,
		Channel:           uuid.New().String(),
		Group:             "",
		MessageSize:       sc.MessageSize,
		SendMessagesCount: sc.MessagesToSend,
		MessageTimeout:    time.Duration(sc.Timeout) * time.Millisecond,
	}
	ig, err := NewInstanceGroup(ctx, igConfig, conn)
	if err != nil {
		log.Println(err)
	} else {
		resultsStr := []string{fmt.Sprintf("Executing %s: Connection: %s, Prodcuers: %d, Consumers: %d, Payload size: %d bytes, Messages count: %d", sc.Name, option.ConnectionTypeNames[conn.Kind], sc.Producers, sc.Consumers, sc.MessageSize, sc.MessagesToSend)}
		result := ig.Execute(ctx)
		if result.HasErrors() {
			resultsStr = append(resultsStr, result.Errors(), "\t----------------------------")
		} else {
			resultsStr = append(resultsStr, result.String(), "\t----------------------------")
		}
		log.Println(strings.Join(resultsStr, "\n"))
	}

}
