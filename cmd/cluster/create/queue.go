package create

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/crd/cluster"
	"github.com/spf13/cobra"
)

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

func defaultQueueConfig(cmd *cobra.Command) *deployQueueOptions {
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

func (o *deployQueueOptions) setConfig(deployment *cluster.KubemqCluster) *deployQueueOptions {
	deployment.Spec.Queue = &cluster.QueueConfig{
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
