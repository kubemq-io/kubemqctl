package connector

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	client2 "github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/connector"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type DeleteOptions struct {
	cfg *config.Config
}

var deleteExamples = `
 	# Delete Kubemq connector
	kubemqctl delete connector
`
var deleteLong = `Delete one or more Kubemq connectors`
var deleteShort = `Delete Kubemq connector`

func NewCmdDelete(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &DeleteOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "connector",
		Aliases: []string{"cn", "con", "connectors"},
		Short:   deleteShort,
		Long:    deleteLong,
		Example: deleteExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}

	return cmd
}

func (o *DeleteOptions) Complete(args []string) error {
	return nil
}

func (o *DeleteOptions) Validate() error {

	return nil
}

func (o *DeleteOptions) Run(ctx context.Context) error {
	client, err := client2.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	connectorManager, err := connector.NewManager(client)
	if err != nil {
		return err
	}

	connectors, err := connectorManager.GetKubemqConnectors()
	if err != nil {
		return err
	}
	if len(connectors.List()) == 0 {
		return fmt.Errorf("no Kubemq connectors were found to delete")
	}

	selection := []string{}
	multiSelected := &survey.MultiSelect{
		Renderer:      survey.Renderer{},
		Message:       "Select Kubemq connectors to delete",
		Options:       connectors.List(),
		Default:       nil,
		Help:          "Select Kubemq connectors to delete",
		PageSize:      0,
		VimMode:       false,
		FilterMessage: "",
		Filter:        nil,
	}
	err = survey.AskOne(multiSelected, &selection)
	if err != nil {
		return err
	}

	areYouSure := false
	promptConfirm := &survey.Confirm{
		Renderer: survey.Renderer{},
		Message:  "Are you sure ?",
		Default:  false,
		Help:     "Confirm Kubemq connector deletion",
	}
	err = survey.AskOne(promptConfirm, &areYouSure)
	if err != nil {
		return err
	}
	if !areYouSure {
		return nil
	}
	for _, selected := range selection {
		err := connectorManager.DeleteKubemqConnector(connectors.Connector(selected))
		if err != nil {
			return err
		}
		utils.Printlnf("Kubemq connector %s deleted.", selected)
	}
	return nil
}
