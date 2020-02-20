package create

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/crd/cluster"
	"github.com/spf13/cobra"
)

type deployImageOptions struct {
	registry   string
	repository string
	tag        string
	pullPolicy string
}

func defaultImageConfig(cmd *cobra.Command) *deployImageOptions {
	o := &deployImageOptions{
		registry:   "",
		repository: "",
		tag:        "",
		pullPolicy: "",
	}
	cmd.PersistentFlags().StringVarP(&o.registry, "image-registry", "", "docker.io", "set image registry")
	cmd.PersistentFlags().StringVarP(&o.registry, "image-repository", "", "kubemq/kubemq", "set image repository")
	cmd.PersistentFlags().StringVarP(&o.registry, "image-tag", "", "latest", "set image tag")
	cmd.PersistentFlags().StringVarP(&o.registry, "image-pull-policy", "", "Always", "set image pull policy")
	return o
}

func (o *deployImageOptions) validate() error {

	return nil
}
func (o *deployImageOptions) complete() error {
	return nil
}

func (o *deployImageOptions) setConfig(deployment *cluster.KubemqCluster) *deployImageOptions {
	deployment.Spec.Image = &cluster.ImageConfig{
		Registry:   o.registry,
		Repository: o.repository,
		Tag:        o.tag,
		PullPolicy: o.pullPolicy,
	}
	return o
}
