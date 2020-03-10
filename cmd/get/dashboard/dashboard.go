package dashboard

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/cmd/get/dashboard/describe"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/dashboard"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"strings"
)

type getOptions struct {
	cfg *config.Config
}

var getExamples = `
	# Get status of Kubemq of dashboards
	kubemqctl get dashboards
`
var getLong = `Get information of Kubemq dashboard resources`
var getShort = `Get information of Kubemq dashboard resources`

func NewCmdGet(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &getOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "dashboard",
		Aliases: []string{"d", "dash", "dashboard"},
		Short:   getShort,
		Long:    getLong,
		Example: getExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.AddCommand(describe.NewCmdDescribe(ctx, cfg))
	return cmd
}

func (o *getOptions) Complete(args []string) error {
	return nil
}

func (o *getOptions) Validate() error {

	return nil
}

func (o *getOptions) Run(ctx context.Context) error {
	client, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	dashboardManager, err := dashboard.NewManager(client)
	if err != nil {
		return err
	}

	dashboards, err := dashboardManager.GetKubemqDashboardes()
	if err != nil {
		return err
	}
	if len(dashboards.List()) == 0 {
		return fmt.Errorf("no Kubemq dashboards were found")
	}
	var ns, name string

	if len(dashboards.List()) == 1 {
		ns, name = StringSplit(dashboards.List()[0])
	} else {
		selection := ""
		selected := &survey.Select{
			Renderer:      survey.Renderer{},
			Message:       "Select dashboard to launch",
			Options:       dashboards.List(),
			Default:       dashboards.List()[0],
			PageSize:      0,
			VimMode:       false,
			FilterMessage: "",
			Filter:        nil,
		}
		err = survey.AskOne(selected, &selection)
		if err != nil {
			return err
		}
		ns, name = StringSplit(selection)
	}

	grafnaPort, _, err := k8s.GetDashboardTransport(ctx, o.cfg, ns, name)
	if err != nil {
		return err
	}
	err = browser.OpenURL(fmt.Sprintf("http://localhost:%s/d/kubemqdashboard/kubemq-dashboard", grafnaPort))
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

func StringSplit(input string) (string, string) {
	pair := strings.Split(input, "/")
	if len(pair) == 2 {
		return pair[0], pair[1]
	}
	return "", ""
}
