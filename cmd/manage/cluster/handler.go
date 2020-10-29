package cluster

import (
	"context"

	builder "github.com/kubemq-hub/builder/cluster"
	"github.com/kubemq-hub/builder/connector"
	contextmanager "github.com/kubemq-io/kubemqctl/cmd/manage/context"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/cluster"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/operator"
	operatorTypes "github.com/kubemq-io/kubemqctl/pkg/k8s/types/operator"
)

type Handler struct {
	ctx              context.Context
	cfg              *config.Config
	contextHandler   *contextmanager.Handler
	connectorHandler connector.ConnectorsHandler
}

func (h *Handler) ConnectorsHandler() connector.ConnectorsHandler {
	return h.connectorHandler
}

func NewHandler() *Handler {
	return &Handler{}
}
func (h *Handler) Name() string {
	return ""
}
func (h *Handler) Init(ctx context.Context, contextHandler *contextmanager.Handler, connectorHandler connector.ConnectorsHandler) error {
	h.connectorHandler = connectorHandler
	h.contextHandler = contextHandler
	return nil
}
func (h *Handler) createOrUpdate(cls *builder.Cluster) error {
	clusterManager, err := cluster.NewManager(h.contextHandler.GetClient())
	if err != nil {
		return err
	}
	operatorManager, err := operator.NewManager(h.contextHandler.GetClient())
	if err != nil {
		return err
	}
	dep := ToDeployment(cls)
	if !operatorManager.IsKubemqOperatorExists(cls.Namespace) {
		operatorDeployment, err := operatorTypes.CreateDeployment("kubemq-operator", dep.Namespace)
		if err != nil {
			return err
		}
		_, _, err = operatorManager.CreateOrUpdateKubemqOperator(operatorDeployment)
		if err != nil {
			return err
		}
	}

	_, _, err = clusterManager.CreateOrUpdateKubemqCluster(dep)
	if err != nil {
		return err
	}

	return err
}
func (h *Handler) Add(cls *builder.Cluster) error {
	return h.createOrUpdate(cls)
}

func (h *Handler) Edit(cls *builder.Cluster) error {
	return h.createOrUpdate(cls)
}

func (h *Handler) Delete(cls *builder.Cluster) error {
	clusterManager, err := cluster.NewManager(h.contextHandler.GetClient())
	if err != nil {
		return err
	}
	return clusterManager.DeleteKubemqCluster(ToDeployment(cls))
}

func (h *Handler) Get(namespace, name string) (*builder.Cluster, error) {
	clusterManager, err := cluster.NewManager(h.contextHandler.GetClient())
	if err != nil {
		return nil, err
	}
	dep, err := clusterManager.GetCluster(name, namespace)
	if err != nil {
		return nil, err
	}
	return FromDeployment(dep), nil
}

func (h *Handler) List() ([]*builder.Cluster, error) {
	clusterManager, err := cluster.NewManager(h.contextHandler.GetClient())
	if err != nil {
		return nil, err
	}
	list, err := clusterManager.GetKubemqClusters()
	if err != nil {
		return nil, err
	}
	return FromDeploymentList(list), nil
}
