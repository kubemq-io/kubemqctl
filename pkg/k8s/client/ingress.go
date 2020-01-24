package client

import (
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	net "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func (c *Client) DeleteIngressForStatefulSet(name string) error {
	pair := strings.Split(name, "/")
	ingressList, err := c.ClientSet.NetworkingV1beta1().Ingresses(pair[0]).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, ingress := range ingressList.Items {
		if strings.Contains(ingress.Name, pair[1]) {
			err := c.ClientSet.NetworkingV1beta1().Ingresses(pair[0]).Delete(ingress.Name, metav1.NewDeleteOptions(0))
			if err != nil {
				utils.Printlnf("Ingress %s/%s not deleted. Error: %s", ingress.Namespace, ingress.Namespace, utils.Title(err.Error()))
				continue
			}
			utils.Printlnf("Ingress %s/%s deleted.", ingress.Namespace, ingress.Name)
		}
	}
	return nil
}

func (c *Client) CreateOrUpdateIngress(ingress *net.Ingress) (*net.Ingress, bool, error) {
	ns, name := ingress.Namespace, ingress.Name
	ingress.ObjectMeta = metav1.ObjectMeta{
		Name:      name,
		Namespace: ns,
	}
	oldIngress, err := c.ClientSet.NetworkingV1beta1().Ingresses(ingress.Namespace).Get(ingress.Name, metav1.GetOptions{})
	if err == nil && oldIngress != nil {
		ingress.ResourceVersion = oldIngress.ResourceVersion

		newIngress, err := c.ClientSet.NetworkingV1beta1().Ingresses(ingress.Namespace).Update(ingress)
		if err != nil {
			return nil, true, err
		}
		return newIngress, true, nil
	}

	createdIngress, err := c.ClientSet.NetworkingV1beta1().Ingresses(ingress.Namespace).Create(ingress)
	if err != nil {
		return nil, false, err
	}
	return createdIngress, false, nil
}

func (c *Client) GetIngress(ns string, stsName string) ([]*net.Ingress, error) {
	ingressList := []*net.Ingress{}
	list, err := c.ClientSet.NetworkingV1beta1().Ingresses(ns).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	if list != nil {
		for i := 0; i < len(list.Items); i++ {
			ingress := list.Items[i]
			ingress.APIVersion = "networking.k8s.io/v1beta1"
			ingress.Kind = "Ingress"
			if strings.Contains(ingress.Name, stsName) {
				ingressList = append(ingressList, &ingress)
			}
		}

	}
	return ingressList, nil
}
