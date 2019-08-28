package logs

import (
	"context"
	"fmt"
	"regexp"

	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/typed/core/v1"
)

// Target is a target to attach
type Target struct {
	Namespace string
	Pod       string
	Container string
}

// GetID returns the ID of the object
func (t *Target) GetID() string {
	return fmt.Sprintf("%s-%s-%s", t.Namespace, t.Pod, t.Container)
}

// Watch starts listening to Kubernetes events and emits modified
// containers/pods. The first result is targets added, the second is targets
// removed
func Watch(ctx context.Context, i v1.PodInterface, podFilter *regexp.Regexp, containerFilter *regexp.Regexp, labelSelector labels.Selector) (chan *Target, chan *Target, error) {
	watcher, err := i.Watch(metav1.ListOptions{Watch: true, LabelSelector: labelSelector.String()})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to set up attach")
	}

	added := make(chan *Target)
	removed := make(chan *Target)

	go func() {
		for {
			select {
			case e := <-watcher.ResultChan():
				if e.Object == nil {
					// Closed because of error
					return
				}

				pod := e.Object.(*corev1.Pod)

				if !podFilter.MatchString(pod.Name) {
					continue
				}

				switch e.Type {
				case watch.Added, watch.Modified:
					var statuses []corev1.ContainerStatus
					statuses = append(statuses, pod.Status.InitContainerStatuses...)
					statuses = append(statuses, pod.Status.ContainerStatuses...)

					for _, c := range statuses {
						if !containerFilter.MatchString(c.Name) {
							continue
						}

						if c.State.Running != nil {
							t := &Target{
								Namespace: pod.Namespace,
								Pod:       pod.Name,
								Container: c.Name,
							}
							added <- t

						}
					}
				case watch.Deleted:
					var containers []corev1.Container
					containers = append(containers, pod.Spec.Containers...)
					containers = append(containers, pod.Spec.InitContainers...)

					for _, c := range containers {
						if !containerFilter.MatchString(c.Name) {
							continue
						}

						removed <- &Target{
							Namespace: pod.Namespace,
							Pod:       pod.Name,
							Container: c.Name,
						}
					}
				}
			case <-ctx.Done():
				watcher.Stop()
				close(added)
				close(removed)
				return
			}
		}
	}()

	return added, removed, nil
}
