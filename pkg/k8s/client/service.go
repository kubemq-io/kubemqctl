package client

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/utils"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"strings"
	"time"
)

func (c *Client) DeleteServicesForStatefulSet(name string) error {
	pair := strings.Split(name, "/")
	svcs, err := c.ClientSet.CoreV1().Services(pair[0]).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, svc := range svcs.Items {
		if strings.Contains(svc.Name, pair[1]) {
			err := c.ClientSet.CoreV1().Services(pair[0]).Delete(svc.Name, metav1.NewDeleteOptions(0))
			if err != nil {
				utils.Printlnf("Service %s/%s not deleted. Error: %s", svc.Namespace, svc.Namespace, utils.Title(err.Error()))
				continue
			}
			utils.Printlnf("Service %s/%s deleted.", svc.Namespace, svc.Name)
		}
	}
	return nil
}

func (c *Client) CreateOrUpdateService(svc *apiv1.Service) (*apiv1.Service, bool, error) {
	ns, name := svc.Namespace, svc.Name
	svc.ObjectMeta = metav1.ObjectMeta{
		Name:      name,
		Namespace: ns,
	}
	oldSvc, err := c.ClientSet.CoreV1().Services(svc.Namespace).Get(svc.Name, metav1.GetOptions{})
	if err == nil && oldSvc != nil {
		svc.ResourceVersion = oldSvc.ResourceVersion
		svc.Spec.ClusterIP = oldSvc.Spec.ClusterIP
		newSvc, err := c.ClientSet.CoreV1().Services(svc.Namespace).Update(svc)
		if err != nil {
			return nil, true, err
		}
		return newSvc, true, nil
	}

	createdSvc, err := c.ClientSet.CoreV1().Services(svc.Namespace).Create(svc)
	if err != nil {
		return nil, false, err
	}
	return createdSvc, false, nil
}

func svcsToStatus(svcs []apiv1.Service) []*ServiceStatus {
	list := []*ServiceStatus{}
	for _, svc := range svcs {
		ss := &ServiceStatus{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			Type:      string(svc.Spec.Type),
			ClusterIP: svc.Spec.ClusterIP,
			ExternalP: "",
			Ports:     "",
			Age:       time.Now().Sub(svc.CreationTimestamp.Time),
		}

		if string(svc.Spec.Type) == "LoadBalancer" {
			if len(svc.Status.LoadBalancer.Ingress) > 0 {
				ss.ExternalP = svc.Status.LoadBalancer.Ingress[0].IP
			}
		}
		portList := []string{}
		for _, port := range svc.Spec.Ports {
			if port.NodePort > 0 {
				portList = append(portList, fmt.Sprintf("%d:%d", port.Port, port.NodePort))
			} else {
				portList = append(portList, fmt.Sprintf("%d", port.Port))
			}
		}
		ss.Ports = strings.Join(portList, ",")

		list = append(list, ss)
	}
	return list
}

func (c *Client) GetServices(ns string, labels map[string]string) ([]*apiv1.Service, error) {
	svcList := []*apiv1.Service{}
	svcs, err := c.ClientSet.CoreV1().Services(ns).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	if svcs != nil {
		for i := 0; i < len(svcs.Items); i++ {
			svc := svcs.Items[i]
			svc.APIVersion = "v1"
			svc.Kind = "Service"
			for key, value := range svc.Spec.Selector {
				if labels[key] == value {
					svcList = append(svcList, &svc)
					continue
				}
			}
		}

	}
	return svcList, nil
}
