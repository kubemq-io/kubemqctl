package create

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/deployment"
	"github.com/spf13/cobra"
	"strings"
)

type deployNodeSelectorOptions struct {
	enabled bool
	keys    map[string]string
}

func defaultNodeSelectorOptions(cmd *cobra.Command) *deployNodeSelectorOptions {
	o := &deployNodeSelectorOptions{
		enabled: false,
		keys:    map[string]string{},
	}
	cmd.PersistentFlags().BoolVarP(&o.enabled, "node-selectors-enabled", "", false, "enable statefulset node selectors configuration")
	cmd.PersistentFlags().StringToStringVarP(&o.keys, "node-selectors-keys", "", map[string]string{}, "set statefulset node selectors key-value (map)")
	return o
}

func (o *deployNodeSelectorOptions) validate() error {
	if !o.enabled {
		return nil
	}
	return nil
}

func (o *deployNodeSelectorOptions) complete() error {
	if !o.enabled {
		return nil
	}
	return nil
}

func (o *deployNodeSelectorOptions) setConfig(config *deployment.KubeMQManifestConfig) *deployNodeSelectorOptions {
	if !o.enabled {
		return nil
	}
	tmpl := []string{"      nodeSelector:\n"}
	for key, value := range o.keys {
		tmpl = append(tmpl, fmt.Sprintf("        %s:%s\n", key, value))
	}
	fmt.Println(strings.Join(tmpl, ""))
	config.StatefulSet.SetNodeSelectors(strings.Join(tmpl, ""))
	return o
}
