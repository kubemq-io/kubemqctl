package logs

import (
	"k8s.io/api/core/v1"
)

type ContainerState string

const (
	RUNNING    = "running"
	WAITING    = "waiting"
	TERMINATED = "terminated"
)

func (stateConfig ContainerState) Match(containerState v1.ContainerState) bool {
	return (stateConfig == RUNNING && containerState.Running != nil) ||
		(stateConfig == WAITING && containerState.Waiting != nil) ||
		(stateConfig == TERMINATED && containerState.Terminated != nil)
}
