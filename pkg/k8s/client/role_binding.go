package client

import (
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) CreateOrUpdateRoleBinding(roleBinding *rbac.RoleBinding) (*rbac.RoleBinding, bool, error) {

	found, err := c.ClientSet.RbacV1().RoleBindings(roleBinding.Namespace).Get(roleBinding.Name, metav1.GetOptions{})

	if err == nil && found != nil {
		updatedRoleBinding, err := c.ClientSet.RbacV1().RoleBindings(roleBinding.Namespace).Update(roleBinding)
		if err != nil {
			return nil, true, err
		}
		return updatedRoleBinding, true, nil
	}

	newRoleBinding, err := c.ClientSet.RbacV1().RoleBindings(roleBinding.Namespace).Create(roleBinding)
	if err != nil {
		return nil, false, err
	}
	return newRoleBinding, false, nil
}

func (c *Client) DeleteRoleBinding(roleBinding *rbac.RoleBinding) (bool, error) {
	if roleBinding == nil {
		return true, nil
	}
	rolesBindingList, err := c.ClientSet.RbacV1().RoleBindings(roleBinding.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, rb := range rolesBindingList.Items {
		if rb.Name == roleBinding.Name {
			err := c.ClientSet.RbacV1().RoleBindings(rb.Namespace).Delete(rb.Name, metav1.NewDeleteOptions(0))
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

func (c *Client) GetRoleBinding(name, namespace string) (*rbac.RoleBinding, error) {
	role, err := c.ClientSet.RbacV1().RoleBindings(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return role, nil
}
