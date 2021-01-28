package routing

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
)

type RoutingOptions struct {
	cfg    *config.Config
	verify bool
}

var policyExamples = `
	# Execute generate smart routing file
 	kubemqctl generate routes
`
var policyLong = `Generate KubeMQ Smart Routing file`
var policyShort = `Generate KubeMQ Smart Routing file`

func NewCmdRouting(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &RoutingOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{
		Use:     "routes",
		Aliases: []string{"r", "route"},
		Short:   policyShort,
		Long:    policyLong,
		Example: policyExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}

	return cmd
}

func (o *RoutingOptions) Complete(args []string, transport string) error {
	return nil
}

func (o *RoutingOptions) Validate() error {
	return nil
}
func (o *RoutingOptions) Run(ctx context.Context) error {
	var routes []*Route
	utils.Println("Create first route:")
	for {
		r, err := getRoute()
		if err != nil {
			return err
		}
		routes = append(routes, r)
		addMoreRule := true
		addMorePrompt := &survey.Confirm{
			Message: "Add more routes ?",
			Default: true,
			Help:    "",
		}
		err = survey.AskOne(addMorePrompt, &addMoreRule)
		if err != nil {
			return err
		}
		if !addMoreRule {
			goto save
		}
		utils.Println("Create next route:")
	}
save:
	data, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		return err
	}
	utils.Println("Routing Rules:")
	fmt.Println(string(data))
	err = ioutil.WriteFile("routes.json", data, 0600)
	if err != nil {
		return err
	}
	utils.Println("Routing data save to routes.json file")
	return nil
}
