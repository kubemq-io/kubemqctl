package client

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) CreateOrUpdateOperator(operator *appsv1.Deployment) (*appsv1.Deployment, bool, error) {
	found, err := c.ClientSet.AppsV1().Deployments(operator.Namespace).Get(operator.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		updatedOperator, err := c.ClientSet.AppsV1().Deployments(operator.Namespace).Update(operator)
		if err != nil {
			return nil, true, err
		}
		return updatedOperator, true, nil
	}

	newOperator, err := c.ClientSet.AppsV1().Deployments(operator.Namespace).Create(operator)
	if err != nil {
		return nil, false, err
	}
	return newOperator, false, nil
}

func (c *Client) DeleteOperator(operator *appsv1.Deployment) (bool, error) {
	if operator == nil {
		return true, nil
	}
	operatorsList, err := c.ClientSet.AppsV1().Deployments(operator.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, op := range operatorsList.Items {
		if op.Name == operator.Name {
			err := c.ClientSet.AppsV1().Deployments(op.Namespace).Delete(op.Name, metav1.NewDeleteOptions(0))
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

func (c *Client) GetOperator(name, namespace string) (*appsv1.Deployment, error) {
	operator, err := c.ClientSet.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return operator, nil
}
