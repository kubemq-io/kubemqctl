package authorization

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

type PolicyOptions struct {
	cfg    *config.Config
	verify bool
}

var policyExamples = `
	# Execute generate authorization policy file
 	kubemqctl generate az policy
`
var policyLong = `Generate KubeMQ Authorization access control file`
var policyShort = `Generate KubeMQ Authorization access control file`

func NewCmdPolicy(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &PolicyOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{
		Use:     "policy",
		Aliases: []string{"p"},
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

func (o *PolicyOptions) Complete(args []string, transport string) error {
	return nil
}

func (o *PolicyOptions) Validate() error {
	return nil
}
func (o *PolicyOptions) Run(ctx context.Context) error {
	var rules []*Rule
	utils.Println("Create first rule:")
	for {
		r, err := getRule()
		if err != nil {
			return err
		}
		rules = append(rules, r)
		addMoreRule := true
		addMorePrompt := &survey.Confirm{
			Message: "Add more rules to policy?",
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
		utils.Println("Create next rule:")
	}
save:
	data, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return err
	}
	utils.Println("Policy Rules:")
	fmt.Println(string(data))
	err = ioutil.WriteFile("policy.json", data, 0600)
	if err != nil {
		return err
	}
	utils.Println("Policy data save to policy.json file")
	return nil
}
