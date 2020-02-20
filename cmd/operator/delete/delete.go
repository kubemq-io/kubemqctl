package delete

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type DeleteOptions struct {
	cfg *config.Config
}

var deleteExamples = `
	# Delete KubeMQ operators 
	kubemqctl operator delete  
`
var deleteLong = `Get command display all operators deployed across all namespaces`
var deleteShort = `Get KubeMQ Operators List`

func NewCmdDelete(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &DeleteOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "delete",
		Aliases: []string{"del"},
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
	mng, err := manager.NewManager(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	utils.Println("Getting KubeMQ Operators List...")
	list, err := mng.GetOperatorList()
	if err != nil {
		return err
	}
	if list == nil {
		return fmt.Errorf("no KubeMQ operators found in cluster")
	}

	selection := []string{}
	multiSelected := &survey.MultiSelect{
		Renderer:      survey.Renderer{},
		Message:       "Select KubeMQ operator namespace to delete",
		Options:       list,
		Default:       nil,
		Help:          "Select KubeMQ operator namespace to delete",
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
		Help:     "Confirm KubeMQ operator deletion",
	}
	err = survey.AskOne(promptConfirm, &areYouSure)
	if err != nil {
		return err
	}
	if !areYouSure {
		return nil
	}
	for _, ns := range selection {
		bundle, err := mng.GetOperator(ns)
		if err != nil {
			return fmt.Errorf("cannot get operator information for namespace %s, error: %s", ns, err.Error())
		}
		err = mng.DeleteOperator(bundle)
		if err != nil {
			return fmt.Errorf("error delete operator for namespace %s, error: %s", ns, err.Error())
		}
		utils.Printlnf("KubeMQ Operator at namespace %s deleted.", ns)
	}
	return nil
}
