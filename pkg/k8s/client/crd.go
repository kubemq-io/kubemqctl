package client

import (
	v1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) CreateOrUpdateCRD(crd *v1beta1.CustomResourceDefinition) (*v1beta1.CustomResourceDefinition, bool, error) {
	found, err := c.ClientApiExtension.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crd.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		updatedCrd, err := c.ClientApiExtension.ApiextensionsV1beta1().CustomResourceDefinitions().Update(crd)
		if err != nil {
			return nil, true, err
		}
		return updatedCrd, true, nil
	}

	newCrd, err := c.ClientApiExtension.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil {
		return nil, false, err
	}
	return newCrd, false, nil
}

func (c *Client) DeleteCrd(crd *v1beta1.CustomResourceDefinition) (bool, error) {
	if crd == nil {
		return true, nil
	}
	crdList, err := c.ClientApiExtension.ApiextensionsV1beta1().CustomResourceDefinitions().List(metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, cr := range crdList.Items {
		if cr.Name == crd.Name {
			err := c.ClientApiExtension.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(cr.Name, metav1.NewDeleteOptions(0))
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

func (c *Client) GetCrd(name, namespace string) (*v1beta1.CustomResourceDefinition, error) {
	crd, err := c.ClientApiExtension.ApiextensionsV1beta1().CustomResourceDefinitions().Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return crd, nil
}
