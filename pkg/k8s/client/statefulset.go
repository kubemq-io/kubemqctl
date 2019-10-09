package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/yaml"
	"strings"
	"time"
)

func (c *Client) GetStatefulSets(ns string, contains ...string) (map[string]appsv1.StatefulSet, error) {

	sts, err := c.ClientSet.AppsV1().StatefulSets(ns).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	list := map[string]appsv1.StatefulSet{}
	for _, item := range sts.Items {
		if contains == nil {
			list[fmt.Sprintf("%s/%s", item.Namespace, item.Name)] = item
		} else {
			for _, name := range contains {
				if strings.Contains(item.Name, name) {
					list[fmt.Sprintf("%s/%s", item.Namespace, item.Name)] = item
				}
			}
		}

	}
	return list, err
}

func (c *Client) DeleteStatefulSet(name string) error {
	pair := strings.Split(name, "/")
	return c.ClientSet.AppsV1().StatefulSets(pair[0]).Delete(pair[1], metav1.NewDeleteOptions(0))
}

func (c *Client) CreateOrUpdateStatefulSet(sts *appsv1.StatefulSet) (*appsv1.StatefulSet, bool, error) {
	ns, name := sts.Namespace, sts.Name
	sts.ObjectMeta = metav1.ObjectMeta{
		Name:      name,
		Namespace: ns,
	}
	oldSts, err := c.ClientSet.AppsV1().StatefulSets(sts.Namespace).Get(sts.Name, metav1.GetOptions{})
	if err == nil && oldSts != nil {
		oldSts.Spec.Replicas = sts.Spec.Replicas
		oldSts.Spec.Template = sts.Spec.Template
		oldSts.Spec.UpdateStrategy = sts.Spec.UpdateStrategy
		newSts, err := c.ClientSet.AppsV1().StatefulSets(sts.Namespace).Update(oldSts)
		if err != nil {
			return nil, true, err
		}
		return newSts, true, nil
	}
	createdSts, err := c.ClientSet.AppsV1().StatefulSets(sts.Namespace).Create(sts)
	if err != nil {
		return nil, false, err
	}
	return createdSts, false, nil
}

func (c *Client) GetStatefulSet(ns, name string) (*appsv1.StatefulSet, error) {
	sts, err := c.ClientSet.AppsV1().StatefulSets(ns).Get(name, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}
	sts.APIVersion = "apps/v1"
	sts.Kind = "StatefulSet"
	return sts, nil
}

func (c *Client) GetStatefulSetDeployment(ns, name string) (*StatefulSetDeploymentStatus, error) {

	sts, err := c.ClientSet.AppsV1().StatefulSets(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	labels := sts.Spec.Template.ObjectMeta.Labels
	svcList := []apiv1.Service{}
	vpcList := []apiv1.PersistentVolumeClaim{}

	svcs, _ := c.ClientSet.CoreV1().Services(ns).List(metav1.ListOptions{})
	if svcs != nil {
		for _, svc := range svcs.Items {
			for key, value := range svc.Spec.Selector {
				if labels[key] == value {
					svcList = append(svcList, svc)
					continue
				}
			}
		}
	}

	vpcs, _ := c.ClientSet.CoreV1().PersistentVolumeClaims(ns).List(metav1.ListOptions{})
	if vpcs != nil {
		for _, vpc := range vpcs.Items {
			for key, value := range vpc.Spec.Selector.MatchLabels {
				if labels[key] == value {
					vpcList = append(vpcList, vpc)
					continue
				}
			}
		}
	}

	dep := &StatefulSetDeploymentStatus{
		Namespace:         ns,
		Name:              name,
		Labels:            labels,
		StatefulSet:       sts,
		Services:          svcList,
		VolumesClaims:     vpcList,
		StatefulSetStatus: stsToStatus(sts),
		ServicesStatus:    svcsToStatus(svcList),
	}
	return dep, nil
}

func stsToStatus(sts *appsv1.StatefulSet) *StatefulSetStatus {
	stss := &StatefulSetStatus{
		Name:      sts.Name,
		Namespace: sts.Namespace,
		Desired:   *sts.Spec.Replicas,
		Running:   sts.Status.Replicas,
		Ready:     sts.Status.ReadyReplicas,
		Image:     "",
		Age:       time.Since(sts.CreationTimestamp.Time),
		PVC:       false,
	}
	for _, container := range sts.Spec.Template.Spec.Containers {
		stss.Image = container.Image
	}
	if len(sts.Spec.VolumeClaimTemplates) > 0 {
		stss.PVC = true
	}
	return stss
}

func (c *Client) GetStatefulSetEvents(ctx context.Context, evt chan *appsv1.StatefulSet, done chan struct{}) {
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(c.ClientSet, time.Second*5)
	stsInformer := kubeInformerFactory.Apps().V1().StatefulSets().Informer()
	stop := make(chan struct{})
	defer close(stop)
	stsInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			sts, ok := obj.(*appsv1.StatefulSet)
			if ok {
				evt <- sts
			}

		},
		DeleteFunc: func(obj interface{}) {
			sts, ok := obj.(*appsv1.StatefulSet)
			if ok {
				evt <- sts
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			sts, ok := newObj.(*appsv1.StatefulSet)
			if ok {
				evt <- sts
			}
		},
	})
	kubeInformerFactory.Start(stop)
	for {
		select {
		case <-done:
			return
		case <-ctx.Done():
			return

		}
	}
}

func (c *Client) Scale(ctx context.Context, ns, name string, replicas int32) error {
	currentScale, err := c.ClientSet.AppsV1().StatefulSets(ns).GetScale(name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	currentScale.Spec.Replicas = replicas
	_, err = c.ClientSet.AppsV1().StatefulSets(ns).UpdateScale(name, currentScale)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DescribeStatefulSet(ns, name string) (string, error) {

	sts, err := c.ClientSet.AppsV1().StatefulSets(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(sts)
	if err != nil {
		return "", err
	}
	y, err := yaml.JSONToYAML(data)
	if err != nil {
		return "", err
	}
	return string(y), nil
}

func (c *Client) DeleteVolumeClaimsForStatefulSet(name string) error {
	pair := strings.Split(name, "/")
	vpcs, err := c.ClientSet.CoreV1().PersistentVolumeClaims(pair[0]).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, vpc := range vpcs.Items {
		if strings.Contains(vpc.Name, pair[1]) {
			err := c.ClientSet.CoreV1().PersistentVolumeClaims(pair[0]).Delete(vpc.Name, metav1.NewDeleteOptions(0))
			if err != nil {
				utils.Printlnf("Persistence Volume Claim %s/%s not deleted. Error: %s", vpc.Namespace, vpc.Namespace, utils.Title(err.Error()))
				continue
			}
			utils.Printlnf("Persistence Volume Claim %s/%s deleted.", vpc.Namespace, vpc.Name)
		}
	}
	return nil
}

func (c *Client) PrintStatefulSetStatus(ctx context.Context, desired int32, namespace, name string) {
	stsDone := make(chan struct{})
	stsCh := make(chan *appsv1.StatefulSet)

	go c.GetStatefulSetEvents(ctx, stsCh, stsDone)
	last := ""
	for {
		select {
		case sts := <-stsCh:
			if sts.Name == name && sts.Namespace == namespace {
				current := fmt.Sprintf("[StatefulSet] -> Desired - %d Current - %d Ready - %d", desired, sts.Status.Replicas, sts.Status.ReadyReplicas)
				//utils.Printlnf("[StatefulSet] -> Desired - %d Current - %d Ready - %d", desired, sts.Status.Replicas, sts.Status.ReadyReplicas)
				if last != current {
					utils.Println(current)
					last = current
				}
			}
		case <-ctx.Done():
			stsDone <- struct{}{}
			return
		}
	}

}
