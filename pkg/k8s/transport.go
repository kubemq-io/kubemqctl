package k8s

import (
	"context"
	"errors"
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	v1 "k8s.io/api/core/v1"
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
	utils.Print("connecting to kubernetes cluster... ")
	c, err := client.NewClient(cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	podNameSpace, podName, err := GetRunningPod(c, cfg.CurrentNamespace, cfg.CurrentStatefulSet)
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
	err = c.ForwardPorts(podNameSpace, podName, ports, stopCh, outCh, errOutCh)
	if err != nil {
		return err
	}
	select {
	case <-outCh:
		utils.Printlnf("->  connected to %s/%s at gRPC Port %s Rest Port %s Api Port %s, ok", podNameSpace, podName, ports[0], ports[1], ports[2])
	case errstr := <-errOutCh:
		return fmt.Errorf(errstr)

	case <-time.After(30 * time.Second):
		return fmt.Errorf("timeout during setting of transport layer to kubernetes cluster")

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

func GetRunningPod(client *client.Client, ns, sts string) (string, string, error) {
	pods, err := client.GetPods(ns, sts)
	if err != nil {
		return "", "", err
	}
	list := NewRandomList()
	for _, pod := range pods {
		if pod.Status.Phase == v1.PodRunning {
			list.Add(pod.Name)
		}
	}
	randPort := list.Random()
	if randPort != "" {
		return ns, randPort, nil
	}

	return "", "", fmt.Errorf("no running pods available in %s/%s. you can change the current context with 'kubemqctl config' command", ns, sts)
}
