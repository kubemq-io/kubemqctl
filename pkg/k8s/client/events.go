package client

import (
	"context"
	apiv1 "k8s.io/api/core/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"strings"
	"time"
)

func (c *Client) GetEvents(ctx context.Context, evtCh chan *apiv1.Event, done chan struct{}) error {
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(c.ClientSet, time.Second*5)
	eventsInformer := kubeInformerFactory.Core().V1().Events().Informer()
	stop := make(chan struct{})
	defer close(stop)
	eventsInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ev, ok := obj.(*apiv1.Event)
			if ok {
				evtCh <- ev
			}

		},
		DeleteFunc: func(obj interface{}) {
			ev, ok := obj.(*apiv1.Event)
			if ok {
				evtCh <- ev
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			ev, ok := newObj.(*apiv1.Event)
			if ok {
				evtCh <- ev
			}
		},
	})
	kubeInformerFactory.Start(stop)
	for {
		select {
		case <-done:
			return nil
		case <-ctx.Done():
			return nil

		}
	}
}

func FormatEventSource(es apiv1.EventSource) string {
	EventSourceString := []string{es.Component}
	if len(es.Host) > 0 {
		EventSourceString = append(EventSourceString, es.Host)
	}
	return strings.Join(EventSourceString, ", ")
}
