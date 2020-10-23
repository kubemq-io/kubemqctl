package build

import (
	"context"
	"fmt"

	"github.com/kubemq-hub/builder/connector/common"

	"github.com/kubemq-hub/builder/survey"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/cluster"

	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type buildOptions struct {
	cfg           *config.Config
	output        string
	injectOptions common.DefaultOptions
	clusters      *ClustersBuilder
	connectors    *ConnectorsBuilder
	deployer      *deploy
	resources     *resources
}

var buildExamples = `
	# Execute build Kubemq components
	kubemqctl build	
	
	# Execute build and export yaml
	kubemqctl build -o deploy.yaml
`
var buildLong = `Executes Kubemq build command`
var buildShort = `Executes Kubemq build command`

func NewCmdBuild(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &buildOptions{
		cfg:           cfg,
		output:        "",
		injectOptions: common.NewDefaultOptions(),
		clusters:      nil,
		connectors:    nil,
		deployer:      nil,
		resources:     nil,
	}
	cmd := &cobra.Command{

		Use:     "build",
		Aliases: []string{"b"},
		Short:   buildShort,
		Long:    buildLong,
		Example: buildExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))

		},
	}
	cmd.PersistentFlags().StringVarP(&o.output, "output", "o", "", "set output yaml file name")
	return cmd
}
func (o *buildOptions) Complete(args []string) error {
	o.resources = newResources()
	return nil
}

func (o *buildOptions) Validate() error {
	return nil
}

func (o *buildOptions) askBuild() error {
	for {
		val := ""
		err := survey.NewString().
			SetKind("string").
			SetName("select-component").
			SetMessage("Select Build Action").
			SetDefault("Add KubeMQ Cluster").
			SetHelp("Sets Build Action").
			SetOptions([]string{
				"Add KubeMQ Cluster",
				"Add KubeMQ Connector",
				"Deploy",
				"Exit",
			}).
			SetRequired(true).
			Render(&val)
		if err != nil {
			return err
		}
		switch val {
		case "Add KubeMQ Cluster":
			if err := o.clusters.add(); err != nil {
				return err
			}
		case "Add KubeMQ Connector":
			if err := o.connectors.add(); err != nil {
				return err
			}
		case "Deploy":
			return o.deployer.Do()
		case "Exit":
			return nil
		}

	}
}

func (o *buildOptions) Run(ctx context.Context) error {
	client, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	//if err := o.updateInjectOption(client); err != nil {
	//	return err
	//}
	o.clusters = newClustersBuilder().
		SetResources(o.resources)

	o.connectors = newConnectorsBuilder().
		SetResources(o.resources)

	o.deployer = newDeploy().
		SetClient(client).
		SetClusters(o.clusters).
		SetConnectors(o.connectors)

	return o.askBuild()
}

//func (o *buildOptions) saveDeployments() error {
//	output := strings.Join(o.deployments, "\n---\n")
//	return ioutil.WriteFile(o.output, []byte(output), 0644)
//}

func (o *buildOptions) updateInjectOption(client *client.Client) error {
	clusterManager, err := cluster.NewManager(client)
	if err != nil {
		return err
	}
	clusters, err := clusterManager.GetKubemqClusters()
	if err != nil {
		return err
	}
	var kubemqAddress []string
	for _, c := range clusters.Items {
		kubemqAddress = append(kubemqAddress, fmt.Sprintf("%s-grpc.%s.svc.local", c.Name, c.Namespace))
	}

	o.injectOptions.Add("kubemq-address", kubemqAddress)
	nsList, err := client.GetNamespaceList()
	if err != nil {
		return nil
	}
	o.injectOptions.Add("namespaces", nsList)
	return nil
}
