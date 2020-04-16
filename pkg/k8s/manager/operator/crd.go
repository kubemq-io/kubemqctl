package operator

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	v1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type crdManager struct {
	*client.Client
}

func (m *crdManager) CreateOrUpdateCRD(crd *v1beta1.CustomResourceDefinition) (*v1beta1.CustomResourceDefinition, bool, error) {
	found, err := m.ClientApiExtension.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crd.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		crd.ResourceVersion = found.ResourceVersion
		updatedCrd, err := m.ClientApiExtension.ApiextensionsV1beta1().CustomResourceDefinitions().Update(crd)
		if err != nil {
			return nil, true, err
		}
		return updatedCrd, true, nil
	}

	newCrd, err := m.ClientApiExtension.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil {
		return nil, false, err
	}
	return newCrd, false, nil
}

func (m *crdManager) DeleteCrd(crd *v1beta1.CustomResourceDefinition) error {
	found, err := m.ClientApiExtension.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crd.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		return m.ClientApiExtension.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(crd.Name, metav1.NewDeleteOptions(0))
	}
	return nil
}

func (m *crdManager) GetCrd(name, namespace string) (*v1beta1.CustomResourceDefinition, error) {
	crd, err := m.ClientApiExtension.ApiextensionsV1beta1().CustomResourceDefinitions().Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return crd, nil
}
