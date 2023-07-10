package builder

import (
	"encoding/json"
	"fmt"
	"github.com/kubemq-io/k8s/api/v1beta1"
	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/config"
	"strings"
)

type Cluster struct {
	Kind       string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	Metadata   struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
	Spec *v1beta1.KubemqClusterSpec `json:"spec"`
}

func NewCluster() *Cluster {
	c := &Cluster{
		Kind:       "KubemqCluster",
		ApiVersion: "core.k8s.kubemq.io/v1beta1",
		Spec: &v1beta1.KubemqClusterSpec{
			Replicas:              new(int32),
			ConfigData:            "",
			License:               "",
			Key:                   "",
			Standalone:            false,
			Volume:                nil,
			Image:                 nil,
			Api:                   nil,
			Rest:                  nil,
			Grpc:                  nil,
			Tls:                   nil,
			Resources:             nil,
			NodeSelectors:         nil,
			Authentication:        nil,
			Authorization:         nil,
			Health:                nil,
			Routing:               nil,
			Log:                   nil,
			Notification:          nil,
			Store:                 nil,
			Queue:                 nil,
			StatefulSetConfigData: "",
		},
	}
	return c
}

func (c *Cluster) FromDTO(dto *ClusterDTO) *Cluster {
	c.SetBasic(dto).
		SetGRPC(dto).
		SetREST(dto).
		SetAPI(dto).
		SetSecurity(dto).
		SetAuthentication(dto).
		SetAuthorization(dto).
		SetRouting(dto).
		SetVolume(dto).
		SetImage(dto).
		SetHealth(dto).
		SetResources(dto).
		SetNodes(dto).
		SetStore(dto).
		SetQueues(dto)
	return c
}
func (c *Cluster) Key() string {
	return c.Spec.Key
}

func (c *Cluster) SetBasic(dto *ClusterDTO) *Cluster {
	c.Metadata.Name = dto.Deployment.Base.ClusterName
	c.Metadata.Namespace = dto.Deployment.Base.ClusterNamespace
	c.Spec.Key = dto.Deployment.Base.LicenseKey

	switch dto.Deployment.Spec.Replicas {
	case "1 Node":
		*c.Spec.Replicas = 1
	case "3 Nodes":
		*c.Spec.Replicas = 3
	case "5 Nodes":
		*c.Spec.Replicas = 5
	case "7 Nodes":
		*c.Spec.Replicas = 7
	default:
		*c.Spec.Replicas = 3
	}

	if dto.Deployment.Spec.Mode == "Standalone" {
		c.Spec.Standalone = true
	}
	return c
}

func (c *Cluster) SetGRPC(dto *ClusterDTO) *Cluster {
	if dto.GrpcInterface == nil {
		return c
	}
	c.Spec.Grpc = &config.GrpcConfig{
		Expose:   dto.GrpcInterface.Mode,
		NodePort: int32(dto.GrpcInterface.NodePort),
	}
	return c
}
func (c *Cluster) SetREST(dto *ClusterDTO) *Cluster {
	if dto.RestInterface == nil {
		return c
	}

	c.Spec.Rest = &config.RestConfig{
		Expose:   dto.RestInterface.Mode,
		NodePort: int32(dto.RestInterface.NodePort),
	}
	return c
}
func (c *Cluster) SetAPI(dto *ClusterDTO) *Cluster {
	if dto.APIInterface == nil {
		return c
	}

	c.Spec.Api = &config.ApiConfig{
		Expose:   dto.APIInterface.Mode,
		NodePort: int32(dto.APIInterface.NodePort),
	}
	return c
}

func (c *Cluster) SetSecurity(dto *ClusterDTO) *Cluster {
	if dto.Security == nil {
		return c
	}
	c.Spec.Tls = &config.TlsConfig{
		Cert: dto.Security.Cert,
		Key:  dto.Security.Key,
		Ca:   dto.Security.Ca,
	}
	return c
}

func (c *Cluster) SetAuthentication(dto *ClusterDTO) *Cluster {
	if dto.Authentication == nil {
		return c
	}
	c.Spec.Authentication = &config.AuthenticationConfig{
		Key:  dto.Authentication.PublicKey,
		Type: dto.Authentication.PublicKeyType,
	}
	return c
}
func (c *Cluster) SetAuthorization(dto *ClusterDTO) *Cluster {
	if dto.Authorization == nil {
		return c
	}
	c.Spec.Authorization = &config.AuthorizationConfig{}
	switch dto.Authorization.Mode {
	case "withUrl":
		c.Spec.Authorization.Url = dto.Authorization.URL
		c.Spec.Authorization.AutoReload = int32(dto.Authorization.FetchInterval)
	case "withPolicy":
		if len(dto.Authorization.Policy.Rules) > 0 {
			data, err := json.Marshal(dto.Authorization.Policy.Rules)
			if err != nil {
				return c
			}
			c.Spec.Authorization.Policy = string(data)
		}
	}

	return c
}

func (c *Cluster) SetRouting(dto *ClusterDTO) *Cluster {
	if dto.Routing == nil {
		return c
	}
	c.Spec.Routing = &config.RoutingConfig{}
	switch dto.Routing.Mode {
	case "withUrl":
		c.Spec.Routing.Url = dto.Routing.URL
		c.Spec.Routing.AutoReload = int32(dto.Routing.FetchInterval)
	case "withRoutes":
		var routesMap []map[string]string
		for _, route := range dto.Routing.Routes.KeyRoutes {

			var routes []string
			for _, e := range strings.Split(route.Events, ";") {
				routes = append(routes, fmt.Sprintf("events:%s", e))
			}
			for _, es := range strings.Split(route.EventsStore, ";") {
				routes = append(routes, fmt.Sprintf("events_store:%s", es))
			}
			for _, q := range strings.Split(route.Queues, ";") {
				routes = append(routes, fmt.Sprintf("queues:%s", q))
			}
			kv := map[string]string{}
			kv["Key"] = route.Key
			kv["Routes"] = strings.Join(routes, ";")
			routesMap = append(routesMap, kv)
		}

		data, _ := json.Marshal(routesMap)
		c.Spec.Routing.Data = string(data)
	}

	return c
}
func (c *Cluster) SetImage(dto *ClusterDTO) *Cluster {
	if dto.Image == nil {
		return c
	}
	c.Spec.Image = &config.ImageConfig{
		Image:      dto.Image.Image,
		PullPolicy: dto.Image.PullPolicy,
	}
	return c
}

func (c *Cluster) SetVolume(dto *ClusterDTO) *Cluster {
	if dto.Volume == nil {
		return c
	}
	c.Spec.Volume = &config.VolumeConfig{
		Size:         fmt.Sprintf("%dGi", dto.Volume.Size),
		StorageClass: dto.Volume.StorageClass,
	}
	return c
}
func (c *Cluster) SetHealth(dto *ClusterDTO) *Cluster {
	if dto.Health == nil {
		return c
	}
	c.Spec.Health = &config.HealthConfig{
		Enabled:             true,
		InitialDelaySeconds: int32(dto.Health.InitialDelaySeconds),
		PeriodSeconds:       int32(dto.Health.PeriodSeconds),
		TimeoutSeconds:      int32(dto.Health.TimeoutSeconds),
		SuccessThreshold:    int32(dto.Health.SuccessThreshold),
		FailureThreshold:    int32(dto.Health.FailureThreshold),
	}
	return c
}
func (c *Cluster) SetResources(dto *ClusterDTO) *Cluster {
	if dto.Resources == nil {
		return c
	}
	c.Spec.Resources = &config.ResourceConfig{
		LimitsCpu:                fmt.Sprintf("%d", dto.Resources.LimitsCPU),
		LimitsMemory:             fmt.Sprintf("%dGi", dto.Resources.LimitsMemory),
		LimitsEphemeralStorage:   fmt.Sprintf("%dGi", dto.Resources.LimitsEphemeralStorage),
		RequestsCpu:              fmt.Sprintf("%d", dto.Resources.RequestCPU),
		RequestsMemory:           fmt.Sprintf("%dGi", dto.Resources.RequestMemory),
		RequestsEphemeralStorage: fmt.Sprintf("%dGi", dto.Resources.RequestsEphemeralStorage),
	}
	return c
}
func (c *Cluster) SetNodes(dto *ClusterDTO) *Cluster {
	if dto.Nodes == nil {
		return c
	}
	keys := map[string]string{}
	for _, item := range dto.Nodes.Items.Kv {
		keys[item.Key] = item.Value
	}
	c.Spec.NodeSelectors = &config.NodeSelectorConfig{
		Keys: keys,
	}
	return c
}
func (c *Cluster) SetStore(dto *ClusterDTO) *Cluster {
	if dto.Store == nil {
		return c
	}
	c.Spec.Store = &config.StoreConfig{
		MaxChannels:              new(int32),
		MaxSubscribers:           new(int32),
		MaxMessages:              new(int32),
		MaxChannelSize:           new(int32),
		MessagesRetentionMinutes: new(int32),
		PurgeInactiveMinutes:     new(int32),
	}
	*c.Spec.Store.MaxChannels = int32(dto.Store.MaxChannels)
	*c.Spec.Store.MaxSubscribers = int32(dto.Store.MaxSubscribers)
	*c.Spec.Store.MaxChannelSize = int32(dto.Store.MaxChannelSize)
	*c.Spec.Store.MessagesRetentionMinutes = int32(dto.Store.MessagesRetentionMinutes)
	*c.Spec.Store.PurgeInactiveMinutes = int32(dto.Store.PurgeInactiveMinutes)
	return c
}
func (c *Cluster) SetQueues(dto *ClusterDTO) *Cluster {
	if dto.Queues == nil {
		return c
	}
	c.Spec.Queue = &config.QueueConfig{
		MaxReceiveMessagesRequest: new(int32),
		MaxWaitTimeoutSeconds:     new(int32),
		MaxExpirationSeconds:      new(int32),
		MaxDelaySeconds:           new(int32),
		MaxReQueues:               new(int32),
		MaxVisibilitySeconds:      new(int32),
		DefaultVisibilitySeconds:  new(int32),
		DefaultWaitTimeoutSeconds: new(int32),
	}
	*c.Spec.Queue.MaxReceiveMessagesRequest = int32(dto.Queues.MaxReceiveMessagesRequest)
	*c.Spec.Queue.MaxWaitTimeoutSeconds = int32(dto.Queues.MaxWaitTimeoutSeconds)
	*c.Spec.Queue.MaxExpirationSeconds = int32(dto.Queues.MaxExpirationSeconds)
	*c.Spec.Queue.MaxDelaySeconds = int32(dto.Queues.MaxDelaySeconds)
	*c.Spec.Queue.MaxReQueues = int32(dto.Queues.MaxReQueues)
	*c.Spec.Queue.MaxVisibilitySeconds = int32(dto.Queues.MaxVisibilitySeconds)
	*c.Spec.Queue.DefaultVisibilitySeconds = int32(dto.Queues.DefaultVisibilitySeconds)
	*c.Spec.Queue.DefaultWaitTimeoutSeconds = int32(dto.Queues.DefaultWaitTimeoutSeconds)
	return c
}
