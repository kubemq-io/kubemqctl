package operator

import (
	"context"

	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type serviceAccountManager struct {
	*client.Client
}

func (m *serviceAccountManager) CreateOrUpdateServiceAccount(serviceAccount *apiv1.ServiceAccount) (*apiv1.ServiceAccount, bool, error) {
	found, err := m.ClientSet.CoreV1().ServiceAccounts(serviceAccount.Namespace).Get(context.Background(), serviceAccount.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		updatedServiceAccount, err := m.ClientSet.CoreV1().ServiceAccounts(serviceAccount.Namespace).Update(context.Background(), serviceAccount, metav1.UpdateOptions{})
		if err != nil {
			return nil, true, err
		}
		return updatedServiceAccount, true, nil
	}

	newServiceAccount, err := m.ClientSet.CoreV1().ServiceAccounts(serviceAccount.Namespace).Create(context.Background(), serviceAccount, metav1.CreateOptions{})
	if err != nil {
		return nil, false, err
	}
	return newServiceAccount, false, nil
}

func (m *serviceAccountManager) DeleteServiceAccount(serviceAccount *apiv1.ServiceAccount) error {
	found, err := m.ClientSet.CoreV1().ServiceAccounts(serviceAccount.Namespace).Get(context.Background(), serviceAccount.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		return m.ClientSet.CoreV1().ServiceAccounts(serviceAccount.Namespace).Delete(context.Background(), serviceAccount.Name, metav1.DeleteOptions{})
	}
	return nil
}

func (m *serviceAccountManager) GetServiceAccount(name, namespace string) (*apiv1.ServiceAccount, error) {
	serviceAccount, err := m.ClientSet.CoreV1().ServiceAccounts(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return serviceAccount, nil
}
