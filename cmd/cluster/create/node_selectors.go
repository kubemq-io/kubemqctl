package create

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/deployment"
	"github.com/spf13/cobra"
	"strings"
)

type deployNodeSelectorOptions struct {
	keys map[string]string
}

func defaultNodeSelectorOptions(cmd *cobra.Command) *deployNodeSelectorOptions {
	o := &deployNodeSelectorOptions{
		keys: map[string]string{},
	}
	cmd.PersistentFlags().StringToStringVarP(&o.keys, "node-selectors-keys", "", map[string]string{}, "set statefulset node selectors key-value (map)")
	return o
}

func (o *deployNodeSelectorOptions) validate() error {

	return nil
}

func (o *deployNodeSelectorOptions) complete() error {
	return nil
}

func (o *deployNodeSelectorOptions) setConfig(config *deployment.KubeMQManifestConfig) *deployNodeSelectorOptions {
	if len(o.keys) > 0 {
		return nil
	}
	tmpl := []string{"      nodeSelector:\n"}
	for key, value := range o.keys {
		tmpl = append(tmpl, fmt.Sprintf("        %s: %s\n", key, value))
	}
	config.StatefulSet.SetNodeSelectors(strings.Join(tmpl, ""))
	return o
}
