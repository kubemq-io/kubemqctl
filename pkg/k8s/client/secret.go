package client

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) CreateOrUpdateSecret(sec *apiv1.Secret) (*apiv1.Secret, bool, error) {
	ns, name := sec.Namespace, sec.Name
	sec.ObjectMeta = metav1.ObjectMeta{
		Name:      name,
		Namespace: ns,
	}
	oldSec, err := c.ClientSet.CoreV1().Secrets(sec.Namespace).Get(sec.Name, metav1.GetOptions{})
	if err == nil && oldSec != nil {
		newSec, err := c.ClientSet.CoreV1().Secrets(sec.Namespace).Update(sec)
		if err != nil {
			return nil, true, err
		}
		return newSec, true, nil
	}
	createSec, err := c.ClientSet.CoreV1().Secrets(sec.Namespace).Create(sec)
	if err != nil {
		return nil, false, err
	}
	return createSec, false, nil
}

func (c *Client) GetSecrets(ns string, volumes []apiv1.Volume) ([]*apiv1.Secret, error) {

	secList := []*apiv1.Secret{}
	secs, err := c.ClientSet.CoreV1().Secrets(ns).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	if secs != nil {
		for i := 0; i < len(secs.Items); i++ {
			sec := secs.Items[i]
			sec.APIVersion = "v1"
			sec.Kind = "Secret"
			for _, vm := range volumes {
				if vm.Secret != nil {
					if vm.Secret.SecretName == sec.Name {
						secList = append(secList, &sec)
						continue
					}

				}

			}
		}
	}
	return secList, nil
}
