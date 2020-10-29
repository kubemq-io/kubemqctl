package connector

import (
	"context"
	builder "github.com/kubemq-hub/builder/connector"
	contextmanager "github.com/kubemq-io/kubemqctl/cmd/manage/context"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/connector"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/operator"
	operatorTypes "github.com/kubemq-io/kubemqctl/pkg/k8s/types/operator"
)

type Handler struct {
	contextHandler *contextmanager.Handler
}

func NewHandler() *Handler {
	return &Handler{}
}
func (h *Handler) Name() string {
	return ""
}
func (h *Handler) Init(ctx context.Context, contextHandler *contextmanager.Handler) error {
	h.contextHandler = contextHandler

	return nil
}
func (h *Handler) createOrUpdate(con *builder.Connector) error {
	connectorManager, err := connector.NewManager(h.contextHandler.GetClient())
	if err != nil {
		return err
	}
	operatorManager, err := operator.NewManager(h.contextHandler.GetClient())
	if err != nil {
		return err
	}

	dep := ToDeployment(con)
	if !operatorManager.IsKubemqOperatorExists(con.Namespace) {
		operatorDeployment, err := operatorTypes.CreateDeployment("kubemq-operator", dep.Namespace)
		if err != nil {
			return err
		}
		_, _, err = operatorManager.CreateOrUpdateKubemqOperator(operatorDeployment)
		if err != nil {
			return err
		}

	}

	_, _, err = connectorManager.CreateOrUpdateKubemqConnector(dep)
	if err != nil {
		return err
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
	connectorManager, err := connector.NewManager(h.contextHandler.GetClient())
	if err != nil {
		return err
	}
	return connectorManager.DeleteKubemqConnector(ToDeployment(con))
}

func (h *Handler) Get(namespace, name string) (*builder.Connector, error) {
	connectorManager, err := connector.NewManager(h.contextHandler.GetClient())
	if err != nil {
		return nil, err
	}
	dep, err := connectorManager.GetConnector(name, namespace)
	if err != nil {
		return nil, err
	}
	return FromDeployment(dep), nil
}

func (h *Handler) List() ([]*builder.Connector, error) {
	connectorManager, err := connector.NewManager(h.contextHandler.GetClient())
	if err != nil {
		return nil, err
	}
	list, err := connectorManager.GetKubemqConnectors()
	if err != nil {
		return nil, err
	}
	return FromDeploymentList(list), nil
}
