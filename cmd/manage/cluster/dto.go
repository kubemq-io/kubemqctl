package cluster

import (
	builder "github.com/kubemq-hub/builder/cluster"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/cluster"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ToDeployment(cluster *builder.Cluster) *kubemqcluster.KubemqCluster {
	deployment := &kubemqcluster.KubemqCluster{
		TypeMeta: v1.TypeMeta{
			Kind:       "KubemqCluster",
			APIVersion: "core.k8s.kubemq.io/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      cluster.Name,
			Namespace: cluster.Namespace,
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
	return deployment
}

func FromDeployment(deployment *kubemqcluster.KubemqCluster) *builder.Cluster {
	cluster := &builder.Cluster{
		Name:           deployment.Name,
		Namespace:      deployment.Namespace,
		Replicas:       0,
		Authentication: nil,
		Authorization:  nil,
		Health:         nil,
		Image:          nil,
		License:        deployment.Spec.License,
		Log:            nil,
		NodeSelectors:  nil,
		Notification:   nil,
		Queue:          nil,
		Resource:       nil,
		Api:            nil,
		Grpc:           nil,
		Rest:           nil,
		Routing:        nil,
		Store:          nil,
		Tls:            nil,
		Volume:         nil,
		Status:         nil,
	}
	spec := deployment.Spec
	if spec.Replicas != nil {
		cluster.Replicas = int(*spec.Replicas)
	}

	if spec.Authentication != nil {
		cluster.Authentication = &builder.Authentication{
			Key:  spec.Authentication.Key,
			Type: spec.Authentication.Type,
		}
	}

	if spec.Authorization != nil {
		cluster.Authorization = &builder.Authorization{
			Policy:     spec.Authorization.Policy,
			Url:        spec.Authorization.Url,
			AutoReload: int(spec.Authorization.AutoReload),
		}
	}

	if spec.Health != nil {
		cluster.Health = &builder.Health{
			Enabled:             spec.Health.Enabled,
			InitialDelaySeconds: int(spec.Health.InitialDelaySeconds),
			PeriodSeconds:       int(spec.Health.PeriodSeconds),
			TimeoutSeconds:      int(spec.Health.TimeoutSeconds),
			SuccessThreshold:    int(spec.Health.SuccessThreshold),
			FailureThreshold:    int(spec.Health.FailureThreshold),
		}
	}
	if spec.Image != nil {
		cluster.Image = &builder.Image{
			Image:      spec.Image.Image,
			PullPolicy: spec.Image.PullPolicy,
		}
	}

	if spec.Log != nil {
		cluster.Log = &builder.Log{
			Level: int(*spec.Log.Level),
		}
	}
	if spec.NodeSelectors != nil {
		cluster.NodeSelectors = map[string]string{}
		for key, val := range spec.NodeSelectors.Keys {
			cluster.NodeSelectors[key] = val
		}
	}
	if spec.Notification != nil {
		cluster.Notification = &builder.Notification{
			Prefix:  spec.Notification.Prefix,
			Enabled: spec.Notification.Enabled,
		}
	}
	if spec.Queue != nil {
		cluster.Queue = &builder.Queue{
			MaxReceiveMessagesRequest: int(*spec.Queue.MaxReceiveMessagesRequest),
			MaxWaitTimeoutSeconds:     int(*spec.Queue.MaxWaitTimeoutSeconds),
			MaxExpirationSeconds:      int(*spec.Queue.MaxExpirationSeconds),
			MaxDelaySeconds:           int(*spec.Queue.MaxDelaySeconds),
			MaxReQueues:               int(*spec.Queue.MaxReQueues),
			MaxVisibilitySeconds:      int(*spec.Queue.MaxVisibilitySeconds),
			DefaultVisibilitySeconds:  int(*spec.Queue.DefaultVisibilitySeconds),
			DefaultWaitTimeoutSeconds: int(*spec.Queue.DefaultWaitTimeoutSeconds),
		}
	}
	if spec.Resources != nil {
		cluster.Resource = &builder.Resource{
			LimitsCpu:      spec.Resources.LimitsCpu,
			LimitsMemory:   spec.Resources.LimitsMemory,
			RequestsCpu:    spec.Resources.RequestsCpu,
			RequestsMemory: spec.Resources.RequestsMemory,
		}
	}
	if spec.Api != nil {
		cluster.Api = &builder.Service{
			NodePort:   int(spec.Api.NodePort),
			Expose:     spec.Api.Expose,
			BufferSize: 0,
			BodyLimit:  0,
		}
	}
	if spec.Grpc != nil {
		cluster.Grpc = &builder.Service{
			NodePort:   int(spec.Grpc.NodePort),
			Expose:     spec.Grpc.Expose,
			BufferSize: int(spec.Grpc.BufferSize),
			BodyLimit:  int(spec.Grpc.BodyLimit),
		}
	}
	if spec.Rest != nil {
		cluster.Rest = &builder.Service{
			NodePort:   int(spec.Rest.NodePort),
			Expose:     spec.Rest.Expose,
			BufferSize: int(spec.Rest.BufferSize),
			BodyLimit:  int(spec.Rest.BodyLimit),
		}
	}
	if spec.Routing != nil {
		cluster.Routing = &builder.Routing{
			Data:       spec.Routing.Data,
			Url:        spec.Routing.Url,
			AutoReload: int(spec.Routing.AutoReload),
		}
	}
	if spec.Store != nil {
		cluster.Store = &builder.Store{
			Clean:                    spec.Store.Clean,
			Path:                     spec.Store.Path,
			MaxChannels:              int(*spec.Store.MaxChannels),
			MaxSubscribers:           int(*spec.Store.MaxSubscribers),
			MaxMessages:              int(*spec.Store.MaxMessages),
			MaxChannelSize:           int(*spec.Store.MaxChannelSize),
			MessagesRetentionMinutes: int(*spec.Store.MessagesRetentionMinutes),
			PurgeInactiveMinutes:     int(*spec.Store.PurgeInactiveMinutes),
		}
	}
	if spec.Tls != nil {
		cluster.Tls = &builder.Tls{
			Cert: spec.Tls.Cert,
			Key:  spec.Tls.Key,
			Ca:   spec.Tls.Ca,
		}
	}
	if spec.Volume != nil {
		cluster.Volume = &builder.Volume{
			Size:         spec.Volume.Size,
			StorageClass: spec.Volume.StorageClass,
		}
	}

	status := deployment.Status
	cluster.Status = &builder.Status{

		Version:       status.Version,
		Ready:         status.Ready,
		Grpc:          status.Grpc,
		Rest:          status.Rest,
		Api:           status.Api,
		Selector:      status.Selector,
		LicenseType:   status.LicenseType,
		LicenseTo:     status.LicenseTo,
		LicenseExpire: status.LicenseExpire,
		Status:        status.Status,
	}
	if status.Replicas != nil {
		cluster.Status.Replicas = *status.Replicas
	}
	return cluster
}

func FromDeploymentList(list *cluster.KubemqClusters) []*builder.Cluster {
	var clustersList []*builder.Cluster
	for _, deployment := range list.Items() {
		clustersList = append(clustersList, FromDeployment(deployment))
	}
	return clustersList
}
