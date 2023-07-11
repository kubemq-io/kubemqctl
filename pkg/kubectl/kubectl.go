package kubectl

import (
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"os"
	"path/filepath"
)

type Kubectl struct {
	configPath string
	apiConfig  *clientcmdapi.Config
	rawConfig *clientcmdapi.Config
}

func NewKubectl() *Kubectl {
	return &Kubectl{}
}

func (k *Kubectl) Init(kubeConfigPath string) error {
	k.configPath = kubeConfigPath
	var err error
	k.apiConfig, err = k.loadConfig()
	if err != nil {
		return err
	}
	return nil
}
func (k *Kubectl) loadConfig() (*clientcmdapi.Config, error) {
	var kubeconfig string
	if k.configPath != "" {
		kubeconfig = k.configPath
	} else {
		if home := homeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		} else {
			return nil, fmt.Errorf("no kubeconfig available")
		}
	}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: ""}})
	apiCfg, err := clientConfig.ConfigAccess().GetStartingConfig()
	if err != nil {
		return nil, fmt.Errorf("load config error: %s", err.Error())
	}
	return apiCfg, nil
}
func (k *Kubectl) GetContextsList() ([]string, error) {
	var err error
	k.apiConfig, err = k.loadConfig()
	if err != nil {
		return nil, err
	}
	var contexts []string
	for context := range k.apiConfig.Contexts {
		contexts = append(contexts, context)
	}
	return contexts, nil
}

func (k *Kubectl) GetCurrentContext() (string, error) {
	var err error
	k.apiConfig, err = k.loadConfig()
	if err != nil {
		return "", err
	}
	return k.apiConfig.CurrentContext, nil
}

func (k *Kubectl) SetContext(context string) error {
	var err error
	k.apiConfig, err = k.loadConfig()
	if err != nil {
		return  err
	}
	rawConfig, err := k.apiConfig.
	apiConfig.CurrentContext = context
	return nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
