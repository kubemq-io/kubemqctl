package apply

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s/deployment"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
)

type ApplyOptions struct {
	cfg        *config.Config
	importData string
}

var applyExamples = `
	# Apply KubeMQ cluster deployment
	# kubetools cluster apply kubemq-cluster.yaml 
`
var applyLong = `Apply a KubeMQ cluster`
var applyShort = `Apply a KubeMQ cluster`

func NewCmdApply(cfg *config.Config) *cobra.Command {
	o := &ApplyOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "apply",
		Aliases: []string{"a", "ap"},
		Short:   applyShort,
		Long:    applyLong,
		Example: applyExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}

	return cmd
}

func (o *ApplyOptions) Complete(args []string) error {
	if len(args) != 0 {
		buff, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}
		o.importData = string(buff)
		return nil
	}
	return fmt.Errorf("invalid argument, no yaml file specified")
}

func (o *ApplyOptions) Validate() error {
	if o.importData == "" {
		return fmt.Errorf("configuration file empty")
	}
	return nil
}

func (o *ApplyOptions) Run(ctx context.Context) error {
	sd, err := deployment.NewStatefulSetDeployment(o.cfg)
	if err != nil {
		return err
	}
	err = sd.Import(o.importData)
	if err != nil {
		return err
	}
	utils.Printlnf("Apply started...")

	executed, err := sd.Execute(sd.StatefulSet.Name, sd.StatefulSet.Namespace)
	if err != nil {
		return err
	}
	if !executed {
		return nil
	}
	utils.Printlnf("Apply StatefulSet %s/%s progress:", sd.StatefulSet.Namespace, sd.StatefulSet.Name)
	done := make(chan struct{})
	evt := make(chan *appsv1.StatefulSet)
	go sd.Client.GetStatefulSetEvents(ctx, evt, done)

	for {
		select {
		case sts := <-evt:
			if sts.Name == sd.StatefulSet.Name && sts.Namespace == sd.StatefulSet.Namespace {
				rep := *sd.StatefulSet.Spec.Replicas
				if rep == sts.Status.Replicas && sts.Status.Replicas == sts.Status.ReadyReplicas {
					utils.Printlnf("Desired:%d Current:%d Ready:%d", rep, sts.Status.Replicas, sts.Status.ReadyReplicas)
					done <- struct{}{}
					return nil
				} else {
					utils.Printlnf("Desired:%d Current:%d Ready:%d", rep, sts.Status.Replicas, sts.Status.ReadyReplicas)
				}
			}
		case <-ctx.Done():
			return nil
		}
	}

}
