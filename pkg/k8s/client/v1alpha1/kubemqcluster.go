package v1alpha1

import (
	"context"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"time"
)

type KubemqClusterInterface interface {
	List(opts metav1.ListOptions) (*kubemqcluster.KubemqClusterList, error)
	Get(name string, options metav1.GetOptions) (*kubemqcluster.KubemqCluster, error)
	Create(cluster *kubemqcluster.KubemqCluster) (*kubemqcluster.KubemqCluster, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Update(cluster *kubemqcluster.KubemqCluster) (*kubemqcluster.KubemqCluster, error)
	GetScale(deploymentName string, options metav1.GetOptions) (*autoscalingv1.Scale, error)
	UpdateScale(deploymentName string, scale *autoscalingv1.Scale) (*autoscalingv1.Scale, error)
}

type kubemqCluster struct {
	client rest.Interface
	ns     string
}

func (c *kubemqCluster) List(opts metav1.ListOptions) (*kubemqcluster.KubemqClusterList, error) {
	result := kubemqcluster.KubemqClusterList{}
	err := c.client.
		Get().
		Namespace(c.ns).
		Resource("kubemqclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *kubemqCluster) Get(clusterName string, opts metav1.GetOptions) (*kubemqcluster.KubemqCluster, error) {
	result := kubemqcluster.KubemqCluster{}
	err := c.client.
		Get().
		Namespace(c.ns).
		Resource("kubemqclusters").
		Name(clusterName).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *kubemqCluster) Create(cluster *kubemqcluster.KubemqCluster) (*kubemqcluster.KubemqCluster, error) {
	result := &kubemqcluster.KubemqCluster{}
	err := c.client.
		Post().
		Namespace(c.ns).
		Resource("kubemqclusters").
		Body(cluster).
		Do(context.Background()).
		Into(result)

	return result, err
}

func (c *kubemqCluster) Update(cluster *kubemqcluster.KubemqCluster) (*kubemqcluster.KubemqCluster, error) {
	result := &kubemqcluster.KubemqCluster{}
	err := c.client.
		Put().
		Namespace(c.ns).
		Name(cluster.Name).
		Resource("kubemqclusters").
		Body(cluster).
		Do(context.Background()).
		Into(result)

	return result, err
}
func (c *kubemqCluster) Delete(name string, opts *metav1.DeleteOptions) error {
	err := c.client.Delete().
		Namespace(c.ns).
		Resource("kubemqclusters").
		Name(name).
		Body(opts).
		Do(context.Background()).
		Error()

	return err
}
func (c *kubemqCluster) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	return c.client.
		Get().
		Namespace(c.ns).
		Resource("kubemqclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(context.Background())
}

func (c *kubemqCluster) GetScale(clusterName string, options metav1.GetOptions) (result *autoscalingv1.Scale, err error) {
	result = &autoscalingv1.Scale{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("kubemqclusters").
		Name(clusterName).
		SubResource("scale").
		VersionedParams(&options, scheme.ParameterCodec).
		Do(context.Background()).
		Into(result)
	return
}

// UpdateScale takes the top resource name and the representation of a scale and updates it. Returns the server's representation of the scale, and an error, if there is any.
func (c *kubemqCluster) UpdateScale(clusterName string, scale *autoscalingv1.Scale) (result *autoscalingv1.Scale, err error) {
	result = &autoscalingv1.Scale{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("kubemqclusters").
		Name(clusterName).
		SubResource("scale").
		Body(scale).
		Do(context.Background()).
		Into(result)
	return
}
