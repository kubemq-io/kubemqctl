package register

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/kubemq"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"sort"
	"strings"
)

type RegisterOptions struct {
	cfg *config.Config
}

var registerExamples = `
 	# Register KubeMQ cluster
	kubemqctl utils register
`
var registerLong = `Register a trial version of KubeMQ cluster`
var registerShort = `Register KubeMQ cluster command`

func NewCmdRegister(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &RegisterOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "register",
		Aliases: []string{"reg"},
		Short:   registerShort,
		Long:    registerLong,
		Example: registerExamples,
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

func (o *RegisterOptions) Complete(args []string) error {
	return nil
}

func (o *RegisterOptions) Validate() error {

	return nil
}

func (o *RegisterOptions) Run(ctx context.Context) error {
	c, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}

	sets, err := c.GetStatefulSets("")
	if err != nil {
		return err
	}
	utils.Println("Please wait, fetching KubeMQ clusters registration data...")
	if len(sets) == 0 {
		return fmt.Errorf("no KubeMQ clusters were found to register")
	}
	keys := make(map[string]string)
	var list []string
	for key, set := range sets {
		regKey := set.Spec.Template.Annotations["kubemq.io/deploy-reg-id"]
		for _, container := range set.Spec.Template.Spec.Containers {
			if strings.Contains(container.Image, "kubemq") && regKey != "" {
				for _, env := range container.Env {
					if env.Name == "KUBEMQ_TOKEN" {
						token := env.Value
						err := kubemq.ValidateRegisterKubeMQ(token, regKey)
						if err == nil {
							keys[regKey] = token
							list = append(list, key)
						}
						continue
					}
				}
			}
		}

	}
	if len(list) == 0 {
		return fmt.Errorf("KubeMQ clusters were found but none of them need to register")
	}

	sort.Strings(list)

	var selected string
	selection := &survey.Select{
		Renderer:      survey.Renderer{},
		Message:       "Select KubeMQ clusters to register",
		Options:       list,
		Default:       nil,
		Help:          "Select KubeMQ clusters to register",
		PageSize:      0,
		VimMode:       false,
		FilterMessage: "",
		Filter:        nil,
	}
	err = survey.AskOne(selection, &selected)
	if err != nil {
		return err
	}
	form := struct {
		Name     string
		Email    string
		Password string
	}{}
	var formEntry = []*survey.Question{
		{
			Name:      "name",
			Prompt:    &survey.Input{Message: "What is your name?"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name:      "email",
			Prompt:    &survey.Input{Message: "What is your Email address? (We will send a confirmation email to complete registration to this address)"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name:      "password",
			Prompt:    &survey.Password{Message: "Please type your password"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
	}
	err = survey.Ask(formEntry, &form)
	if err != nil {
		return err
	}

	regRequest := &kubemq.RegistrationRequest{
		Name:            form.Name,
		Username:        form.Email,
		Password:        form.Password,
		GeneratedKey:    "",
		RegistrationKey: "",
	}

	set := sets[selected]
	regRequest.RegistrationKey = set.Spec.Template.Annotations["kubemq.io/deploy-reg-id"]
	regRequest.GeneratedKey = keys[regRequest.RegistrationKey]
	utils.Println("Please wait while attempting to register your KubeMQ cluster")
	err = kubemq.RegisterKubeMQ(regRequest)
	if err != nil {
		return err
	}

	utils.Printlnf("Please check your email account (%s) for a confirmation email and complete registration", regRequest.Username)
	return nil
}
