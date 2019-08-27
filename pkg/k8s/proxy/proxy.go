package proxy

import (
	"context"
	"errors"
	"fmt"
	"github.com/kubemq-io/kubetools/pkg/k8s/client"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"time"
)

type ProxyOptions struct {
	KubeConfig string
	Namespace  string
	Pod        string
	Ports      []string
}

func SetProxy(ctx context.Context, opts *ProxyOptions) error {
	utils.Print("connecting to kuberenets cluster... ")
	c, err := client.NewClient(opts.KubeConfig)
	if err != nil {
		return err
	}
	stopCh := make(chan struct{})
	outCh, errOutCh := make(chan string, 1), make(chan string, 1)
	err = c.ForwardPorts(opts.Namespace, opts.Pod, opts.Ports, stopCh, outCh, errOutCh)
	if err != nil {
		return err
	}
	select {
	case <-outCh:
		utils.Println("ok")
	case errstr := <-errOutCh:
		return fmt.Errorf(errstr)

	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout during setting of proxy layer to Kubeernetes cluster")

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
