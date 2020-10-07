package build

import (
	"context"
	"fmt"
	clusterbuilder "github.com/kubemq-hub/builder/cluster"
	"github.com/kubemq-hub/builder/common"
	connectorbuilder "github.com/kubemq-hub/builder/connector"
	"github.com/kubemq-hub/builder/survey"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/cluster"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/connector"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/operator"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqconnector"
	operatorTypes "github.com/kubemq-io/kubemqctl/pkg/k8s/types/operator"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type buildOptions struct {
	cfg           *config.Config
	output        string
	injectOptions common.DefaultOptions
	deployments   []string
}

var buildExamples = `
	# Execute build Kubemq components
	kubemqctl build	
	
	# Execute build and export yaml
	kubemqctl build -o deploy.yaml
`
var buildLong = `Executes Kubemq build command`
var buildShort = `Executes Kubemq build command`

func NewCmdBuild(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &buildOptions{
		cfg:           cfg,
		output:        "",
		injectOptions: common.NewDefaultOptions(),
		deployments:   nil,
	}
	cmd := &cobra.Command{

		Use:     "build",
		Aliases: []string{"b"},
		Short:   buildShort,
		Long:    buildLong,
		Example: buildExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))

		},
	}
	cmd.PersistentFlags().StringVarP(&o.output, "output", "o", "", "set output yaml file name")
	return cmd
}
func (o *buildOptions) Complete(args []string) error {
	return nil
}

func (o *buildOptions) Validate() error {
	return nil
}

func (o *buildOptions) Run(ctx context.Context) error {
	client, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	if err := o.updateInjectOption(client); err != nil {
		return err
	}
	comp, err := o.askComponent()
	if err != nil {
		return err
	}
	isDeploy := o.output == ""
	switch comp {
	case "KubeMQ Cluster":
		deployment, err := o.renderCluster()
		if err != nil {
			return err
		}
		if isDeploy {
			if err := o.deployCluster(client, deployment); err != nil {
				return err
			}
		} else {
			o.deployments = append(o.deployments, deployment.String())
		}
	case "KubeMQ Connector":
		deployment, err := o.renderConnector()
		if err != nil {
			return err
		}
		if isDeploy {
			if err := o.deployConnector(client, deployment); err != nil {
				return err
			}
		} else {
			o.deployments = append(o.deployments, deployment.String())
		}
	}

	if len(o.deployments) > 0 {
		return o.saveDeployments()
	}
	return nil
}
func (o *buildOptions) saveDeployments() error {
	output := strings.Join(o.deployments, "\n---\n")
	return ioutil.WriteFile(o.output, []byte(output), 0644)
}
func (o *buildOptions) deployOperator(client *client.Client, namespace string) error {
	operatorManager, err := operator.NewManager(client)
	if err != nil {
		return err
	}
	if !operatorManager.IsKubemqOperatorExists(namespace) {
		operatorDeployment, err := operatorTypes.CreateDeployment("kubemq-operator", namespace)
		if err != nil {
			return err
		}
		_, _, err = operatorManager.CreateOrUpdateKubemqOperator(operatorDeployment)
		if err != nil {
			return err
		}
		utils.Printlnf("Kubemq operator %s/kubemq-operator created.", namespace)
	} else {
		utils.Printlnf("Kubemq operator %s/kubemq-operator exists", namespace)
	}
	return nil
}
func (o *buildOptions) deployCluster(client *client.Client, deployment *kubemqcluster.KubemqCluster) error {
	utils.Printlnf("Deploying KubeMQ Cluster...")
	clusterManager, err := cluster.NewManager(client)
	if err != nil {
		return err
	}
	if err := o.deployOperator(client, deployment.Namespace); err != nil {
		return err
	}
	cluster, isUpdate, err := clusterManager.CreateOrUpdateKubemqCluster(deployment)
	if err != nil {
		return err
	}
	if isUpdate {
		utils.Printlnf("kubemq cluster %s/%s configured.", cluster.Namespace, cluster.Name)
	} else {
		utils.Printlnf("kubemq cluster %s/%s created.", cluster.Namespace, cluster.Name)
	}
	return nil
}

func (o *buildOptions) deployConnector(client *client.Client, deployment *kubemqconnector.KubemqConnector) error {
	connectorManager, err := connector.NewManager(client)
	if err != nil {
		return err
	}
	if err := o.deployOperator(client, deployment.Namespace); err != nil {
		return err
	}
	connector, isUpdate, err := connectorManager.CreateOrUpdateKubemqConnector(deployment)
	if err != nil {
		return err
	}
	if isUpdate {
		utils.Printlnf("kubemq connector %s/%s configured.", connector.Namespace, connector.Name)
	} else {
		utils.Printlnf("kubemq connector %s/%s created.", connector.Namespace, connector.Name)
	}
	return nil
}

func (o *buildOptions) updateInjectOption(client *client.Client) error {
	clusterManager, err := cluster.NewManager(client)
	if err != nil {
		return err
	}
	clusters, err := clusterManager.GetKubemqClusters()
	if err != nil {
		return err
	}
	var kubemqAddress []string
	for _, c := range clusters.Items {
		kubemqAddress = append(kubemqAddress, fmt.Sprintf("%s-grpc.%s.svc.local", c.Name, c.Namespace))
	}

	o.injectOptions.Add("kubemq-address", kubemqAddress)
	nsList, err := client.GetNamespaceList()
	if err != nil {
		return nil
	}
	o.injectOptions.Add("namespaces", nsList)
	return nil
}

func (o *buildOptions) askComponent() (string, error) {
	val := ""
	err := survey.NewString().
		SetKind("string").
		SetName("select-component").
		SetMessage("Select KubeMQ Component to build").
		SetDefault("KubeMQ Cluster").
		SetHelp("Sets KubeMQ Component to build").
		SetOptions([]string{"KubeMQ Cluster", "KubeMQ Connector"}).
		SetRequired(true).
		Render(&val)
	if err != nil {
		return "", err
	}
	return val, nil
}

func (o *buildOptions) renderCluster() (*kubemqcluster.KubemqCluster, error) {
	cluster, err := clusterbuilder.
		NewCluster().
		SetNamespaces(o.injectOptions["namespaces"]).
		Render()
	if err != nil {
		return nil, err
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
	return deployment, nil
}

func (o *buildOptions) renderConnector() (*kubemqconnector.KubemqConnector, error) {
	connector, err := connectorbuilder.NewConnector().
		SetDefaultOptions(o.injectOptions).
		Render()
	if err != nil {
		return nil, err
	}
	deployment := &kubemqconnector.KubemqConnector{
		TypeMeta: v1.TypeMeta{
			Kind:       "KubemqConnector",
			APIVersion: "core.k8s.kubemq.io/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", connector.Name, connector.Type),
			Namespace: connector.Namespace,
		},
		Spec: kubemqconnector.KubemqConnectorSpec{
			Replicas:    new(int32),
			Type:        connector.Type,
			Image:       connector.Image,
			Config:      connector.Config,
			NodePort:    int32(connector.NodePort),
			ServiceType: connector.ServiceType,
		},
		Status: kubemqconnector.KubemqConnectorStatus{},
	}

	return deployment, nil

}
