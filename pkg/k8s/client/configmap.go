package client

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

func (c *Client) GetConfigMaps(ns string, volumes []apiv1.Volume) ([]*apiv1.ConfigMap, error) {
	cmList := []*apiv1.ConfigMap{}
	cms, err := c.ClientSet.CoreV1().ConfigMaps(ns).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	if cms != nil {
		for i := 0; i < len(cms.Items); i++ {
			cm := cms.Items[i]
			cm.APIVersion = "v1"
			cm.Kind = "ConfigMap"
			for _, vm := range volumes {
				if vm.ConfigMap != nil {
					if vm.ConfigMap.Name == cm.Name {
						cmList = append(cmList, &cm)
						continue
					}

				}
			}
		}
	}
	return cmList, nil
}
