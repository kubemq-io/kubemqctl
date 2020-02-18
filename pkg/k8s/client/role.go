package client

import (
	apiv1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) CreateOrUpdateRole(role *rbac.Role) (*rbac.Role, bool, error) {

	found, err := c.ClientSet.RbacV1().Roles(role.Namespace).Get(role.Name, metav1.GetOptions{})

	if err == nil && found != nil {
		updatedRole, err := c.ClientSet.RbacV1().Roles(role.Namespace).Update(role)
		if err != nil {
			return nil, true, err
		}
		return updatedRole, true, nil
	}

	newRole, err := c.ClientSet.RbacV1().Roles(role.Namespace).Create(role)
	if err != nil {
		return nil, false, err
	}
	return newRole, false, nil
}

func (c *Client) DeleteRole(role *rbac.Role) (bool, error) {
	if role == nil {
		return true, nil
	}
	rolesList, err := c.ClientSet.RbacV1().Roles(role.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, r := range rolesList.Items {
		if r.Name == role.Name {
			err := c.ClientSet.RbacV1().Roles(r.Namespace).Delete(r.Name, metav1.NewDeleteOptions(0))
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

func (c *Client) GetRole(name, namespace string) (*rbac.Role, error) {
	role, err := c.ClientSet.RbacV1().Roles(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return role, nil
}
