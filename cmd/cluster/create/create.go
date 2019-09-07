package create

import (
	"bufio"
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/ghodss/yaml"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s/client"
	conf "github.com/kubemq-io/kubetools/pkg/k8s/config"
	"github.com/skratchdot/open-golang/open"
	"k8s.io/api/core/v1"
	"os"

	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
)

type CreateOptions struct {
	cfg           *config.Config
	defaultMode   bool
	exportFile    bool
	token         string
	replicas      int
	version       string
	namespace     string
	name          string
	appsVersion   string
	coreVersion   string
	volume        int
	isNodePort    bool
	isLoadBalance bool
}

var createExamples = `
	# Create default KubeMQ cluster
	# kubetools cluster create b33600cc-93ef-4395-bba3-13131eb27d5e -d

	# Create KubeMQ cluster with options  
	# kubetools cluster create b3330scc-93ef-4395-bba3-13131sb2785e

	# Export KubeMQ cluster yaml file (Dry-Run)    
	# kubetools cluster create b3330scc-93ef-4395-bba3-13131sb2785e -f 
`
var createLong = `Create a KubeMQ cluster`
var createShort = `Create a KubeMQ cluster`

func NewCmdCreate(cfg *config.Config) *cobra.Command {
	o := &CreateOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "create",
		Aliases: []string{"c"},
		Short:   createShort,
		Long:    createLong,
		Example: createExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}

	cmd.PersistentFlags().BoolVarP(&o.defaultMode, "default", "d", false, "Create KubeMQ cluster with default parameters")
	cmd.PersistentFlags().BoolVarP(&o.exportFile, "file", "f", false, "Generate yaml configuration file")

	return cmd
}

func (o *CreateOptions) Complete(args []string) error {
	if len(args) > 0 {
		o.token = args[0]
	} else {
		toRegister := true
		promptConfirm := &survey.Confirm{
			Renderer: survey.Renderer{},
			Message:  "No KubeMQ token provided, want to open the registration form ?",
			Default:  true,
			Help:     "",
		}
		err := survey.AskOne(promptConfirm, &toRegister)
		if err != nil {
			return err
		}
		err = open.Run("https://account.kubemq.io/login/register")
		if err != nil {
			return err
		}
		utils.Println("")
	}
	if o.defaultMode {
		return o.setDefaultOptions()
	}

	return o.askOptions()
}

func (o *CreateOptions) Validate() error {
	if o.token == "" {
		return fmt.Errorf("no KubeMQ token provided")
	}

	return nil
}

func (o *CreateOptions) Run(ctx context.Context) error {
	if o.exportFile {
		return o.export()
	}
	deployment := &StatefulSetDeployment{
		Namespace:   nil,
		StatefulSet: nil,
		Services:    map[string]*v1.Service{},
	}
	c, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	var created bool
	deployment.Namespace, created, err = c.CheckAndCreateNamespace(o.namespace)
	if err != nil {
		return err
	}
	if created {
		utils.Printlnf("Namespace %s created", o.namespace)
	}
	spec, err := NewStatefulSetConfig(o).Spec()
	if err != nil {
		return err
	}
	isCreated := false
	deployment.StatefulSet, err = c.CreateStatefulSet(spec)
	if err != nil {
		utils.Printlnf("StatefulSet %s/%s not created. Error: %s", o.namespace, o.name, utils.Title(err.Error()))
	} else {
		isCreated = true
		utils.Printlnf("StatefulSet %s/%s created", o.namespace, o.name)
	}

	for _, cfg := range NewServiceConfigs(o) {
		spec, err := cfg.Spec()
		svc, err := c.CreateService(spec)
		if err != nil {
			utils.Printlnf("Service %s/%s not created. Error: %s", cfg.Namespace, cfg.Name, utils.Title(err.Error()))
		} else {
			if svc != nil {
				utils.Printlnf("Service %s/%s created", cfg.Namespace, cfg.Name)
				deployment.Services[svc.Name] = svc
			}
		}

	}
	if !isCreated {
		return nil
	}
	utils.Printlnf("StatefulSet %s/%s list:", o.namespace, o.name)
	done := make(chan struct{})
	evt := make(chan *appsv1.StatefulSet)
	go c.GetStatefulSetEvents(ctx, evt, done)

	for {
		select {
		case sts := <-evt:
			if int32(o.replicas) == sts.Status.Replicas && sts.Status.Replicas == sts.Status.ReadyReplicas {
				utils.Printlnf("Desired:%d Current:%d Ready:%d", o.replicas, sts.Status.Replicas, sts.Status.ReadyReplicas)
				done <- struct{}{}
				return nil
			} else {
				utils.Printlnf("Desired:%d Current:%d Ready:%d", o.replicas, sts.Status.Replicas, sts.Status.ReadyReplicas)
			}
		case <-ctx.Done():
			return nil
		}
	}
	return nil
}
func (o *CreateOptions) export() error {
	deployment := struct {
		Namespace   []byte
		StatefulSet []byte
		Services    [][]byte
	}{}
	c, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	ns, existed, err := c.GetNamespace(o.namespace)
	if err != nil {
		return err
	}
	if !existed {
		deployment.Namespace, err = yaml.Marshal(ns)
		if err != nil {
			return err
		}
	}

	deployment.StatefulSet, err = NewStatefulSetConfig(o).Spec()
	if err != nil {
		return err
	}
	for _, cfg := range NewServiceConfigs(o) {
		spec, err := cfg.Spec()
		if err != nil {
			return err
		}
		deployment.Services = append(deployment.Services, spec)
	}
	f, err := os.Create(fmt.Sprintf("%s.yaml", o.name))
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	if deployment.Namespace != nil {
		w.Write(deployment.Namespace)
		w.WriteString(fmt.Sprintf("---"))
	}
	if deployment.StatefulSet != nil {
		w.Write(deployment.StatefulSet)
		w.WriteString(fmt.Sprintf("---"))
	}

	for _, svc := range deployment.Services {
		w.Write(svc)
		w.WriteString(fmt.Sprintf("---"))
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	utils.Printlnf("create yaml file was exported to %s.yaml", o.name)
	return nil
}

func (o *CreateOptions) askOptions() error {
	answers := struct {
		Namespace string
		Name      string
		Version   string
		Replicas  int
		Volume    int
		Service   string
	}{}

	qs := []*survey.Question{
		{
			Name: "namespace",
			Prompt: &survey.Input{
				Renderer: survey.Renderer{},
				Message:  "Enter namespace of KubeMQ cluster creation:",
				Default:  "default",
				Help:     "",
			},
			Validate:  survey.Validator(conf.IsRequired()),
			Transform: nil,
		},
		{
			Name: "name",
			Prompt: &survey.Input{
				Renderer: survey.Renderer{},
				Message:  "Enter KubeMQ cluster name:",
				Default:  "kubemq-cluster",
				Help:     "",
			},
			Validate:  survey.Validator(conf.IsRequired()),
			Transform: nil,
		},
		{
			Name: "version",
			Prompt: &survey.Input{
				Renderer: survey.Renderer{},
				Message:  "Set KubeMQ image version:",
				Default:  "latest",
				Help:     "",
			},
			Validate:  survey.Validator(conf.IsRequired()),
			Transform: nil,
		},
		{
			Name: "replicas",
			Prompt: &survey.Input{
				Renderer: survey.Renderer{},
				Message:  "Set KubeMQ cluster nodes:",
				Default:  "3",
				Help:     "",
			},
			Validate:  survey.Validator(conf.IsUint()),
			Transform: nil,
		},
		{
			Name: "volume",
			Prompt: &survey.Input{
				Renderer: survey.Renderer{},
				Message:  "Set KubeMQ cluster persistence volume claim size (0 - no persistence claims):",
				Default:  "0",
				Help:     "",
			},
			Validate:  survey.Validator(conf.IsUint()),
			Transform: nil,
		},
		{
			Name: "service",
			Prompt: &survey.Select{
				Renderer: survey.Renderer{},
				Message:  "Expose services as :",
				Options:  []string{"ClusterIP", "NodePort", "LoadBalancer"},
				Default:  "ClusterIP",
				Help:     "",
			},
			Validate:  nil,
			Transform: nil,
		},
	}
	err := survey.Ask(qs, &answers)
	if err != nil {
		return err
	}
	o.appsVersion = "apps/v1"
	o.coreVersion = "v1"
	o.name = answers.Name
	o.namespace = answers.Namespace
	o.version = answers.Version
	o.replicas = answers.Replicas
	o.volume = answers.Volume
	switch answers.Service {
	case "NodePort":
		o.isNodePort = true
	case "LoadBalancer":
		o.isLoadBalance = true
	}

	return nil
}

func (o *CreateOptions) setDefaultOptions() error {

	o.appsVersion = "apps/v1"
	o.coreVersion = "v1"
	o.name = "kubemq-cluster"
	o.namespace = "default"
	o.version = "latest"
	o.replicas = 3
	o.volume = 0
	return nil
}
