package k8s

import (
	"context"
	"errors"
	"fmt"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s/client"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"net"
	"time"
)

func getFreePorts(count int) ([]int, error) {
	var ports []int
	for i := 0; i < count; i++ {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			return nil, err
		}

		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return nil, err
		}
		defer l.Close()
		ports = append(ports, l.Addr().(*net.TCPAddr).Port)
	}
	return ports, nil
}

func SetTransport(ctx context.Context, cfg *config.Config) error {
	if !cfg.AutoIntegrated {
		return nil
	}
	utils.Print("connecting to kuberenets cluster... ")
	c, err := client.NewClient(cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	freePorts, err := getFreePorts(3)
	if err != nil {
		return err
	}

	ports := []string{
		fmt.Sprintf("%d:%d", freePorts[0], cfg.GrpcPort),
		fmt.Sprintf("%d:%d", freePorts[1], cfg.RestPort),
		fmt.Sprintf("%d:%d", freePorts[2], cfg.ApiPort),
	}
	cfg.GrpcPort = freePorts[0]
	cfg.RestPort = freePorts[1]
	cfg.ApiPort = freePorts[2]

	stopCh := make(chan struct{})
	outCh, errOutCh := make(chan string, 1), make(chan string, 1)
	err = c.ForwardPorts(cfg.CurrentNamespace, cfg.CurrentStatefulSet+"-0", ports, stopCh, outCh, errOutCh)
	if err != nil {
		return err
	}
	select {
	case <-outCh:
		utils.Printlnf("-> gRPC Port %s Rest Port %s Api Port %s, ok", ports[0], ports[1], ports[2])
	case errstr := <-errOutCh:
		return fmt.Errorf(errstr)

	case <-time.After(30 * time.Second):
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
