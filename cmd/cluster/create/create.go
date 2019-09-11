package create

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubetools/pkg/config"
	conf "github.com/kubemq-io/kubetools/pkg/k8s/config"
	"os"

	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
)

type CreateOptions struct {
	cfg           *config.Config
	setOptions    bool
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
	optionsMenu   *conf.Menu
	deployment    *StatefulSetDeployment
}

var createExamples = `
	# Create default KubeMQ cluster
	# kubetools cluster create -t b33600cc-93ef-4395-bba3-13131eb27d5e 

	# Create KubeMQ cluster with options  
	# kubetools cluster create -t b3330scc-93ef-4395-bba3-13131sb2785e -o

	# Export KubeMQ cluster yaml file (Dry-Run)    
	# kubetools cluster create -t b3330scc-93ef-4395-bba3-13131sb2785e -f 
`
var createLong = `Create a KubeMQ cluster`
var createShort = `Create a KubeMQ cluster`

func NewCmdCreate(cfg *config.Config) *cobra.Command {
	o := &CreateOptions{
		cfg:         cfg,
		optionsMenu: conf.CreateMenu,
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

	cmd.PersistentFlags().StringVarP(&o.token, "token", "t", "", "Set KubeMQ Token")
	cmd.PersistentFlags().BoolVarP(&o.setOptions, "options", "o", false, "Create KubeMQ cluster with options")
	cmd.PersistentFlags().BoolVarP(&o.exportFile, "file", "f", false, "Generate yaml configuration file")

	return cmd
}

func (o *CreateOptions) Complete(args []string) error {
	if o.token == "" {
		return fmt.Errorf("No KubeMQ token provided")
	}
	if o.setOptions {
		err := conf.CreateMenu.Run()
		if err != nil {
			return err
		}
		o.appsVersion = "apps/v1"
		o.coreVersion = "v1"
		o.name = conf.CreateBasicOptions.Name
		o.namespace = conf.CreateBasicOptions.Namespace
		o.version = conf.CreateBasicOptions.Image
		o.replicas = conf.CreateBasicOptions.Replicas
		o.volume = conf.CreateBasicOptions.Vol
		switch conf.CreateBasicOptions.ServiceMode {
		case "NodePort":
			o.isNodePort = true
		case "LoadBalancer":
			o.isLoadBalance = true
		}
		return nil
	}
	return o.setDefaultOptions()
}

func (o *CreateOptions) Validate() error {
	if o.token == "" {
		return fmt.Errorf("no KubeMQ token provided")
	}

	return nil
}

func (o *CreateOptions) Run(ctx context.Context) error {
	utils.Printlnf("\n")
	sd, err := CreateStatefulSetDeployment(o)
	if err != nil {
		return err
	}
	if o.exportFile {

		f, err := os.Create(fmt.Sprintf("%s.yaml", o.name))
		if err != nil {
			return err
		}
		return sd.Export(f, o)
	}

	executed, err := sd.Execute(o)
	if err != nil {
		return err
	}
	if !executed {
		return nil
	}
	utils.Printlnf("Create StatefulSet %s/%s progress:", o.namespace, o.name)
	done := make(chan struct{})
	evt := make(chan *appsv1.StatefulSet)
	go sd.client.GetStatefulSetEvents(ctx, evt, done)

	for {
		select {
		case sts := <-evt:
			if sts.Name == sd.StatefulSet.Name && sts.Namespace == sd.StatefulSet.Namespace {
				if int32(o.replicas) == sts.Status.Replicas && sts.Status.Replicas == sts.Status.ReadyReplicas {
					utils.Printlnf("Desired:%d Current:%d Ready:%d", o.replicas, sts.Status.Replicas, sts.Status.ReadyReplicas)
					done <- struct{}{}
					return nil
				} else {
					utils.Printlnf("Desired:%d Current:%d Ready:%d", o.replicas, sts.Status.Replicas, sts.Status.ReadyReplicas)
				}
			}
		case <-ctx.Done():
			return nil
		}
	}

}
func (o *CreateOptions) setDefaultOptions() error {

	o.appsVersion = "apps/v1"
	o.coreVersion = "v1"
	o.name = "kubemq-cluster"
	o.namespace = "default"
	o.version = "latest"
	o.replicas = 3
	o.volume = 0
	utils.Printlnf("Create KubeMQ cluster with default options:")
	utils.Printlnf("\tKubeMQ Token: %s", o.token)
	utils.Printlnf("\tCluster Name: %s", o.name)
	utils.Printlnf("\tCluster Namespace: %s", o.namespace)
	utils.Printlnf("\tCluster Docker Image: kubemq/kubemq:%s", o.version)
	utils.Printlnf("\tCluster Replicas: %d", o.replicas)
	utils.Printlnf("\tCluster PVC Size: %d", o.volume)

	return nil
}
