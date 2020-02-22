package cluster

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	"github.com/spf13/cobra"
)

var defaultQueueConfig = &deployQueueOptions{
	maxReceiveMessagesRequest: 1024,
	maxWaitTimeoutSeconds:     3600,
	maxExpirationSeconds:      43200,
	maxDelaySeconds:           43200,
	maxReQueues:               1024,
	maxVisibilitySeconds:      43200,
	defaultVisibilitySeconds:  60,
	defaultWaitTimeoutSeconds: 1,
}

type deployQueueOptions struct {
	maxReceiveMessagesRequest int32
	maxWaitTimeoutSeconds     int32
	maxExpirationSeconds      int32
	maxDelaySeconds           int32
	maxReQueues               int32
	maxVisibilitySeconds      int32
	defaultVisibilitySeconds  int32
	defaultWaitTimeoutSeconds int32
}

func setQueueConfig(cmd *cobra.Command) *deployQueueOptions {
	o := &deployQueueOptions{
		maxReceiveMessagesRequest: 0,
		maxWaitTimeoutSeconds:     0,
		maxExpirationSeconds:      0,
		maxDelaySeconds:           0,
		maxReQueues:               0,
		maxVisibilitySeconds:      0,
		defaultVisibilitySeconds:  0,
		defaultWaitTimeoutSeconds: 0,
	}

	cmd.PersistentFlags().Int32VarP(&o.maxReceiveMessagesRequest, "queue-max-receive-messages-request", "", 1024, "set max of sending / receiving batch of queue message ")
	cmd.PersistentFlags().Int32VarP(&o.maxWaitTimeoutSeconds, "queue-max-wait-timeout-seconds", "", 3600, "set max wait timeout allowed for message")
	cmd.PersistentFlags().Int32VarP(&o.maxExpirationSeconds, "queue-max-expiration-seconds", "", 43200, "set max expiration allowed for message")
	cmd.PersistentFlags().Int32VarP(&o.maxDelaySeconds, "queue-max-delay-seconds", "", 43200, "set max delay seconds allowed for message")
	cmd.PersistentFlags().Int32VarP(&o.maxReQueues, "queue-max-requeues", "", 1024, "set max retires to receive message before discard")
	cmd.PersistentFlags().Int32VarP(&o.maxVisibilitySeconds, "queue-max-visibility-seconds", "", 43200, "set max time of hold received message before returning to queue")
	cmd.PersistentFlags().Int32VarP(&o.defaultVisibilitySeconds, "queue-default-visibility-seconds", "", 60, "set default time of hold received message before returning to queue")
	cmd.PersistentFlags().Int32VarP(&o.defaultWaitTimeoutSeconds, "queue-default-wait-timeout-seconds", "", 1, "set default time to wait for a message in a queue")

	return o
}

func (o *deployQueueOptions) validate() error {

	return nil
}
func (o *deployQueueOptions) complete() error {
	return nil
}

func (o *deployQueueOptions) setConfig(deployment *kubemqcluster.KubemqCluster) *deployQueueOptions {
	if isDefault(o, defaultQueueConfig) {
		return o
	}

	deployment.Spec.Queue = &kubemqcluster.QueueConfig{
		MaxReceiveMessagesRequest: new(int32),
		MaxWaitTimeoutSeconds:     new(int32),
		MaxExpirationSeconds:      new(int32),
		MaxDelaySeconds:           new(int32),
		MaxReQueues:               new(int32),
		MaxVisibilitySeconds:      new(int32),
		DefaultVisibilitySeconds:  new(int32),
		DefaultWaitTimeoutSeconds: new(int32),
	}
	*deployment.Spec.Queue.MaxReceiveMessagesRequest = o.maxReceiveMessagesRequest
	*deployment.Spec.Queue.MaxWaitTimeoutSeconds = o.maxWaitTimeoutSeconds
	*deployment.Spec.Queue.MaxExpirationSeconds = o.maxExpirationSeconds
	*deployment.Spec.Queue.MaxDelaySeconds = o.maxDelaySeconds
	*deployment.Spec.Queue.MaxReQueues = o.maxReQueues
	*deployment.Spec.Queue.MaxVisibilitySeconds = o.maxVisibilitySeconds
	*deployment.Spec.Queue.DefaultVisibilitySeconds = o.defaultVisibilitySeconds
	*deployment.Spec.Queue.DefaultWaitTimeoutSeconds = o.defaultWaitTimeoutSeconds
	return o
}
