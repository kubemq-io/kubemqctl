package create

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	conf "github.com/kubemq-io/kubemqctl/pkg/k8s/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/deployment"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

type CreateOptions struct {
	cfg           *config.Config
	setOptions    bool
	exportFile    bool
	file          string
	importData    string
	watch         bool
	status        bool
	optionsMenu   *conf.Menu
	deployOptions *deployment.Options
}

var createExamples = `
	# Create default KubeMQ cluster
	# kubemqctl cluster create -t b33600cc-93ef-4395-bba3-13131eb27d5e 

	# Create default KubeMQ cluster and watch events and status
	# kubemqctl cluster create -t b3d3600cc-93ef-4395-bba3-13131eb27d5e -w -s

	# Import KubeMQ cluster yaml file  
	# kubemqctl cluster create -f kubemq-cluster.yaml

	# Create KubeMQ cluster with options
	# kubemqctl cluster create -t b33d30scc-93ef-43565-bba3-13131sb2785e -o

	# Export KubeMQ cluster yaml file    
	# kubemqctl cluster create -t b3d330scc-93qf-4395-bba3-13131sb2785e -e 
`
var createLong = `Create command allows to deploy a KubeMQ cluster with configuration options`
var createShort = `Create a KubeMQ cluster command`

func NewCmdCreate(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &CreateOptions{
		cfg:           cfg,
		optionsMenu:   conf.CreateMenu,
		deployOptions: &deployment.Options{},
	}
	cmd := &cobra.Command{

		Use:     "create",
		Aliases: []string{"c", "cr"},
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

	cmd.PersistentFlags().StringVarP(&o.deployOptions.Token, "token", "t", "", "set KubeMQ Token")
	cmd.PersistentFlags().BoolVarP(&o.setOptions, "options", "o", false, "create KubeMQ cluster with options")
	cmd.PersistentFlags().BoolVarP(&o.exportFile, "export", "e", false, "generate yaml configuration file output (exporting)")
	cmd.PersistentFlags().StringVarP(&o.file, "file", "f", "", "import configuration yaml file")
	cmd.PersistentFlags().BoolVarP(&o.watch, "watch", "w", false, "stream real-time events during KubeMQ cluster Create command")
	cmd.PersistentFlags().BoolVarP(&o.status, "status", "s", false, "stream real-time status events during KubeMQ cluster Create command")

	return cmd
}

func (o *CreateOptions) Complete(args []string) error {

	if o.file != "" {
		buff, err := ioutil.ReadFile(o.file)
		if err != nil {
			return err
		}
		o.importData = string(buff)
		return nil
	}

	if o.deployOptions.Token == "" {
		if o.cfg.DefaultToken != "" {
			o.deployOptions.Token = o.cfg.DefaultToken
			utils.Printlnf("Create using save KubeMQ token: %s", o.cfg.DefaultToken)
		} else {
			return fmt.Errorf("No KubeMQ token provided")
		}

	} else {
		o.cfg.DefaultToken = o.deployOptions.Token
		_ = o.cfg.Save()
	}

	if o.setOptions {
		err := conf.CreateMenu.Run()
		if err != nil {
			return err
		}
		o.deployOptions.AppsVersion = "apps/v1"
		o.deployOptions.CoreVersion = "v1"
		o.deployOptions.Name = conf.CreateBasicOptions.Name
		o.deployOptions.Namespace = conf.CreateBasicOptions.Namespace
		o.deployOptions.Version = conf.CreateBasicOptions.Image
		o.deployOptions.Replicas = conf.CreateBasicOptions.Replicas
		o.deployOptions.Volume = conf.CreateBasicOptions.Vol
		switch conf.CreateBasicOptions.ServiceMode {
		case "NodePort":
			o.deployOptions.IsNodePort = true
		case "LoadBalancer":
			o.deployOptions.IsLoadBalance = true
		}
		return nil
	}
	return o.setDefaultOptions()
}

func (o *CreateOptions) Validate() error {
	return nil
}

func (o *CreateOptions) Run(ctx context.Context) error {
	sd, err := deployment.NewStatefulSetDeployment(o.cfg)
	if err != nil {
		return err
	}

	if o.importData != "" {
		err := sd.Import(o.importData)
		if err != nil {
			return err
		}
	} else {
		var err error
		err = sd.CreateStatefulSetDeployment(o.deployOptions, o.optionsMenu)
		if err != nil {
			return err
		}
	}

	utils.Printlnf("Create started...")

	stsName := sd.StatefulSet.Name
	stsNamespace := sd.StatefulSet.Namespace
	if o.exportFile {

		f, err := os.Create(fmt.Sprintf("%s.yaml", stsName))
		if err != nil {
			return err
		}
		err = sd.Export(f)
		if err != nil {
			utils.Printlnf("export to file %s failed", stsName)
			return err
		}
		utils.Printlnf("export to file %s.yaml completed", stsName)
		return nil
	}

	executed, err := sd.Execute(sd.StatefulSet.Name, sd.StatefulSet.Namespace)
	if err != nil {
		return err
	}
	if !executed {
		return nil
	}

	if o.watch {
		go sd.Client.PrintEvents(ctx, stsNamespace, stsName)
	}

	if o.status {
		go sd.Client.PrintStatefulSetStatus(ctx, *sd.StatefulSet.Spec.Replicas, stsNamespace, stsName)
	}
	if o.status || o.watch {
		<-ctx.Done()

	}

	return nil

}
func (o *CreateOptions) setDefaultOptions() error {

	o.deployOptions.AppsVersion = "apps/v1"
	o.deployOptions.CoreVersion = "v1"
	o.deployOptions.Name = "kubemq-cluster"
	o.deployOptions.Namespace = "default"
	o.deployOptions.Version = "latest"
	o.deployOptions.Replicas = 3
	o.deployOptions.Volume = 0
	utils.Printlnf("Create KubeMQ cluster with default options:")
	utils.Printlnf("\tKubeMQ Token: %s", o.deployOptions.Token)
	utils.Printlnf("\tCluster Name: %s", o.deployOptions.Name)
	utils.Printlnf("\tCluster Namespace: %s", o.deployOptions.Namespace)
	utils.Printlnf("\tCluster Docker Image: kubemq/kubemq:%s", o.deployOptions.Version)
	utils.Printlnf("\tCluster Replicas: %d", o.deployOptions.Replicas)
	utils.Printlnf("\tCluster PVC Size: %d", o.deployOptions.Volume)

	return nil
}
