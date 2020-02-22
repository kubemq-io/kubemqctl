package cluster

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	"github.com/spf13/cobra"
)

var defaultImageConfig = &deployImageOptions{
	registry:   "docker.io",
	repository: "kubemq/kubemq",
	tag:        "latest",
	pullPolicy: "Always",
}

type deployImageOptions struct {
	registry   string
	repository string
	tag        string
	pullPolicy string
}

func setImageConfig(cmd *cobra.Command) *deployImageOptions {
	o := &deployImageOptions{
		registry:   "",
		repository: "",
		tag:        "",
		pullPolicy: "",
	}
	cmd.PersistentFlags().StringVarP(&o.registry, "image-registry", "", "docker.io", "set image registry")
	cmd.PersistentFlags().StringVarP(&o.repository, "image-repository", "", "kubemq/kubemq", "set image repository")
	cmd.PersistentFlags().StringVarP(&o.tag, "image-tag", "T", "latest", "set image tag")
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
		Registry:   o.registry,
		Repository: o.repository,
		Tag:        o.tag,
		PullPolicy: o.pullPolicy,
	}
	return o
}
