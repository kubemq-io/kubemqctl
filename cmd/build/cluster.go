package build

import (
	clusterbuilder "github.com/kubemq-hub/builder/cluster"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClustersBuilder struct {
	deployments []*kubemqcluster.KubemqCluster
	resources   *resources
}

func newClustersBuilder() *ClustersBuilder {
	return &ClustersBuilder{}
}
func (c *ClustersBuilder) SetResources(value *resources) *ClustersBuilder {
	c.resources = value
	return c
}

func (c *ClustersBuilder) render() error {
	cluster, err := clusterbuilder.
		NewCluster().
		Render()
	if err != nil {
		return err
	}
	if err := c.resources.updateClusterName(cluster.Namespace, cluster.Name); err != nil {
		return err
	}

	deployment := &kubemqcluster.KubemqCluster{
		TypeMeta: v1.TypeMeta{
			Kind:       "KubemqCluster",
			APIVersion: "core.k8s.kubemq.io/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      cluster.Name,
			Namespace: cluster.Name,
		},
		Spec: kubemqcluster.KubemqClusterSpec{
			Replicas:       new(int32),
			License:        "",
			ConfigData:     "",
			Volume:         nil,
			Image:          nil,
			Api:            nil,
			Rest:           nil,
			Grpc:           nil,
			Tls:            nil,
			Resources:      nil,
			NodeSelectors:  nil,
			Authentication: nil,
			Authorization:  nil,
			Health:         nil,
			Routing:        nil,
			Log:            nil,
			Notification:   nil,
			Store:          nil,
			Queue:          nil,
		},
		Status: kubemqcluster.KubemqClusterStatus{},
	}

	*deployment.Spec.Replicas = int32(cluster.Replicas)
	deployment.Spec.License = cluster.License

	if cluster.Volume != nil {
		deployment.Spec.Volume = &kubemqcluster.VolumeConfig{
			Size:         cluster.Volume.Size,
			StorageClass: "",
		}
	}

	if cluster.Image != nil {
		deployment.Spec.Image = &kubemqcluster.ImageConfig{
			Image:      cluster.Image.Image,
			PullPolicy: cluster.Image.PullPolicy,
		}
	}

	if cluster.Api != nil {
		deployment.Spec.Api = &kubemqcluster.ApiConfig{
			Disabled: false,
			Port:     8080,
			Expose:   cluster.Api.Expose,
			NodePort: int32(cluster.Api.NodePort),
		}
	}

	if cluster.Grpc != nil {
		deployment.Spec.Grpc = &kubemqcluster.GrpcConfig{
			Disabled:   false,
			Port:       50000,
			Expose:     cluster.Grpc.Expose,
			NodePort:   int32(cluster.Grpc.NodePort),
			BufferSize: int32(cluster.Grpc.BufferSize),
			BodyLimit:  int32(cluster.Grpc.BodyLimit),
		}
	}

	if cluster.Rest != nil {
		deployment.Spec.Rest = &kubemqcluster.RestConfig{
			Disabled:   false,
			Port:       9090,
			Expose:     cluster.Rest.Expose,
			NodePort:   int32(cluster.Rest.NodePort),
			BufferSize: int32(cluster.Rest.BufferSize),
			BodyLimit:  int32(cluster.Rest.BodyLimit),
		}
	}
	if cluster.Tls != nil {
		deployment.Spec.Tls = &kubemqcluster.TlsConfig{
			Cert: cluster.Tls.Cert,
			Key:  cluster.Tls.Key,
			Ca:   cluster.Tls.Ca,
		}
	}

	if cluster.Resource != nil {
		deployment.Spec.Resources = &kubemqcluster.ResourceConfig{
			LimitsCpu:      cluster.Resource.LimitsCpu,
			LimitsMemory:   cluster.Resource.LimitsMemory,
			RequestsCpu:    cluster.Resource.RequestsCpu,
			RequestsMemory: cluster.Resource.RequestsMemory,
		}
	}

	if cluster.NodeSelectors != nil {
		deployment.Spec.NodeSelectors = &kubemqcluster.NodeSelectorConfig{
			Keys: cluster.NodeSelectors,
		}
	}
	if cluster.Authentication != nil {
		deployment.Spec.Authentication = &kubemqcluster.AuthenticationConfig{
			Key:  cluster.Authentication.Key,
			Type: cluster.Authentication.Type,
		}
	}

	if cluster.Authorization != nil {
		deployment.Spec.Authorization = &kubemqcluster.AuthorizationConfig{
			Policy:     cluster.Authorization.Policy,
			Url:        cluster.Authorization.Url,
			AutoReload: int32(cluster.Authorization.AutoReload),
		}
	}
	if cluster.Health != nil {
		deployment.Spec.Health = &kubemqcluster.HealthConfig{
			Enabled:             cluster.Health.Enabled,
			InitialDelaySeconds: int32(cluster.Health.InitialDelaySeconds),
			PeriodSeconds:       int32(cluster.Health.PeriodSeconds),
			TimeoutSeconds:      int32(cluster.Health.TimeoutSeconds),
			SuccessThreshold:    int32(cluster.Health.SuccessThreshold),
			FailureThreshold:    int32(cluster.Health.FailureThreshold),
		}
	}
	if cluster.Routing != nil {
		deployment.Spec.Routing = &kubemqcluster.RoutingConfig{
			Data:       cluster.Routing.Data,
			Url:        cluster.Routing.Url,
			AutoReload: int32(cluster.Routing.AutoReload),
		}
	}
	if cluster.Log != nil {
		deployment.Spec.Log = &kubemqcluster.LogConfig{
			Level: new(int32),
			File:  "",
		}
		*deployment.Spec.Log.Level = int32(cluster.Log.Level)
	}
	if cluster.Notification != nil {
		deployment.Spec.Notification = &kubemqcluster.NotificationConfig{
			Enabled: cluster.Notification.Enabled,
			Prefix:  cluster.Notification.Prefix,
			Log:     false,
		}
	}
	if cluster.Store != nil {
		deployment.Spec.Store = &kubemqcluster.StoreConfig{
			Clean:                    cluster.Store.Clean,
			Path:                     cluster.Store.Path,
			MaxChannels:              new(int32),
			MaxSubscribers:           new(int32),
			MaxMessages:              new(int32),
			MaxChannelSize:           new(int32),
			MessagesRetentionMinutes: new(int32),
			PurgeInactiveMinutes:     new(int32),
		}
		*deployment.Spec.Store.MaxChannels = int32(cluster.Store.MaxChannels)
		*deployment.Spec.Store.MaxSubscribers = int32(cluster.Store.MaxSubscribers)
		*deployment.Spec.Store.MaxMessages = int32(cluster.Store.MaxMessages)
		*deployment.Spec.Store.MaxChannelSize = int32(cluster.Store.MaxChannelSize)
		*deployment.Spec.Store.MessagesRetentionMinutes = int32(cluster.Store.MessagesRetentionMinutes)
		*deployment.Spec.Store.PurgeInactiveMinutes = int32(cluster.Store.PurgeInactiveMinutes)
	}
	if cluster.Queue != nil {
		deployment.Spec.Queue = &kubemqcluster.QueueConfig{
			MaxReceiveMessagesRequest: new(int32),
			MaxWaitTimeoutSeconds:     new(int32),
			MaxExpirationSeconds:      new(int32),
			MaxDelaySeconds:           new(int32),
			MaxReQueues:               new(int32),
			MaxVisibilitySeconds:      new(int32),
			DefaultVisibilitySeconds:  new(int32),
			DefaultWaitTimeoutSeconds: new(int32),
		}
		*deployment.Spec.Queue.MaxReceiveMessagesRequest = int32(cluster.Queue.MaxReceiveMessagesRequest)
		*deployment.Spec.Queue.MaxWaitTimeoutSeconds = int32(cluster.Queue.MaxWaitTimeoutSeconds)
		*deployment.Spec.Queue.MaxExpirationSeconds = int32(cluster.Queue.MaxExpirationSeconds)
		*deployment.Spec.Queue.MaxDelaySeconds = int32(cluster.Queue.MaxDelaySeconds)
		*deployment.Spec.Queue.MaxReQueues = int32(cluster.Queue.MaxReQueues)
		*deployment.Spec.Queue.MaxVisibilitySeconds = int32(cluster.Queue.MaxVisibilitySeconds)
		*deployment.Spec.Queue.DefaultVisibilitySeconds = int32(cluster.Queue.DefaultVisibilitySeconds)
		*deployment.Spec.Queue.DefaultWaitTimeoutSeconds = int32(cluster.Queue.DefaultWaitTimeoutSeconds)
	}
	c.deployments = append(c.deployments, deployment)
	return nil
}
func (c *ClustersBuilder) add() error {
	for {
		utils.Println("Adding new KubeMQ Cluster:")
		if err := c.render(); err != nil {
			utils.Printlnf("Error adding new KubeMQ Cluster: %s", err.Error())
		} else {
			utils.Println("Adding new KubeMQ Cluster completed successfully!")
			return nil
		}

	}

}
