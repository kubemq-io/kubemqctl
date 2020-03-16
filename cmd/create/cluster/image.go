package cluster

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	"github.com/spf13/cobra"
)

var defaultImageConfig = &deployImageOptions{
	image:      "docker.io",
	pullPolicy: "Always",
}

type deployImageOptions struct {
	image      string
	pullPolicy string
}

func setImageConfig(cmd *cobra.Command) *deployImageOptions {
	o := &deployImageOptions{
		image:      "",
		pullPolicy: "",
	}
	cmd.PersistentFlags().StringVarP(&o.image, "image", "", "docker.io/kubemq/kubemq:latest", "set image registry/repository:tag")
	cmd.PersistentFlags().StringVarP(&o.pullPolicy, "image-pull-policy", "", "Always", "set image pull policy")
	return o
}

func (o *deployImageOptions) validate() error {

	return nil
}
func (o *deployImageOptions) complete() error {
	return nil
}

func (o *deployImageOptions) setConfig(deployment *kubemqcluster.KubemqCluster) *deployImageOptions {
	if isDefault(o, defaultImageConfig) {
		return o
	}
	deployment.Spec.Image = &kubemqcluster.ImageConfig{
		Image:      o.image,
		PullPolicy: o.pullPolicy,
	}
	return o
}
