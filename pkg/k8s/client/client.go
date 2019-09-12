package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kubemq-io/kubetools/pkg/utils"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"sort"
	"strings"

	"time"
)

type Client struct {
	ClientSet    *kubernetes.Clientset
	ClientConfig clientcmd.ClientConfig
}

func NewClient(kubeConfigPath string) (*Client, error) {

	var kubeconfig string
	if kubeConfigPath != "" {
		kubeconfig = kubeConfigPath
	} else {
		if home := homeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		} else {
			return nil, errors.New("no kubeconfig available")
		}
	}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: ""}})
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("kubernetes config file: %s", err)
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	c := &Client{
		ClientSet:    clientset,
		ClientConfig: clientConfig,
	}
	return c, nil
}

func (c *Client) GetConfigClusters() (map[string]*clientcmdapi.Cluster, error) {
	config, err := c.ClientConfig.ConfigAccess().GetStartingConfig()
	if err != nil {
		return nil, err
	}
	return config.Clusters, nil
}

func (c *Client) GetConfigContext() (map[string]*clientcmdapi.Context, string, error) {
	config, err := c.ClientConfig.ConfigAccess().GetStartingConfig()
	if err != nil {
		return nil, "", err
	}
	return config.Contexts, config.CurrentContext, nil
}

func (c *Client) GetCurrentContext() (string, error) {
	config, err := c.ClientConfig.ConfigAccess().GetStartingConfig()
	if err != nil {
		return "", err
	}
	return config.CurrentContext, nil
}

func (c *Client) SwitchContext(contextName string) error {
	config, err := c.ClientConfig.ConfigAccess().GetStartingConfig()
	if err != nil {
		return err
	}
	config.CurrentContext = contextName
	err = clientcmd.ModifyConfig(c.ClientConfig.ConfigAccess(), *config, true)
	return err
}
func (c *Client) GetNamespace(name string) (*apiv1.Namespace, bool, error) {
	ns, err := c.ClientSet.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	if err == nil && ns != nil {
		return ns, true, nil
	}
	newNs := &apiv1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec:   apiv1.NamespaceSpec{},
		Status: apiv1.NamespaceStatus{},
	}
	return newNs, false, nil

}
func (c *Client) CheckAndCreateNamespace(name string) (*apiv1.Namespace, bool, error) {
	ns, err := c.ClientSet.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	if err == nil && ns != nil {
		return ns, false, nil
	}
	newNs := &apiv1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec:   apiv1.NamespaceSpec{},
		Status: apiv1.NamespaceStatus{},
	}
	ns, err = c.ClientSet.CoreV1().Namespaces().Create(newNs)

	if err != nil {
		return nil, false, err
	}

	return ns, true, nil
}

func (c *Client) CreateOrUpdateStatefulSet(sts *appsv1.StatefulSet) (*appsv1.StatefulSet, bool, error) {

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
func (c *Client) DeleteStatefulSet(name string) error {
	pair := strings.Split(name, "/")
	return c.ClientSet.AppsV1().StatefulSets(pair[0]).Delete(pair[1], metav1.NewDeleteOptions(0))
}
func (c *Client) DeleteServicesForStatefulSet(name string) error {
	pair := strings.Split(name, "/")
	svcs, err := c.ClientSet.CoreV1().Services(pair[0]).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, svc := range svcs.Items {
		if strings.Contains(svc.Name, pair[1]) {
			err := c.ClientSet.CoreV1().Services(pair[0]).Delete(svc.Name, metav1.NewDeleteOptions(0))
			if err != nil {
				utils.Printlnf("Service %s/%s not deleted. Error: %s", svc.Namespace, svc.Namespace, utils.Title(err.Error()))
				continue
			}
			utils.Printlnf("Service %s/%s deleted.", svc.Namespace, svc.Name)
		}
	}
	return nil
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

func (c *Client) CreateOrUpdateService(svc *apiv1.Service) (*apiv1.Service, bool, error) {

	oldSvc, err := c.ClientSet.CoreV1().Services(svc.Namespace).Get(svc.Name, metav1.GetOptions{})
	if err == nil && oldSvc != nil {
		svc.ResourceVersion = oldSvc.ResourceVersion
		svc.Spec.ClusterIP = oldSvc.Spec.ClusterIP
		newSvc, err := c.ClientSet.CoreV1().Services(svc.Namespace).Update(svc)
		if err != nil {
			return nil, true, err
		}
		return newSvc, true, nil
	}

	createdSvc, err := c.ClientSet.CoreV1().Services(svc.Namespace).Create(svc)
	if err != nil {
		return nil, false, err
	}
	return createdSvc, false, nil
}

func (c *Client) CreateOrUpdateConfigMap(cm *apiv1.ConfigMap) (*apiv1.ConfigMap, bool, error) {
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
func (c *Client) CreateOrUpdateSecret(sec *apiv1.Secret) (*apiv1.Secret, bool, error) {
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

func (c *Client) GetConfigMap(namespace, name string) (*apiv1.ConfigMap, error) {

	cm, err := c.ClientSet.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return cm, nil
}

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

func (c *Client) DescribeStatefulSet(ns, name string) (string, error) {

	sts, err := c.ClientSet.AppsV1().StatefulSets(ns).Get(name, metav1.GetOptions{})
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

func (c *Client) GetServices(ns string, contains ...string) (map[string]apiv1.Service, error) {
	svcs, err := c.ClientSet.CoreV1().Services(ns).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	list := map[string]apiv1.Service{}
	for _, item := range svcs.Items {
		if contains == nil {
			list[fmt.Sprintf("%s/%s", item.Namespace, item.Name)] = item
		} else {
			for _, name := range contains {
				if strings.Contains(item.Name, name) {
					list[fmt.Sprintf("%s/%s", item.Namespace, item.Name)] = item
					continue
				}
			}
		}

	}
	return list, err
}

func (c *Client) GetPods(ns string, contains ...string) (map[string]apiv1.Pod, error) {
	pods, err := c.ClientSet.CoreV1().Pods(ns).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	list := map[string]apiv1.Pod{}
	for _, item := range pods.Items {
		if contains == nil {
			list[fmt.Sprintf("%s/%s", item.Namespace, item.Name)] = item
		} else {

			for _, name := range contains {
				if strings.Contains(item.Name, name) {
					list[fmt.Sprintf("%s/%s", item.Namespace, item.Name)] = item
					continue
				}
			}
		}

	}
	return list, err
}

func (c *Client) ForwardPorts(ns string, name string, ports []string, stopChan chan struct{}, outCh chan string, errOutCh chan string) error {
	restConfig, err := c.ClientConfig.ClientConfig()
	if err != nil {
		return err
	}
	roundTripper, upgrader, err := spdy.RoundTripperFor(restConfig)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", ns, name)
	hostIP := strings.TrimLeft(restConfig.Host, "https:/")
	serverURL := url.URL{Scheme: "https", Path: path, Host: hostIP}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodPost, &serverURL)
	readyChan := make(chan struct{}, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)
	forwarder, err := portforward.New(dialer, ports, stopChan, readyChan, out, errOut)
	if err != nil {
		return err
	}

	go func() {
		for range readyChan { // Kubernetes will close this channel when it has something to tell us.
		}
		if len(errOut.String()) != 0 {
			errOutCh <- errOut.String()
			close(stopChan)
		} else if len(out.String()) != 0 {

			outCh <- out.String()
		}

	}()

	go func() {
		if err = forwarder.ForwardPorts(); err != nil { // Locks until stopChan is closed.
			errOutCh <- err.Error()

		}
	}()

	return nil
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
	done := make(chan struct{})
	evt := make(chan *appsv1.StatefulSet)
	go c.GetStatefulSetEvents(ctx, evt, done)

	for {
		select {
		case sts := <-evt:
			if replicas == sts.Status.Replicas && sts.Status.Replicas == sts.Status.ReadyReplicas {
				utils.Printlnf("Desired:%d Current:%d Ready:%d", replicas, sts.Status.Replicas, sts.Status.ReadyReplicas)
				done <- struct{}{}
				return nil
			} else {
				utils.Printlnf("Desired:%d Current:%d Ready:%d", replicas, sts.Status.Replicas, sts.Status.ReadyReplicas)
			}
		case <-ctx.Done():
			return nil
		}

	}

}

func (c *Client) GetKubeMQClusters() ([]string, error) {

	sets, err := c.GetStatefulSets("")
	if err != nil {
		return nil, err
	}
	var list []string
	for key, set := range sets {
		for _, container := range set.Spec.Template.Spec.Containers {
			if strings.Contains(container.Image, "kubemq") {
				list = append(list, key)
				continue
			}
		}

	}
	sort.Strings(list)
	return list, nil
}

func (c *Client) GetKubeMQClustersStatus() ([]*StatefulSetStatus, error) {

	sets, err := c.GetStatefulSets("")
	if err != nil {
		return nil, err
	}
	var list []*StatefulSetStatus
	for _, set := range sets {
		for _, container := range set.Spec.Template.Spec.Containers {
			if strings.Contains(container.Image, "kubemq") {
				sts := &StatefulSetStatus{
					Name:      set.Name,
					Namespace: set.Namespace,
					Desired:   *set.Spec.Replicas,
					Running:   set.Status.Replicas,
					Ready:     set.Status.ReadyReplicas,
					Image:     container.Image,
					Age:       time.Now().Sub(set.CreationTimestamp.Time),
				}
				list = append(list, sts)
				continue
			}
		}
	}

	return list, nil
}

func (c *Client) GetStatefulSetDeployment(ns, name string) (*StatefulSetDeployment, error) {

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

	dep := &StatefulSetDeployment{
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
		Age:       time.Now().Sub(sts.CreationTimestamp.Time),
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

func svcsToStatus(svcs []apiv1.Service) []*ServiceStatus {
	list := []*ServiceStatus{}
	for _, svc := range svcs {
		ss := &ServiceStatus{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			Type:      string(svc.Spec.Type),
			ClusterIP: svc.Spec.ClusterIP,
			ExternalP: "",
			Ports:     "",
			Age:       time.Now().Sub(svc.CreationTimestamp.Time),
		}

		if string(svc.Spec.Type) == "LoadBalancer" {
			if len(svc.Status.LoadBalancer.Ingress) > 0 {
				ss.ExternalP = svc.Status.LoadBalancer.Ingress[0].IP
			}
		}
		portList := []string{}
		for _, port := range svc.Spec.Ports {
			if port.NodePort > 0 {
				portList = append(portList, fmt.Sprintf("%d:%d", port.Port, port.NodePort))
			} else {
				portList = append(portList, fmt.Sprintf("%d", port.Port))
			}
		}
		ss.Ports = strings.Join(portList, ",")

		list = append(list, ss)
	}
	return list
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func (c *Client) GetStatefulSetEvents(ctx context.Context, evt chan *appsv1.StatefulSet, done chan struct{}) error {
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
			return nil
		case <-ctx.Done():
			return nil

		}
	}
}
