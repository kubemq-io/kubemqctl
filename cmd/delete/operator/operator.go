package operator

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"

	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/operator"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type DeleteOptions struct {
	cfg   *config.Config
	isAll bool
}

var deleteExamples = `
	# Delete Kubemq operator 
	kubemqctl delete operator  
`
var deleteLong = `Delete one or more Kubemq operators`
var deleteShort = `Delete Kubemq operator`

func NewCmdDelete(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &DeleteOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "operator",
		Aliases: []string{"o", "op", "operators"},
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
	cmd.PersistentFlags().BoolVarP(&o.isAll, "remove-all", "", false, "remove all operator components")
	return cmd
}

func (o *DeleteOptions) Complete(args []string) error {
	return nil
}

func (o *DeleteOptions) Validate() error {
	return nil
}

func (o *DeleteOptions) Run(ctx context.Context) error {
	newClient, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}

	operatorManager, err := operator.NewManager(newClient)
	if err != nil {
		return err
	}

	utils.Println("Getting Kubemq Operators List...")
	operators, err := operatorManager.GetKubemqOperators()
	if err != nil {
		return err
	}
	if len(operators.List()) == 0 {
		return fmt.Errorf("no Kubemq operators found in cluster")
	}

	selection := []string{}
	multiSelected := &survey.MultiSelect{
		Renderer:      survey.Renderer{},
		Message:       "Select Kubemq operator namespace to delete",
		Options:       operators.List(),
		Default:       nil,
		Help:          "Select Kubemq operator namespace to delete",
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
		Help:     "Confirm Kubemq operator deletion",
	}
	err = survey.AskOne(promptConfirm, &areYouSure)
	if err != nil {
		return err
	}
	if !areYouSure {
		return nil
	}
	for _, name := range selection {

		err = operatorManager.DeleteKubemqOperator(operators.Deployment(name), o.isAll)
		if err != nil {
			return fmt.Errorf("error delete operator for namespace %s, error: %s", name, err.Error())
		}
		utils.Printlnf("Kubemq Operator at namespace %s deleted.", name)
	}
	return nil
}
