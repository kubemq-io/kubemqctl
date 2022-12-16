package operator

import (
	"context"

	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type operatorManager struct {
	*client.Client
}

func (m *operatorManager) CreateOrUpdateOperator(operator *appsv1.Deployment) (*appsv1.Deployment, bool, error) {
	found, err := m.ClientSet.AppsV1().Deployments(operator.Namespace).Get(context.Background(), operator.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		updatedOperator, err := m.ClientSet.AppsV1().Deployments(operator.Namespace).Update(context.Background(), operator, metav1.UpdateOptions{})
		if err != nil {
			return nil, true, err
		}
		return updatedOperator, true, nil
	}

	newOperator, err := m.ClientSet.AppsV1().Deployments(operator.Namespace).Create(context.Background(), operator, metav1.CreateOptions{})
	if err != nil {
		return nil, false, err
	}
	return newOperator, false, nil
}

func (m *operatorManager) DeleteOperator(operator *appsv1.Deployment) error {
	found, err := m.ClientSet.AppsV1().Deployments(operator.Namespace).Get(context.Background(), operator.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		return m.ClientSet.AppsV1().Deployments(operator.Namespace).Delete(context.Background(), operator.Name, metav1.DeleteOptions{})
	}
	return nil
}

func (m *operatorManager) GetOperator(name, namespace string) (*appsv1.Deployment, error) {
	operator, err := m.ClientSet.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return operator, nil
}
