package kubetest

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kubemq-io/kubetools/transport/option"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestInstanceGroup_EventsExecute(t *testing.T) {

	tests := []struct {
		name      string
		cfg       *InstanceGroupConfig
		connOpts  *option.Options
		timeout   time.Duration
		successes int
		errors    int
	}{
		{
			name: "events, 1 producer, 1 consumer",
			cfg: &InstanceGroupConfig{
				ScenarioType:      ScenarioTypeEvents,
				Producers:         1,
				Consumers:         1,
				Channel:           uuid.New().String(),
				Group:             "",
				MessageSize:       100,
				SendMessagesCount: 1,
				MessageTimeout:    5 * time.Second,
			},
			connOpts: &option.Options{
				Kind:      option.ConnectionTypeGrpc,
				Host:      "localhost",
				Port:      50000,
				IsSecured: false,
				CertFile:  "",
			},
			timeout:   5 * time.Second,
			successes: 2,
			errors:    0,
		},
		{
			name: "events, 1 producer, 10 consumer",
			cfg: &InstanceGroupConfig{
				ScenarioType:      ScenarioTypeEvents,
				Producers:         1,
				Consumers:         10,
				Channel:           uuid.New().String(),
				Group:             "",
				MessageSize:       100,
				SendMessagesCount: 1,
				MessageTimeout:    5 * time.Second,
			},
			connOpts: &option.Options{
				Kind:      option.ConnectionTypeGrpc,
				Host:      "localhost",
				Port:      50000,
				IsSecured: false,
				CertFile:  "",
			},
			timeout:   5 * time.Second,
			successes: 11,
			errors:    0,
		},
		{
			name: "events, 10 producer, 10 consumer",
			cfg: &InstanceGroupConfig{
				ScenarioType:      ScenarioTypeEvents,
				Producers:         10,
				Consumers:         10,
				Channel:           uuid.New().String(),
				Group:             "",
				MessageSize:       100,
				SendMessagesCount: 1,
				MessageTimeout:    5 * time.Second,
			},
			connOpts: &option.Options{
				Kind:      option.ConnectionTypeGrpc,
				Host:      "localhost",
				Port:      50000,
				IsSecured: false,
				CertFile:  "",
			},
			timeout:   5 * time.Second,
			successes: 20,
			errors:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := context.WithTimeout(context.Background(), tt.timeout)
			ig, err := NewInstanceGroup(ctx, tt.cfg, tt.connOpts)
			require.NoError(t, err)
			results := ig.Execute(ctx)
			assert.Zero(t, results.Producers.Errors())
			assert.Zero(t, results.Consumers.Errors())
			assert.Equal(t, tt.successes, results.Producers.Success()+results.Consumers.Success())
			assert.Equal(t, tt.errors, results.Producers.Errors()+results.Consumers.Errors())
			assert.NotZero(t, results.Consumers.Latency())
		})
	}
}

func TestInstanceGroup_EventsStoreExecute(t *testing.T) {

	tests := []struct {
		name      string
		cfg       *InstanceGroupConfig
		connOpts  *option.Options
		timeout   time.Duration
		successes int
		errors    int
	}{
		{
			name: "events store, 1 producer, 1 consumer",
			cfg: &InstanceGroupConfig{
				ScenarioType:      ScenarioTypeEventsStore,
				Producers:         1,
				Consumers:         1,
				Channel:           uuid.New().String(),
				Group:             "",
				MessageSize:       100,
				SendMessagesCount: 1,
				MessageTimeout:    5 * time.Second,
			},
			connOpts: &option.Options{
				Kind:      option.ConnectionTypeGrpc,
				Host:      "localhost",
				Port:      50000,
				IsSecured: false,
				CertFile:  "",
			},
			timeout:   5 * time.Second,
			successes: 2,
			errors:    0,
		},
		{
			name: "events store, 1 producer, 10 consumer",
			cfg: &InstanceGroupConfig{
				ScenarioType:      ScenarioTypeEventsStore,
				Producers:         1,
				Consumers:         10,
				Channel:           uuid.New().String(),
				Group:             "",
				MessageSize:       100,
				SendMessagesCount: 1,
				MessageTimeout:    5 * time.Second,
			},
			connOpts: &option.Options{
				Kind:      option.ConnectionTypeGrpc,
				Host:      "localhost",
				Port:      50000,
				IsSecured: false,
				CertFile:  "",
			},
			timeout:   5 * time.Second,
			successes: 11,
			errors:    0,
		},
		{
			name: "events store, 10 producer, 10 consumer",
			cfg: &InstanceGroupConfig{
				ScenarioType:      ScenarioTypeEventsStore,
				Producers:         10,
				Consumers:         10,
				Channel:           uuid.New().String(),
				Group:             "",
				MessageSize:       100,
				SendMessagesCount: 1,
				MessageTimeout:    5 * time.Second,
			},
			connOpts: &option.Options{
				Kind:      option.ConnectionTypeGrpc,
				Host:      "localhost",
				Port:      50000,
				IsSecured: false,
				CertFile:  "",
			},
			timeout:   5 * time.Second,
			successes: 20,
			errors:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := context.WithTimeout(context.Background(), tt.timeout)
			ig, err := NewInstanceGroup(ctx, tt.cfg, tt.connOpts)
			require.NoError(t, err)
			results := ig.Execute(ctx)
			assert.Zero(t, results.Producers.Errors())
			assert.Zero(t, results.Consumers.Errors())
			assert.Equal(t, tt.successes, results.Producers.Success()+results.Consumers.Success())
			assert.Equal(t, tt.errors, results.Producers.Errors()+results.Consumers.Errors())
			assert.NotZero(t, results.Consumers.Latency())
		})
	}
}

func TestInstanceGroup_CommandExecute(t *testing.T) {

	tests := []struct {
		name      string
		cfg       *InstanceGroupConfig
		connOpts  *option.Options
		timeout   time.Duration
		successes int
		errors    int
	}{
		{
			name: "command, 1 producer, 1 consumer",
			cfg: &InstanceGroupConfig{
				ScenarioType:      ScenarioTypeCommands,
				Producers:         1,
				Consumers:         1,
				Channel:           uuid.New().String(),
				Group:             "",
				MessageSize:       100,
				SendMessagesCount: 1,
				MessageTimeout:    5 * time.Second,
			},
			connOpts: &option.Options{
				Kind:      option.ConnectionTypeGrpc,
				Host:      "localhost",
				Port:      50000,
				IsSecured: false,
				CertFile:  "",
			},
			timeout:   5 * time.Second,
			successes: 2,
			errors:    0,
		},
		{
			name: "commands, 1 producer, 10 consumers",
			cfg: &InstanceGroupConfig{
				ScenarioType:      ScenarioTypeCommands,
				Producers:         1,
				Consumers:         10,
				Channel:           uuid.New().String(),
				Group:             "",
				MessageSize:       100,
				SendMessagesCount: 1,
				MessageTimeout:    5 * time.Second,
			},
			connOpts: &option.Options{
				Kind:      option.ConnectionTypeGrpc,
				Host:      "localhost",
				Port:      50000,
				IsSecured: false,
				CertFile:  "",
			},
			timeout:   5 * time.Second,
			successes: 11,
			errors:    0,
		},
		{
			name: "commands, 10 producers, 10 consumers",
			cfg: &InstanceGroupConfig{
				ScenarioType:      ScenarioTypeCommands,
				Producers:         10,
				Consumers:         10,
				Channel:           uuid.New().String(),
				Group:             "",
				MessageSize:       100,
				SendMessagesCount: 1,
				MessageTimeout:    5 * time.Second,
			},
			connOpts: &option.Options{
				Kind:      option.ConnectionTypeGrpc,
				Host:      "localhost",
				Port:      50000,
				IsSecured: false,
				CertFile:  "",
			},
			timeout:   5 * time.Second,
			successes: 20,
			errors:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := context.WithTimeout(context.Background(), tt.timeout)
			ig, err := NewInstanceGroup(ctx, tt.cfg, tt.connOpts)
			require.NoError(t, err)
			results := ig.Execute(ctx)
			assert.Zero(t, results.Producers.Errors())
			assert.Zero(t, results.Consumers.Errors())
			assert.Equal(t, tt.successes, results.Producers.Success()+results.Consumers.Success())
			assert.Equal(t, tt.errors, results.Producers.Errors()+results.Consumers.Errors())
			assert.NotZero(t, results.Consumers.Latency())
		})
	}
}
func TestInstanceGroup_QueryExecute(t *testing.T) {

	tests := []struct {
		name      string
		cfg       *InstanceGroupConfig
		connOpts  *option.Options
		timeout   time.Duration
		successes int
		errors    int
	}{
		{
			name: "query, 1 producer, 1 consumer",
			cfg: &InstanceGroupConfig{
				ScenarioType:      ScenarioTypeQueries,
				Producers:         1,
				Consumers:         1,
				Channel:           uuid.New().String(),
				Group:             "",
				MessageSize:       100,
				SendMessagesCount: 1,
				MessageTimeout:    5 * time.Second,
			},
			connOpts: &option.Options{
				Kind:      option.ConnectionTypeRest,
				Host:      "localhost",
				Port:      9090,
				IsSecured: false,
				CertFile:  "",
			},
			timeout:   5 * time.Second,
			successes: 2,
			errors:    0,
		},
		{
			name: "query, 1 producer, 10 consumers",
			cfg: &InstanceGroupConfig{
				ScenarioType:      ScenarioTypeQueries,
				Producers:         1,
				Consumers:         10,
				Channel:           uuid.New().String(),
				Group:             "",
				MessageSize:       100,
				SendMessagesCount: 1,
				MessageTimeout:    5 * time.Second,
			},
			connOpts: &option.Options{
				Kind:      option.ConnectionTypeRest,
				Host:      "localhost",
				Port:      9090,
				IsSecured: false,
				CertFile:  "",
			},
			timeout:   5 * time.Second,
			successes: 11,
			errors:    0,
		},
		{
			name: "query, 10 producers, 10 consumers",
			cfg: &InstanceGroupConfig{
				ScenarioType:      ScenarioTypeQueries,
				Producers:         10,
				Consumers:         10,
				Channel:           uuid.New().String(),
				Group:             "",
				MessageSize:       100,
				SendMessagesCount: 1,
				MessageTimeout:    5 * time.Second,
			},
			connOpts: &option.Options{
				Kind:      option.ConnectionTypeRest,
				Host:      "localhost",
				Port:      9090,
				IsSecured: false,
				CertFile:  "",
			},
			timeout:   5 * time.Second,
			successes: 20,
			errors:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := context.WithTimeout(context.Background(), tt.timeout)
			ig, err := NewInstanceGroup(ctx, tt.cfg, tt.connOpts)
			require.NoError(t, err)
			results := ig.Execute(ctx)
			assert.Zero(t, results.Producers.Errors())
			assert.Zero(t, results.Consumers.Errors())
			assert.Equal(t, tt.successes, results.Producers.Success()+results.Consumers.Success())
			assert.Equal(t, tt.errors, results.Producers.Errors()+results.Consumers.Errors())
			assert.NotZero(t, results.Consumers.Latency())
		})
	}
}
