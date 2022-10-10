package v1alpha1

import (
	"context"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqdashboard"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"time"
)

type KubemqDashboardInterface interface {
	List(opts metav1.ListOptions) (*kubemqdashboard.KubemqDashboardList, error)
	Get(name string, options metav1.GetOptions) (*kubemqdashboard.KubemqDashboard, error)
	Create(cluster *kubemqdashboard.KubemqDashboard) (*kubemqdashboard.KubemqDashboard, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Update(cluster *kubemqdashboard.KubemqDashboard) (*kubemqdashboard.KubemqDashboard, error)
}

type kubemqDashboard struct {
	client rest.Interface
	ns     string
}

func (c *kubemqDashboard) List(opts metav1.ListOptions) (*kubemqdashboard.KubemqDashboardList, error) {
	result := kubemqdashboard.KubemqDashboardList{}
	err := c.client.
		Get().
		Namespace(c.ns).
		Resource("kubemqdashboards").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *kubemqDashboard) Get(clusterName string, opts metav1.GetOptions) (*kubemqdashboard.KubemqDashboard, error) {
	result := kubemqdashboard.KubemqDashboard{}
	err := c.client.
		Get().
		Namespace(c.ns).
		Resource("kubemqdashboards").
		Name(clusterName).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *kubemqDashboard) Create(cluster *kubemqdashboard.KubemqDashboard) (*kubemqdashboard.KubemqDashboard, error) {
	result := &kubemqdashboard.KubemqDashboard{}
	err := c.client.
		Post().
		Namespace(c.ns).
		Resource("kubemqdashboards").
		Body(cluster).
		Do(context.Background()).
		Into(result)

	return result, err
}

func (c *kubemqDashboard) Update(cluster *kubemqdashboard.KubemqDashboard) (*kubemqdashboard.KubemqDashboard, error) {
	result := &kubemqdashboard.KubemqDashboard{}
	err := c.client.
		Put().
		Namespace(c.ns).
		Name(cluster.Name).
		Resource("kubemqdashboards").
		Body(cluster).
		Do(context.Background()).
		Into(result)

	return result, err
}
func (c *kubemqDashboard) Delete(name string, opts *metav1.DeleteOptions) error {
	err := c.client.Delete().
		Namespace(c.ns).
		Resource("kubemqdashboards").
		Name(name).
		Body(opts).
		Do(context.Background()).
		Error()

	return err
}
func (c *kubemqDashboard) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	return c.client.
		Get().
		Namespace(c.ns).
		Resource("kubemqdashboards").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(context.Background())
}
