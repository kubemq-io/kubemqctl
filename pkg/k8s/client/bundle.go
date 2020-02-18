package client

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/operator"
)

func (c *Client) CreateOrUpdateBundle(bundle *operator.Bundle) (*operator.Bundle, bool, error) {
	newBundle := &operator.Bundle{
		Name:           bundle.Name,
		Namespace:      bundle.Namespace,
		CRDs:           nil,
		Deployment:     nil,
		Role:           nil,
		RoleBinding:    nil,
		ServiceAccount: nil,
	}
	isBundleUpdated := false
	serviceAccount, isUpdated, err := c.CreateOrUpdateServiceAccount(bundle.ServiceAccount)
	if err != nil {
		return nil, false, fmt.Errorf("error create or update service account, error: %s", err.Error())
	}
	newBundle.ServiceAccount = serviceAccount
	if isUpdated {
		isBundleUpdated = true
	}

	role, isUpdated, err := c.CreateOrUpdateRole(bundle.Role)
	if err != nil {
		return nil, false, fmt.Errorf("error create or update role, error: %s", err.Error())
	}
	newBundle.Role = role
	if isUpdated {
		isBundleUpdated = true
	}

	roleBinding, isUpdated, err := c.CreateOrUpdateRoleBinding(bundle.RoleBinding)
	if err != nil {
		return nil, false, fmt.Errorf("error create or update role binding, error: %s", err.Error())
	}
	newBundle.RoleBinding = roleBinding
	if isUpdated {
		isBundleUpdated = true
	}

	for _, crd := range bundle.CRDs {
		newCrd, isUpdated, err := c.CreateOrUpdateCRD(crd)
		if err != nil {
			return nil, false, fmt.Errorf("error create or update crd, error: %s", err.Error())
		}
		newBundle.CRDs = append(newBundle.CRDs, newCrd)
		if isUpdated {
			isBundleUpdated = true
		}
	}

	deployment, isUpdated, err := c.CreateOrUpdateOperator(bundle.Deployment)
	if err != nil {
		return nil, false, fmt.Errorf("error create or update operator deployment, error: %s", err.Error())
	}
	newBundle.Deployment = deployment
	if isUpdated {
		isBundleUpdated = true
	}
	return newBundle, isBundleUpdated, nil
}

func (c *Client) DeleteBundle(bundle *operator.Bundle) (bool, error) {
	if bundle == nil {
		return true, nil
	}

	_, err := c.DeleteOperator(bundle.Deployment)
	if err != nil {
		return false, fmt.Errorf("delete operator failed, error: %s", err.Error())
	}

	for _, crd := range bundle.CRDs {
		_, err := c.DeleteCrd(crd)
		if err != nil {
			return false, fmt.Errorf("delete crd failed, error: %s", err.Error())
		}
	}

	_, err = c.DeleteRole(bundle.Role)
	if err != nil {
		return false, fmt.Errorf("delete role failed, error: %s", err.Error())
	}

	_, err = c.DeleteRoleBinding(bundle.RoleBinding)
	if err != nil {
		return false, fmt.Errorf("delete role binding failed, error: %s", err.Error())
	}

	_, err = c.DeleteServiceAccount(bundle.ServiceAccount)
	if err != nil {
		return false, fmt.Errorf("delete service acount failed, error: %s", err.Error())
	}

	return true, nil
}

func (c *Client) GetBundle(name, namespace string) (*operator.Bundle, error) {
	bundle := &operator.Bundle{
		Name:           name,
		Namespace:      namespace,
		CRDs:           nil,
		Deployment:     nil,
		Role:           nil,
		RoleBinding:    nil,
		ServiceAccount: nil,
	}

	bundle.Deployment, _ = c.GetOperator(name, namespace)
	crd, _ := c.GetCrd(name, namespace)
	bundle.CRDs = append(bundle.CRDs, crd)
	bundle.Role, _ = c.GetRole(name, namespace)
	bundle.RoleBinding, _ = c.GetRoleBinding(name, namespace)
	bundle.ServiceAccount, _ = c.GetServiceAccount(name, namespace)

	return bundle, nil
}
