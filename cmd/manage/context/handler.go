package cluster

import (
	"context"
	"github.com/AlecAivazis/survey/v2"
	"sort"

	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
)

type Handler struct {
	ctx        context.Context
	client     *client.Client
	kubeconfig string
	current    string
}

func NewHandler() *Handler {
	return &Handler{}
}
func (h *Handler) Name() string {
	cfg, _ := h.client.ClientConfig.ClientConfig()
	return cfg.Host
}
func (h *Handler) Init(ctx context.Context, kubeconfig string) error {
	var err error
	h.client, err = client.NewClient(kubeconfig)
	if err != nil {
		return err
	}
	_, current, err := h.client.GetConfigContext()
	if err != nil {
		return err
	}
	h.current = current
	return nil
}
func (h *Handler) GetClient() *client.Client {
	return h.client
}
func (h *Handler) Set() error {
	contextMap, current, err := h.client.GetConfigContext()
	if err != nil {
		return err
	}
	list := []string{}
	for key := range contextMap {
		list = append(list, key)
	}
	sort.Strings(list)
	contextSelected := ""
	contextSelect := &survey.Select{
		Renderer:      survey.Renderer{},
		Message:       "Select kubernetes cluster context",
		Options:       list,
		Default:       current,
		Help:          "Set kubernetes connection context",
		PageSize:      0,
		VimMode:       false,
		FilterMessage: "",
		Filter:        nil,
	}
	err = survey.AskOne(contextSelect, &contextSelected)
	if err != nil {
		return err
	}
	err = h.client.SwitchContext(contextSelected)
	if err != nil {
		return err
	}
	h.current = contextSelected
	h.client, err = client.NewClient(h.kubeconfig)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) Get() string {
	return h.current
}
