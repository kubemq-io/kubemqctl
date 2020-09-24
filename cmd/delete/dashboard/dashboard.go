package dashboard

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	client2 "github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/dashboard"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/operator"
	operatorTypes "github.com/kubemq-io/kubemqctl/pkg/k8s/types/operator"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type DeleteOptions struct {
	cfg *config.Config
}

var deleteExamples = `
 	# Delete Kubemq dashboard
	kubemqctl delete dashboard
`
var deleteLong = `Delete one or more Kubemq dashboards`
var deleteShort = `Delete Kubemq dashboard`

func NewCmdDelete(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &DeleteOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "dashboard",
		Aliases: []string{"d", "dash", "dashboards"},
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
	dashboardManager, err := dashboard.NewManager(client)
	if err != nil {
		return err
	}
	operatorManager, err := operator.NewManager(client)
	if err != nil {
		return err
	}
	dashboards, err := dashboardManager.GetKubemqDashboardes()
	if err != nil {
		return err
	}
	if len(dashboards.List()) == 0 {
		return fmt.Errorf("no Kubemq dashboards were found to delete")
	}

	selection := []string{}
	multiSelected := &survey.MultiSelect{
		Renderer:      survey.Renderer{},
		Message:       "Select Kubemq dashboards to delete",
		Options:       dashboards.List(),
		Default:       nil,
		Help:          "Select Kubemq dashboards to delete",
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
		Help:     "Confirm Kubemq dashboard deletion",
	}
	err = survey.AskOne(promptConfirm, &areYouSure)
	if err != nil {
		return err
	}
	if !areYouSure {
		return nil
	}
	for _, selected := range selection {
		dashboard := dashboards.Dashboard(selected)
		if !operatorManager.IsKubemqOperatorExists(dashboard.Namespace) {
			operatorDeployment, err := operatorTypes.CreateDeployment("kubemq-operator", dashboard.Namespace)
			if err != nil {
				return err
			}
			_, _, err = operatorManager.CreateOrUpdateKubemqOperator(operatorDeployment)
			if err != nil {
				return err
			}
			utils.Printlnf("Kubemq operator %s/kubemq-operator created.", dashboard.Namespace)
		} else {
			utils.Printlnf("Kubemq operator %s/kubemq-operator exists", dashboard.Namespace)
		}
		err := dashboardManager.DeleteKubemqDashboard(dashboard)
		if err != nil {
			return err
		}
		utils.Printlnf("Kubemq dashboard %s deleted.", selected)
	}
	return nil
}
