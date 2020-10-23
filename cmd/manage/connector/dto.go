package connector

import (
	builder "github.com/kubemq-hub/builder/connector"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/connector"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqconnector"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ToDeployment(con *builder.Connector) *kubemqconnector.KubemqConnector {
	deployment := &kubemqconnector.KubemqConnector{
		TypeMeta: v1.TypeMeta{
			Kind:       "KubemqConnector",
			APIVersion: "core.k8s.kubemq.io/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      con.Name,
			Namespace: con.Namespace,
		},
		Spec: kubemqconnector.KubemqConnectorSpec{
			Replicas:    new(int32),
			Type:        con.Type,
			Image:       con.Image,
			Config:      con.Config,
			NodePort:    int32(con.NodePort),
			ServiceType: con.ServiceType,
		},
		Status: kubemqconnector.KubemqConnectorStatus{},
	}
	*deployment.Spec.Replicas = int32(con.Replicas)
	return deployment
}

func FromDeployment(deployment *kubemqconnector.KubemqConnector) *builder.Connector {
	con := &builder.Connector{
		Name:        deployment.Name,
		Namespace:   deployment.Namespace,
		Type:        deployment.Spec.Type,
		Replicas:    0,
		Config:      deployment.Spec.Config,
		NodePort:    int(deployment.Spec.NodePort),
		ServiceType: deployment.Spec.ServiceType,
		Image:       deployment.Spec.Image,
	}
	if deployment.Spec.Replicas != nil {
		con.Replicas = int(*deployment.Spec.Replicas)
	}
	return con
}

func FromDeploymentList(list *connector.KubemqConnectors) []*builder.Connector {
	var connectorsList []*builder.Connector
	for _, deployment := range list.Items() {
		connectorsList = append(connectorsList, FromDeployment(deployment))
	}
	return connectorsList
}
