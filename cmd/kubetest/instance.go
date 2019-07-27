package kubetest

import (
	"context"
	"sync"
	"time"

	"github.com/kubemq-io/kubetools/transport"
	"github.com/kubemq-io/kubetools/transport/grpc"
	"github.com/kubemq-io/kubetools/transport/option"
	"github.com/kubemq-io/kubetools/transport/rest"

	"github.com/google/uuid"
)

//ClientType - represent client type
type ClientType int

const (
	ClientTypeEventUnknown ClientType = iota
	ClientTypeEventProducer
	ClientTypeEventConsumer
	ClientTypeEventStoreProducer
	ClientTypeEventStoreConsumer
	ClientTypeCommandProducer
	ClientTypeCommandConsumer
	ClientTypeQueryProducer
	ClientTypeQueryConsumer
)

var ClientTypeNames = map[ClientType]string{
	ClientTypeEventUnknown:       "Unknown",
	ClientTypeEventProducer:      "Event Producer",
	ClientTypeEventConsumer:      "Event Consumer",
	ClientTypeEventStoreProducer: "EventStore Producer",
	ClientTypeEventStoreConsumer: "EventStore Consumer",
	ClientTypeCommandProducer:    "Command Producer",
	ClientTypeCommandConsumer:    "Command Consumer",
	ClientTypeQueryProducer:      "Query Producer",
	ClientTypeQueryConsumer:      "Query Consumer",
}

//InstanceConfig - represents instance config
type InstanceConfig struct {
	Kind              ClientType
	Channel           string
	Group             string
	MessageSize       int
	SendMessagesCount int
	MessageTimeout    time.Duration
}

//InstanceGroupConfig - represents instance group configuration
type InstanceGroupConfig struct {
	ScenarioType      ScenarioType
	Producers         int
	Consumers         int
	Channel           string
	Group             string
	MessageSize       int
	SendMessagesCount int
	MessageTimeout    time.Duration
}

//Instance - represent an instance
type Instance struct {
	id     string
	cfg    *InstanceConfig
	client transport.Transport
}

//NewInstance
func NewInstance(ctx context.Context, cfg *InstanceConfig, connOpts *option.Options) (*Instance, error) {
	ins := &Instance{
		id:     uuid.New().String(),
		cfg:    cfg,
		client: nil,
	}
	var err error
	switch connOpts.Kind {
	case option.ConnectionTypeGrpc:
		ins.client, err = grpc.New(ctx, connOpts)
	case option.ConnectionTypeRest:
		ins.client, err = rest.New(ctx, connOpts)
	}
	if err != nil {
		return nil, err
	}
	return ins, nil
}

func (ins *Instance) executeProducer(ctx context.Context) *Result {

	result := NewResult(ins.id, ins.cfg.SendMessagesCount)
	for i := 0; i < ins.cfg.SendMessagesCount; i++ {
		m := transport.NewMessage(ins.cfg.MessageSize)
		switch ins.cfg.Kind {
		case ClientTypeEventProducer:
			result.SetKind(ClientTypeEventProducer)
			err := ins.client.SendEvent(ctx, ins.cfg.Channel, m)
			if err != nil {
				result.AddError(err)
			}
		case ClientTypeEventStoreProducer:
			result.SetKind(ClientTypeEventStoreProducer)
			err := ins.client.SendEventStore(ctx, ins.cfg.Channel, m)
			if err != nil {
				result.AddError(err)
			}
		case ClientTypeCommandProducer:
			result.SetKind(ClientTypeCommandProducer)
			err := ins.client.SendCommand(ctx, ins.cfg.Channel, m, ins.cfg.MessageTimeout)
			if err != nil {
				result.AddError(err)
			} else {
				result.AddLatency(m.Latency())
			}
		case ClientTypeQueryProducer:
			result.SetKind(ClientTypeQueryProducer)
			resp, err := ins.client.SendQuery(ctx, ins.cfg.Channel, m, ins.cfg.MessageTimeout)
			if err != nil {
				result.AddError(err)
			} else {
				result.AddLatency(resp.Latency())
			}

		}

	}
	return result
}

func (ins *Instance) executeConsumer(ctx context.Context) *Result {
	result := NewResult(ins.id, ins.cfg.SendMessagesCount)
	errCh := make(chan error, 1)
	msgCh := make(chan *transport.Message, 1)
	switch ins.cfg.Kind {
	case ClientTypeEventConsumer:
		result.SetKind(ClientTypeEventConsumer)
		err := ins.client.ReceiveEvent(ctx, ins.cfg.Channel, ins.cfg.Group, msgCh, errCh)
		if err != nil {
			result.AddError(err)
			return result
		}
	case ClientTypeEventStoreConsumer:
		result.SetKind(ClientTypeEventStoreConsumer)
		err := ins.client.ReceiveEventStore(ctx, ins.cfg.Channel, ins.cfg.Group, msgCh, errCh)
		if err != nil {
			result.AddError(err)
			return result
		}

	case ClientTypeCommandConsumer:
		result.SetKind(ClientTypeCommandConsumer)
		err := ins.client.ReceiveCommand(ctx, ins.cfg.Channel, ins.cfg.Group, msgCh, errCh)
		if err != nil {
			result.AddError(err)
			return result
		}

	case ClientTypeQueryConsumer:
		result.SetKind(ClientTypeQueryConsumer)
		err := ins.client.ReceiveQuery(ctx, ins.cfg.Channel, ins.cfg.Group, msgCh, errCh)
		if err != nil {
			result.AddError(err)
			return result
		}
	}

	time.Sleep(100 * time.Millisecond)

	for {

		select {
		case <-errCh:

		case msg := <-msgCh:
			result.AddLatency(msg.Latency())
		case <-ctx.Done():
			return result
		}
	}

}

func (ins *Instance) Execute(ctx context.Context, results chan *Result) {
	defer ins.client.Close()
	switch ins.cfg.Kind {
	case ClientTypeEventProducer, ClientTypeEventStoreProducer, ClientTypeCommandProducer, ClientTypeQueryProducer:
		results <- ins.executeProducer(ctx)
	case ClientTypeEventConsumer, ClientTypeEventStoreConsumer, ClientTypeCommandConsumer, ClientTypeQueryConsumer:
		results <- ins.executeConsumer(ctx)
	}
}

//InstanceGroup - represents an instance group struct
type InstanceGroup struct {
	producers []*Instance
	consumers []*Instance
}

func getInstanceConfig(cfg *InstanceGroupConfig, isProducer bool) *InstanceConfig {
	insCfg := &InstanceConfig{
		Kind:              0,
		Channel:           cfg.Channel,
		Group:             "",
		MessageSize:       cfg.MessageSize,
		SendMessagesCount: cfg.SendMessagesCount,
		MessageTimeout:    cfg.MessageTimeout,
	}
	if isProducer {
		switch cfg.ScenarioType {
		case ScenarioTypeEvents:
			insCfg.Kind = ClientTypeEventProducer
		case ScenarioTypeEventsStore:
			insCfg.Kind = ClientTypeEventStoreProducer
		case ScenarioTypeCommands:
			insCfg.Kind = ClientTypeCommandProducer
		case ScenarioTypeQueries:
			insCfg.Kind = ClientTypeQueryProducer
		}
	} else {
		insCfg.Group = cfg.Group
		switch cfg.ScenarioType {
		case ScenarioTypeEvents:
			insCfg.Kind = ClientTypeEventConsumer
		case ScenarioTypeEventsStore:
			insCfg.Kind = ClientTypeEventStoreConsumer
		case ScenarioTypeCommands:
			insCfg.Kind = ClientTypeCommandConsumer
		case ScenarioTypeQueries:
			insCfg.Kind = ClientTypeQueryConsumer
		}
	}
	return insCfg
}

//NewInstanceGroup
func NewInstanceGroup(ctx context.Context, cfg *InstanceGroupConfig, connOpts *option.Options) (*InstanceGroup, error) {

	ig := &InstanceGroup{
		producers: nil,
		consumers: nil,
	}
	for i := 0; i < cfg.Producers; i++ {
		ins, err := NewInstance(ctx, getInstanceConfig(cfg, true), connOpts)
		if err != nil {
			return nil, err
		}
		ig.producers = append(ig.producers, ins)
	}
	for i := 0; i < cfg.Consumers; i++ {
		ins, err := NewInstance(ctx, getInstanceConfig(cfg, false), connOpts)
		if err != nil {
			return nil, err
		}
		ig.consumers = append(ig.consumers, ins)
	}
	return ig, nil
}

//Execute
func (ig *InstanceGroup) Execute(ctx context.Context) *Results {
	producerResults := NewResultsSet()
	consumerResults := NewResultsSet()
	producerResultsCh := make(chan *Result, len(ig.producers))
	consumerResultsCh := make(chan *Result, len(ig.consumers))

	producerWG := sync.WaitGroup{}
	consumerWG := sync.WaitGroup{}

	producerWG.Add(len(ig.producers))
	consumerWG.Add(len(ig.consumers))
	producersDoneCtx, producersDone := context.WithCancel(ctx)
	go func() {

		for i := 0; i < len(ig.consumers); i++ {
			go ig.consumers[i].Execute(producersDoneCtx, consumerResultsCh)
		}
		for {
			select {
			case result := <-consumerResultsCh:
				consumerResults.Add(result)
				consumerWG.Done()
				if consumerResults.Count() == len(ig.consumers) {
					return
				}

			}
		}

	}()
	time.Sleep(300 * time.Millisecond)
	go func() {
		for i := 0; i < len(ig.producers); i++ {
			go ig.producers[i].Execute(producersDoneCtx, producerResultsCh)
		}
		for {
			select {
			case result := <-producerResultsCh:
				producerResults.Add(result)
				producerWG.Done()
				if producerResults.Count() == len(ig.producers) {
					return
				}

			}
		}

	}()
	producerWG.Wait()
	producersDone()
	consumerWG.Wait()
	return NewResults(producerResults, consumerResults)
}
