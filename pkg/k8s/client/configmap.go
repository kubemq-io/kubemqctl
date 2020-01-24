package client

import (
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func (c *Client) DeleteConfigMapsForStatefulSet(name string) error {
	pair := strings.Split(name, "/")
	cms, err := c.ClientSet.CoreV1().ConfigMaps(pair[0]).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, cm := range cms.Items {
		if strings.Contains(cm.Name, pair[1]) {
			err := c.ClientSet.CoreV1().ConfigMaps(pair[0]).Delete(cm.Name, metav1.NewDeleteOptions(0))
			if err != nil {
				utils.Printlnf("ConfigMap %s/%s not deleted. Error: %s", cm.Namespace, cm.Namespace, utils.Title(err.Error()))
				continue
			}
			utils.Printlnf("ConfigMap %s/%s deleted.", cm.Namespace, cm.Name)
		}
	}
	return nil
}

func (c *Client) CreateOrUpdateConfigMap(cm *apiv1.ConfigMap) (*apiv1.ConfigMap, bool, error) {
	ns, name := cm.Namespace, cm.Name
	cm.ObjectMeta = metav1.ObjectMeta{
		Name:      name,
		Namespace: ns,
	}
	oldCm, err := c.ClientSet.CoreV1().ConfigMaps(cm.Namespace).Get(cm.Name, metav1.GetOptions{})
	if err == nil && oldCm != nil {
		newSvc, err := c.ClientSet.CoreV1().ConfigMaps(cm.Namespace).Update(cm)
		if err != nil {
			return nil, true, err
		}
		return newSvc, true, nil
	}
	createCm, err := c.ClientSet.CoreV1().ConfigMaps(cm.Namespace).Create(cm)
	if err != nil {
		return nil, false, err
	}
	return createCm, false, nil
}
func (c *Client) GetConfigMap(namespace, name string) (*apiv1.ConfigMap, error) {

	cm, err := c.ClientSet.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func (c *Client) GetConfigMaps(ns string, stsName string) ([]*apiv1.ConfigMap, error) {
	cmList := []*apiv1.ConfigMap{}
	cms, err := c.ClientSet.CoreV1().ConfigMaps(ns).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	if cms != nil {
		for i := 0; i < len(cms.Items); i++ {
			cm := cms.Items[i]
			cm.Kind = "ConfigMap"
			cm.APIVersion = "v1"
			if strings.Contains(cm.Name, stsName) {
				cmList = append(cmList, &cm)
			}
		}
	}
	return cmList, nil
}
