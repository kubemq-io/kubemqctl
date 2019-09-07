package create

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

type CreateOptions struct {
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

var createExamples = `
	# Create default KubeMQ cluster
	# kubetools cluster create b33600cc-93ef-4395-bba3-13131eb27d5e

	# Create KubeMQ cluster default namespace with specific cluster name  
	# kubetools cluster create b3330scc-93ef-4395-bba3-13131sb2785e -n kubemq-cluster-1

	# Create KubeMQ cluster with specific cluster name and namespace   
	# kubetools cluster create b3330scc-93ef-4395-bba3-13131sb2785e -n kubemq-cluster-1 -s kubemq-namespace

	# Create default KubeMQ cluster with 5 pods   
	# kubetools cluster create b3330scc-93ef-4395-bba3-13131sb2785e -r 5

	# Create default KubeMQ cluster with persistence volume claims of 10Gi   
	# kubetools cluster create b3330scc-93ef-4395-bba3-13131sb2785e -v 10

	# Create default KubeMQ cluster with specific KubeMQ image version   
	# kubetools cluster create b3330scc-93ef-4395-bba3-13131sb2785e -i v1.6.2
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

	cmd.PersistentFlags().StringVarP(&o.namespace, "namespace", "s", "default", "Set namespace name")
	cmd.PersistentFlags().StringVarP(&o.name, "name", "n", "kubemq-cluster", "Set KubeMQ cluster name")
	cmd.PersistentFlags().StringVarP(&o.version, "image", "i", "latest", "Set KubeMQ image version")
	cmd.PersistentFlags().StringVarP(&o.token, "token", "t", "", "Set KubeMQ token")
	cmd.PersistentFlags().IntVarP(&o.replicas, "replicas", "r", 3, "Set how many replicas in KubeMQ cluster")
	cmd.PersistentFlags().IntVarP(&o.volume, "volume", "v", 0, "Set size of persistence volume")
	cmd.PersistentFlags().StringVarP(&o.appsVersion, "apps-api", "", "apps/v1", "Set api version for kubernetes apps end-point")
	cmd.PersistentFlags().StringVarP(&o.coreVersion, "core-api", "", "v1", "Set api version for kubernetes core end-point")
	cmd.PersistentFlags().BoolVarP(&o.isNodePort, "set-node-port", "", false, "Set expose services with NodePort")
	cmd.PersistentFlags().BoolVarP(&o.isLoadBalance, "set-load-balancer", "", false, "Set expose services with LoadBalancer")

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
	return nil
}

func (o *CreateOptions) Validate() error {
	if o.token == "" {
		return fmt.Errorf("no KubeMQ token provided")
	}

	return nil
}

func (o *CreateOptions) Run(ctx context.Context) error {
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
	isCreateed := false
	deployment.StatefulSet, err = c.CreateStatefulSet(spec)
	if err != nil {
		utils.Printlnf("StatefulSet %s/%s not created. Error: %s", o.namespace, o.name, utils.Title(err.Error()))
	} else {
		isCreateed = true
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
	if !isCreateed {
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
