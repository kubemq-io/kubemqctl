package cluster

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	"github.com/spf13/cobra"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type deployOptions struct {
	configData            string
	configFilename        string
	name                  string
	namespace             string
	replicas              int32
	standalone            bool
	statefulSetConfigData string
	key                   string
	api                   *deployApiOptions
	authentication        *deployAuthenticationOptions
	authorization         *deployAuthorizationOptions
	grpc                  *deployGrpcOptions
	health                *deployHealthOptions
	image                 *deployImageOptions
	license               *deployLicenseOptions
	log                   *deployLogOptions
	nodeSelector          *deployNodeSelectorOptions
	notification          *deployNotificationOptions
	queue                 *deployQueueOptions
	resources             *deployResourceOptions
	rest                  *deployRestOptions
	routing               *deployRoutingOptions
	store                 *deployStoreOptions
	tls                   *deployTlsOptions
	volume                *deployVolumeOptions
}

func defaultDeployOptions(cmd *cobra.Command) *deployOptions {
	o := &deployOptions{
		configData:            "",
		configFilename:        "",
		name:                  "",
		namespace:             "",
		replicas:              0,
		standalone:            false,
		statefulSetConfigData: "",
		key:                   "",
		api:                   setApiConfig(cmd),
		authentication:        setAuthenticationOptions(cmd),
		authorization:         setAuthorizationConfig(cmd),
		grpc:                  setGrpcConfig(cmd),
		health:                setHealthOptions(cmd),
		image:                 setImageConfig(cmd),
		license:               setLicenseConfig(cmd),
		log:                   setLogConfig(cmd),
		nodeSelector:          setNodeSelectorOptions(cmd),
		notification:          setNotificationConfig(cmd),
		queue:                 setQueueConfig(cmd),
		resources:             setResourceOptions(cmd),
		rest:                  setRestConfig(cmd),
		routing:               setRoutingConfig(cmd),
		store:                 setStoreConfig(cmd),
		tls:                   setTolsConfig(cmd),
		volume:                setVolumeConfig(cmd),
	}
	cmd.PersistentFlags().StringVarP(&o.configFilename, "config-file", "c", "", "set kubemq config file")
	cmd.PersistentFlags().StringVarP(&o.name, "name", "", "kubemq-cluster", "set kubemq cluster name")
	cmd.PersistentFlags().StringVarP(&o.namespace, "namespace", "n", "kubemq", "set kubemq cluster namespace")
	cmd.PersistentFlags().StringVarP(&o.key, "key", "", "", "set kubemq license key")
	cmd.PersistentFlags().StringVarP(&o.statefulSetConfigData, "statefulset-config-data", "", "", "set kubemq cluster statefulset configuration data")
	cmd.PersistentFlags().BoolVarP(&o.standalone, "standalone", "", false, "set kubemq cluster standalone mode")
	cmd.PersistentFlags().Int32VarP(&o.replicas, "replicas", "r", 3, "set replicas")

	return o
}

func (o *deployOptions) validate() error {
	if o.name == "" {
		return fmt.Errorf("error setting deploy configuration, missing kubemq cluster name")
	}
	if o.namespace == "" {
		return fmt.Errorf("error setting deploy configuration, missing kubemq cluster namespace")
	}

	if err := o.api.validate(); err != nil {
		return err
	}

	if err := o.authentication.validate(); err != nil {
		return err
	}
	if err := o.authorization.validate(); err != nil {
		return err
	}

	if err := o.grpc.validate(); err != nil {
		return err
	}

	if err := o.health.validate(); err != nil {
		return err
	}
	if err := o.image.validate(); err != nil {
		return err
	}
	if err := o.license.validate(); err != nil {
		return err
	}
	if err := o.log.validate(); err != nil {
		return err
	}
	if err := o.notification.validate(); err != nil {
		return err
	}
	if err := o.resources.validate(); err != nil {
		return err
	}

	if err := o.nodeSelector.validate(); err != nil {
		return err
	}
	if err := o.queue.validate(); err != nil {
		return err
	}
	if err := o.rest.validate(); err != nil {
		return err
	}
	if err := o.routing.validate(); err != nil {
		return err
	}
	if err := o.store.validate(); err != nil {
		return err
	}
	if err := o.tls.validate(); err != nil {
		return err
	}
	if err := o.volume.validate(); err != nil {
		return err
	}
	return nil
}

func (o *deployOptions) complete() error {
	if o.configFilename != "" {
		data, err := ioutil.ReadFile(o.configFilename)
		if err != nil {
			return fmt.Errorf("error config file data: %s", err.Error())
		}
		o.configData = string(data)
	}

	if err := o.api.complete(); err != nil {
		return err
	}

	if err := o.authentication.complete(); err != nil {
		return err
	}
	if err := o.authorization.complete(); err != nil {
		return err
	}
	if err := o.grpc.complete(); err != nil {
		return err
	}

	if err := o.health.complete(); err != nil {
		return err
	}
	if err := o.image.complete(); err != nil {
		return err
	}
	if err := o.license.complete(); err != nil {
		return err
	}
	if err := o.log.complete(); err != nil {
		return err
	}
	if err := o.notification.complete(); err != nil {
		return err
	}
	if err := o.resources.complete(); err != nil {
		return err
	}

	if err := o.nodeSelector.complete(); err != nil {
		return err
	}
	if err := o.queue.complete(); err != nil {
		return err
	}
	if err := o.rest.complete(); err != nil {
		return err
	}
	if err := o.routing.complete(); err != nil {
		return err
	}
	if err := o.store.complete(); err != nil {
		return err
	}
	if err := o.tls.complete(); err != nil {
		return err
	}
	if err := o.volume.complete(); err != nil {
		return err
	}
	return nil
}

func (o *deployOptions) getClusterDeployment() *kubemqcluster.KubemqCluster {

	deployment := &kubemqcluster.KubemqCluster{
		TypeMeta: v1.TypeMeta{
			Kind:       "KubemqCluster",
			APIVersion: "core.k8s.kubemq.io/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      o.name,
			Namespace: o.namespace,
		},
		Spec: kubemqcluster.KubemqClusterSpec{
			Replicas:              new(int32),
			License:               "",
			ConfigData:            o.configData,
			Key:                   o.key,
			Standalone:            o.standalone,
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
			StatefulSetConfigData: o.statefulSetConfigData,
		},
		Status: kubemqcluster.KubemqClusterStatus{},
	}
	*deployment.Spec.Replicas = o.replicas
	o.api.setConfig(deployment)
	o.authentication.setConfig(deployment)
	o.authorization.setConfig(deployment)
	o.grpc.setConfig(deployment)
	o.health.setConfig(deployment)
	o.image.setConfig(deployment)
	o.license.setConfig(deployment)
	o.log.setConfig(deployment)
	o.nodeSelector.setConfig(deployment)
	o.notification.setConfig(deployment)
	o.queue.setConfig(deployment)
	o.rest.setConfig(deployment)
	o.resources.setConfig(deployment)
	o.routing.setConfig(deployment)
	o.store.setConfig(deployment)
	o.tls.setConfig(deployment)
	o.volume.setConfig(deployment)

	return deployment
}
