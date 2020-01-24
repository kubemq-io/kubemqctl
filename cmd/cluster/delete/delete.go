package delete

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type DeleteOptions struct {
	cfg *config.Config
}

var deleteExamples = `
 	# Delete KubeMQ cluster
	kubemqctl cluster delete
`
var deleteLong = `Delete command allows deleting one or more KubeMQ clusters deployments`
var deleteShort = `Delete KubeMQ cluster command`

func NewCmdDelete(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &DeleteOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "delete",
		Aliases: []string{"del", "de", "d"},
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
	c, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}

	list, err := c.GetKubeMQClusters()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return fmt.Errorf("no KubeMQ clusters were found to delete")
	}
	selection := []string{}
	multiSelected := &survey.MultiSelect{
		Renderer:      survey.Renderer{},
		Message:       "Select KubeMQ clusters to delete",
		Options:       list,
		Default:       nil,
		Help:          "Select KubeMQ clusters to delete",
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
		Help:     "Confirm KubeMQ cluster deletion",
	}
	err = survey.AskOne(promptConfirm, &areYouSure)
	if err != nil {
		return err
	}
	if !areYouSure {
		return nil
	}
	for _, sts := range selection {

		err := c.DeleteStatefulSet(sts)
		if err != nil {
			utils.Printlnf("StatefulSet %s not deleted. Error %s", sts, utils.Title(err.Error()))
		} else {
			utils.Printlnf("StatefulSet %s deleted.", sts)
		}
		err = c.DeleteServicesForStatefulSet(sts)
		if err != nil {
			utils.Printlnf("Delete services failed. Error %s", utils.Title(err.Error()))
		}
		err = c.DeleteConfigMapsForStatefulSet(sts)
		if err != nil {
			utils.Printlnf("Delete config maps failed. Error %s", utils.Title(err.Error()))
		}

		err = c.DeleteSecretsForStatefulSet(sts)
		if err != nil {
			utils.Printlnf("Delete secretes failed. Error %s", utils.Title(err.Error()))
		}
		err = c.DeleteIngressForStatefulSet(sts)
		if err != nil {
			utils.Printlnf("Delete ingress failed. Error %s", utils.Title(err.Error()))
		}

		err = c.DeleteVolumeClaimsForStatefulSet(sts)
		if err != nil {
			utils.Printlnf("Delete persistence volume claims failed. Error %s", utils.Title(err.Error()))
		}
	}
	return nil
}
