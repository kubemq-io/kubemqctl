package v1alpha1

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/crd"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type V1Alpha1Interface interface {
	KubemqClusters(namespace string) KubemqClusterInterface
}

type V1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*V1Alpha1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: crd.GroupName, Version: crd.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &V1Alpha1Client{restClient: client}, nil
}

func (c *V1Alpha1Client) KubemqClusters(namespace string) KubemqClusterInterface {
	return &restClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
