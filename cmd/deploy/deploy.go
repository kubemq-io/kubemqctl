package deploy

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s/client"
	"github.com/skratchdot/open-golang/open"
	"k8s.io/api/core/v1"

	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
)

type DeployOptions struct {
	cfg           *config.Config
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

var deployExamples = `
	# Deploy default KubeMQ cluster
	# kubetools deploy b33600cc-93ef-4395-bba3-13131eb27d5e

	# Deploy KubeMQ cluster default namespace with specific cluster name  
	# kubetools deploy b3330scc-93ef-4395-bba3-13131sb2785e -n kubemq-cluster-1

	# Deploy KubeMQ cluster with specific cluster name and namespace   
	# kubetools deploy b3330scc-93ef-4395-bba3-13131sb2785e -n kubemq-cluster-1 -s kubemq-namespace

	# Deploy default KubeMQ cluster with 5 pods   
	# kubetools deploy b3330scc-93ef-4395-bba3-13131sb2785e -r 5

	# Deploy default KubeMQ cluster with persistence volume claims of 10Gi   
	# kubetools deploy b3330scc-93ef-4395-bba3-13131sb2785e -v 10

	# Deploy default KubeMQ cluster with specific KubeMQ image version   
	# kubetools deploy b3330scc-93ef-4395-bba3-13131sb2785e -i v1.6.2
`
var deployLong = `Deploy KubeMQ cluster`
var deployShort = `Deploy KubeMQ cluster`

func NewCmdDeploy(cfg *config.Config) *cobra.Command {
	o := &DeployOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "deploy",
		Aliases: []string{"dep", "dp", "d"},
		Short:   deployShort,
		Long:    deployLong,
		Example: deployExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args))
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}

	cmd.PersistentFlags().StringVarP(&o.namespace, "namespace", "s", "default", "set namespace name")
	cmd.PersistentFlags().StringVarP(&o.name, "name", "n", "kubemq-cluster", "set KubeMQ cluster name")
	cmd.PersistentFlags().StringVarP(&o.version, "image", "i", "latest", "set KubeMQ image version")
	cmd.PersistentFlags().StringVarP(&o.token, "token", "t", "", "set KubeMQ token")
	cmd.PersistentFlags().IntVarP(&o.replicas, "replicas", "r", 3, "set how many replicas in KubeMQ cluster")
	cmd.PersistentFlags().IntVarP(&o.volume, "volume", "v", 0, "set size of persistence volume")
	cmd.PersistentFlags().StringVarP(&o.appsVersion, "apps-api", "", "apps/v1", "set api version for kubernetes apps end-point")
	cmd.PersistentFlags().StringVarP(&o.coreVersion, "core-api", "", "v1", "set api version for kubernetes core end-point")
	cmd.PersistentFlags().BoolVarP(&o.isNodePort, "set-node-port", "", false, "set expose services with NodePort")
	cmd.PersistentFlags().BoolVarP(&o.isLoadBalance, "set-load-balancer", "", false, "set expose services with LoadBalancer")

	return cmd
}

func (o *DeployOptions) Complete(args []string) error {
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
	return nil
}

func (o *DeployOptions) Validate() error {
	if o.token == "" {
		return fmt.Errorf("no KubeMQ token provided")
	}

	return nil
}

func (o *DeployOptions) Run(ctx context.Context) error {
	deployment := &StatefulSetDeployment{
		Namespace:   nil,
		StatefulSet: nil,
		Services:    map[string]*v1.Service{},
	}
	c, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}

	//stsCfg := StatefulSetConfig{
	//	ApiVersion: "apps/v1",
	//	Name:       "kubemq-cluster",
	//	Namespace:  "default",
	//	Replicas:   3,
	//	Token:      "b33600cc-93ef-4395-bba3-13131eb2785e",
	//	Version:    "latest",
	//}
	//t := NewTemplate(defaultStsTemplate, stsCfg)
	//spec, err := t.Get()
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
	isDeployed := false
	deployment.StatefulSet, err = c.CreateStatefulSet(spec)
	if err != nil {
		utils.Printlnf("StatefulSet %s/%s not deployed. Error: %s", o.namespace, o.name, utils.Title(err.Error()))
	} else {
		isDeployed = true
		utils.Printlnf("StatefulSet %s/%s deployed", o.namespace, o.name)
	}

	for _, cfg := range NewServiceConfigs(o) {
		spec, err := cfg.Spec()
		svc, err := c.CreateService(spec)
		if err != nil {
			utils.Printlnf("Service %s/%s not deployed. Error: %s", cfg.Namespace, cfg.Name, utils.Title(err.Error()))
		} else {
			if svc != nil {
				utils.Printlnf("Service %s/%s deployed", cfg.Namespace, cfg.Name)
				deployment.Services[svc.Name] = svc
			}
		}

	}
	if !isDeployed {
		return nil
	}
	utils.Printlnf("StatefulSet %s/%s status:", o.namespace, o.name)
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
