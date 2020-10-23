package connector

import (
	"context"
	builder "github.com/kubemq-hub/builder/connector"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/connector"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/operator"
	operatorTypes "github.com/kubemq-io/kubemqctl/pkg/k8s/types/operator"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
)

type Handler struct {
	ctx              context.Context
	cfg              *config.Config
	client           *client.Client
	connectorManager *connector.Manager
	operatorManager  *operator.Manager
}

func NewHandler() *Handler {
	return &Handler{}
}
func (h *Handler) Init(ctx context.Context, cfg *config.Config) error {
	var err error
	h.client, err = client.NewClient(cfg.KubeConfigPath)
	if err != nil {
		return err
	}

	h.connectorManager, err = connector.NewManager(h.client)
	if err != nil {
		return err
	}

	h.operatorManager, err = operator.NewManager(h.client)
	if err != nil {
		return err
	}
	return nil
}
func (h *Handler) createOrUpdate(con *builder.Connector) error {
	dep := ToDeployment(con)
	if !h.operatorManager.IsKubemqOperatorExists(con.Namespace) {
		operatorDeployment, err := operatorTypes.CreateDeployment("kubemq-operator", dep.Namespace)
		if err != nil {
			return err
		}
		_, _, err = h.operatorManager.CreateOrUpdateKubemqOperator(operatorDeployment)
		if err != nil {
			return err
		}
		utils.Printlnf("<cyan>Kubemq operator %s/kubemq-operator created</>", dep.Namespace)
	}

	connector, isUpdate, err := h.connectorManager.CreateOrUpdateKubemqConnector(dep)
	if err != nil {
		return err
	}
	if isUpdate {
		utils.Printlnf("<cyan>Kubemq connector %s/%s configured</>", connector.Namespace, connector.Name)
	} else {
		utils.Printlnf("<cyan>Kubemq connector %s/%s created</>", connector.Namespace, connector.Name)
	}
	return err
}
func (h *Handler) Add(con *builder.Connector) error {
	return h.createOrUpdate(con)
}

func (h *Handler) Edit(con *builder.Connector) error {
	return h.createOrUpdate(con)
}

func (h *Handler) Delete(con *builder.Connector) error {
	return h.connectorManager.DeleteKubemqConnector(ToDeployment(con))
}

func (h *Handler) Get(namespace, name string) (*builder.Connector, error) {
	dep, err := h.connectorManager.GetConnector(name, namespace)
	if err != nil {
		return nil, err
	}
	return FromDeployment(dep), nil
}

func (h *Handler) List() ([]*builder.Connector, error) {
	list, err := h.connectorManager.GetKubemqConnectors()
	if err != nil {
		return nil, err
	}

	return FromDeploymentList(list), nil
}
