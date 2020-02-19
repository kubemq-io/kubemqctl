package client

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) GetNamespace(name string) (*apiv1.Namespace, bool, error) {
	ns, err := c.ClientSet.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	if err == nil && ns != nil {
		return ns, true, nil
	}
	newNs := &apiv1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec:   apiv1.NamespaceSpec{},
		Status: apiv1.NamespaceStatus{},
	}
	return newNs, false, nil

}

func (c *Client) CheckAndCreateNamespace(namespace *apiv1.Namespace) (*apiv1.Namespace, bool, error) {
	ns, err := c.ClientSet.CoreV1().Namespaces().Get(namespace.Name, metav1.GetOptions{})
	if err == nil && ns != nil {
		return ns, false, nil
	}

	createNs, err := c.ClientSet.CoreV1().Namespaces().Create(namespace)
	if err != nil {
		return nil, false, err
	}

	return createNs, true, nil
}
func (c *Client) CreateIfNotPresentNamespace(namespace string) error {

	ns, err := c.ClientSet.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if err == nil && ns != nil {
		return nil
	}

	ns = &apiv1.Namespace{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
		Spec:   apiv1.NamespaceSpec{},
		Status: apiv1.NamespaceStatus{},
	}
	_, err = c.ClientSet.CoreV1().Namespaces().Create(ns)
	if err != nil {
		return err
	}

	return nil
}
func (c *Client) GetNamespaceList() ([]string, error) {
	list, err := c.ClientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var nsList []string
	for _, item := range list.Items {
		nsList = append(nsList, item.Name)
	}
	return nsList, nil

}
