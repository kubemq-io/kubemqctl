package v1beta1

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type KubemqClustersGetter interface {
	KubemqClusters(namespace string) KubemqClusterInterface
}

type V1Beta1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*V1Beta1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: types.GroupName, Version: types.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &V1Beta1Client{restClient: client}, nil
}

func (c *V1Beta1Client) KubemqClusters(namespace string) KubemqClusterInterface {
	return &kubemqCluster{
		client: c.restClient,
		ns:     namespace,
	}
}

func (c *V1Beta1Client) KubemqConnector(namespace string) KubemqConnectorInterface {
	return &kubemqConnector{
		client: c.restClient,
		ns:     namespace,
	}
}
