package k8s

import (
	"context"
	"errors"
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"time"
)

type ProxyOptions struct {
	KubeConfig  string
	Namespace   string
	StatefulSet string
	Pod         string
	Ports       []string
}

func SetProxy(ctx context.Context, opts *ProxyOptions) error {

	c, err := client.NewClient(opts.KubeConfig)
	if err != nil {
		return err
	}
	if opts.Pod == "" {
		opts.Namespace, opts.Pod, err = GetRunningPod(c, opts.Namespace, opts.StatefulSet)
		if err != nil {
			return err
		}
	}
	utils.Print("connecting to kuberenets cluster... ")
	stopCh := make(chan struct{})
	outCh, errOutCh := make(chan string, 1), make(chan string, 1)
	err = c.ForwardPorts(opts.Namespace, opts.Pod, opts.Ports, stopCh, outCh, errOutCh)
	if err != nil {
		return err
	}

	select {
	case <-outCh:
		utils.Println("ok.")
		utils.Printlnf("start proxy for %s/%s. press CTRL C to close.", opts.Namespace, opts.Pod)
		for _, port := range opts.Ports {
			utils.Printlnf("%s/%s:%s -> 127.0.0.1:%s", opts.Namespace, opts.Pod, port, port)
		}

	case errstr := <-errOutCh:
		return fmt.Errorf(errstr)

	case <-time.After(30 * time.Second):
		return fmt.Errorf("timeout during setting of proxy layer to kubernetes cluster")

	case <-ctx.Done():
		close(stopCh)
		return nil
	}

	for {
		select {
		case str := <-outCh:
			utils.Println(str)
		case errstr := <-errOutCh:
			utils.CheckErr(errors.New(errstr))
			return nil
		case <-ctx.Done():
			close(stopCh)
			return nil
		}
	}

}
