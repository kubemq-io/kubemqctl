package client

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) CreateOrUpdateServiceAccount(serviceAccount *apiv1.ServiceAccount) (*apiv1.ServiceAccount, bool, error) {

	found, err := c.ClientSet.CoreV1().ServiceAccounts(serviceAccount.Namespace).Get(serviceAccount.Name, metav1.GetOptions{})

	if err == nil && found != nil {
		updatedServiceAccount, err := c.ClientSet.CoreV1().ServiceAccounts(serviceAccount.Namespace).Update(serviceAccount)
		if err != nil {
			return nil, true, err
		}
		return updatedServiceAccount, true, nil
	}

	newServiceAccount, err := c.ClientSet.CoreV1().ServiceAccounts(serviceAccount.Namespace).Create(serviceAccount)
	if err != nil {
		return nil, false, err
	}
	return newServiceAccount, false, nil
}

func (c *Client) DeleteServiceAccount(serviceAccount *apiv1.ServiceAccount) (bool, error) {
	if serviceAccount == nil {
		return true, nil
	}
	serviceAccountList, err := c.ClientSet.CoreV1().ServiceAccounts(serviceAccount.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, sa := range serviceAccountList.Items {
		if sa.Name == serviceAccount.Name {
			err := c.ClientSet.CoreV1().ServiceAccounts(sa.Namespace).Delete(sa.Name, metav1.NewDeleteOptions(0))
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

func (c *Client) GetServiceAccount(name, namespace string) (*apiv1.ServiceAccount, error) {
	serviceAccount, err := c.ClientSet.CoreV1().ServiceAccounts(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return serviceAccount, nil
}
