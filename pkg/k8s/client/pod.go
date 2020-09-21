package client

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
)

func PodStatus(p *v1.Pod) string {

	fields := map[string]string{}
	fields["Pod"] = p.Name
	fields["Phase"] = string(p.Status.Phase)

	if len(p.Status.ContainerStatuses) > 0 {
		cs := p.Status.ContainerStatuses[0]
		if cs.State.Running != nil {
			fields["Current Status"] = "Running"
		}
		if cs.State.Waiting != nil {
			fields["Current Status"] = fmt.Sprintf("Waiting (%s)", cs.State.Waiting.Reason)
		}

		if cs.State.Terminated != nil {
			fields["Current Status"] = fmt.Sprintf("Terminated (Code %d) %s, %s", cs.State.Terminated.ExitCode, cs.State.Terminated.Message, cs.State.Terminated.Reason)
		}
		fields["Restarts"] = fmt.Sprintf("%d", cs.RestartCount)
	}

	return fmt.Sprintf("%s Phase: %s Status: %s Restarts: %s", fields["Pod"], fields["Phase"], fields["Current Status"], fields["Restarts"])

}
