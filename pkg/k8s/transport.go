package k8s

import (
	"context"
	"errors"
	"fmt"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s/client"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"time"
)

func SetTransport(ctx context.Context, cfg *config.Config) error {
	if !cfg.AutoIntegrated {
		return nil
	}
	utils.Print("connecting to kuberenets cluster... ")
	c, err := client.NewClient(cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	ports := []string{
		fmt.Sprintf("%d:%d", cfg.GrpcPort, cfg.GrpcPort),
		fmt.Sprintf("%d:%d", cfg.RestPort, cfg.RestPort),
		fmt.Sprintf("%d:%d", cfg.ApiPort, cfg.ApiPort),
	}
	stopCh := make(chan struct{})
	outCh, errOutCh := make(chan string, 1), make(chan string, 1)
	err = c.ForwardPorts(cfg.CurrentNamespace, cfg.CurrentStatefulSet+"-0", ports, stopCh, outCh, errOutCh)
	if err != nil {
		return err
	}
	select {
	case <-outCh:
		utils.Println("ok")
	case errstr := <-errOutCh:
		return fmt.Errorf(errstr)

	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout during setting of transport layer to Kubeernetes cluster")

	case <-ctx.Done():
		close(stopCh)
		return nil
	}

	go func() {
		for {
			select {
			case <-outCh:

			case errstr := <-errOutCh:
				utils.CheckErr(errors.New(errstr))
				return
			case <-ctx.Done():
				close(stopCh)
				return
			}
		}
	}()
	return nil
}
