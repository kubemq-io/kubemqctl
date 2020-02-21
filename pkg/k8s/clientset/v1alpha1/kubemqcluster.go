package v1alpha1

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/crd/cluster"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type KubemqClusterInterface interface {
	List(opts metav1.ListOptions) (*cluster.KubemqClusterList, error)
	Get(name string, options metav1.GetOptions) (*cluster.KubemqCluster, error)
	Create(kubemqCluster *cluster.KubemqCluster) (*cluster.KubemqCluster, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Delete(name string, options metav1.DeleteOptions) error
}

type restClient struct {
	restClient rest.Interface
	ns         string
}

func (c *restClient) List(opts metav1.ListOptions) (*cluster.KubemqClusterList, error) {
	result := cluster.KubemqClusterList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("kubemqclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *restClient) Get(name string, opts metav1.GetOptions) (*cluster.KubemqCluster, error) {
	result := cluster.KubemqCluster{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("kubemqclusters").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *restClient) Create(kubemqCluster *cluster.KubemqCluster) (*cluster.KubemqCluster, error) {
	result := cluster.KubemqCluster{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("kubemqclusters").
		Body(kubemqCluster).
		Do().
		Into(&result)

	return &result, err
}
func (c *restClient) Delete(name string, opts metav1.DeleteOptions) error {

	err := c.restClient.
		Delete().
		Namespace(c.ns).
		Resource("kubemqclusters").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().Error()

	return err
}
func (c *restClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("kubemqclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
