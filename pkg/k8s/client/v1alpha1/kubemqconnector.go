package v1alpha1

import (
	"context"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqconnector"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"time"
)

type KubemqConnectorInterface interface {
	List(opts metav1.ListOptions) (*kubemqconnector.KubemqConnectorList, error)
	Get(name string, options metav1.GetOptions) (*kubemqconnector.KubemqConnector, error)
	Create(cluster *kubemqconnector.KubemqConnector) (*kubemqconnector.KubemqConnector, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Update(cluster *kubemqconnector.KubemqConnector) (*kubemqconnector.KubemqConnector, error)
}

type kubemqConnector struct {
	client rest.Interface
	ns     string
}

func (c *kubemqConnector) List(opts metav1.ListOptions) (*kubemqconnector.KubemqConnectorList, error) {
	result := kubemqconnector.KubemqConnectorList{}
	err := c.client.
		Get().
		Namespace(c.ns).
		Resource("kubemqconnectors").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *kubemqConnector) Get(clusterName string, opts metav1.GetOptions) (*kubemqconnector.KubemqConnector, error) {
	result := kubemqconnector.KubemqConnector{}
	err := c.client.
		Get().
		Namespace(c.ns).
		Resource("kubemqconnectors").
		Name(clusterName).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *kubemqConnector) Create(cluster *kubemqconnector.KubemqConnector) (*kubemqconnector.KubemqConnector, error) {
	result := &kubemqconnector.KubemqConnector{}
	err := c.client.
		Post().
		Namespace(c.ns).
		Resource("kubemqconnectors").
		Body(cluster).
		Do(context.Background()).
		Into(result)

	return result, err
}

func (c *kubemqConnector) Update(cluster *kubemqconnector.KubemqConnector) (*kubemqconnector.KubemqConnector, error) {
	result := &kubemqconnector.KubemqConnector{}
	err := c.client.
		Put().
		Namespace(c.ns).
		Name(cluster.Name).
		Resource("kubemqconnectors").
		Body(cluster).
		Do(context.Background()).
		Into(result)

	return result, err
}
func (c *kubemqConnector) Delete(name string, opts *metav1.DeleteOptions) error {
	err := c.client.Delete().
		Namespace(c.ns).
		Resource("kubemqconnectors").
		Name(name).
		Body(opts).
		Do(context.Background()).
		Error()

	return err
}
func (c *kubemqConnector) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	return c.client.
		Get().
		Namespace(c.ns).
		Resource("kubemqconnectors").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(context.Background())
}
