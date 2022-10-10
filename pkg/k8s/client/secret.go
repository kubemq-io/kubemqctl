package client

import (
	"context"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func (c *Client) DeleteSecretsForStatefulSet(name string) error {
	pair := strings.Split(name, "/")
	secs, err := c.ClientSet.CoreV1().Secrets(pair[0]).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, sec := range secs.Items {
		if strings.Contains(sec.Name, pair[1]) {
			err := c.ClientSet.CoreV1().Secrets(pair[0]).Delete(context.Background(), sec.Name, metav1.DeleteOptions{})
			if err != nil {
				utils.Printlnf("Secret %s/%s not deleted. Error: %s", sec.Namespace, sec.Namespace, utils.Title(err.Error()))
				continue
			}
			utils.Printlnf("Secret %s/%s deleted.", sec.Namespace, sec.Name)
		}
	}
	return nil
}

func (c *Client) CreateOrUpdateSecret(sec *apiv1.Secret) (*apiv1.Secret, bool, error) {
	ns, name := sec.Namespace, sec.Name
	sec.ObjectMeta = metav1.ObjectMeta{
		Name:      name,
		Namespace: ns,
	}
	oldSec, err := c.ClientSet.CoreV1().Secrets(sec.Namespace).Get(context.Background(), sec.Name, metav1.GetOptions{})
	if err == nil && oldSec != nil {
		newSec, err := c.ClientSet.CoreV1().Secrets(sec.Namespace).Update(context.Background(), sec, metav1.UpdateOptions{})
		if err != nil {
			return nil, true, err
		}
		return newSec, true, nil
	}
	createSec, err := c.ClientSet.CoreV1().Secrets(sec.Namespace).Create(context.Background(), sec, metav1.CreateOptions{})
	if err != nil {
		return nil, false, err
	}
	return createSec, false, nil
}

func (c *Client) GetSecrets(ns string, stsName string) ([]*apiv1.Secret, error) {

	secList := []*apiv1.Secret{}
	secs, err := c.ClientSet.CoreV1().Secrets(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	if secs != nil {
		for i := 0; i < len(secs.Items); i++ {

			sec := secs.Items[i]
			sec.Kind = "Secret"
			sec.APIVersion = "v1"
			if strings.Contains(sec.Name, stsName) {
				secList = append(secList, &sec)
			}
		}
	}
	return secList, nil
}
