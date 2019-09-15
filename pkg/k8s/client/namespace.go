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

func (c *Client) CheckAndCreateNamespace(name string) (*apiv1.Namespace, bool, error) {
	ns, err := c.ClientSet.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	if err == nil && ns != nil {
		return ns, false, nil
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
	ns, err = c.ClientSet.CoreV1().Namespaces().Create(newNs)

	if err != nil {
		return nil, false, err
	}

	return ns, true, nil
}
