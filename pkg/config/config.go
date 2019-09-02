package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var DefaultConfig = &Config{
	AutoIntegrated:     true,
	CurrentNamespace:   "default",
	CurrentStatefulSet: "kubemq-cluster",
	Host:               "localhost",
	GrpcPort:           50000,
	RestPort:           9090,
	ApiPort:            8080,
	IsSecured:          false,
	CertFile:           "",
	KubeConfigPath:     "",
	ConnectionType:     "grpc",
}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
func CheckConfigFile() (*Config, error) {
	if !exists(".kubetools.yaml") {
		data, err := yaml.Marshal(DefaultConfig)
		if err != nil {
			return DefaultConfig, err
		}
		err = ioutil.WriteFile(".kubetools.yaml", data, 0644)
		if err != nil {
			return DefaultConfig, err
		}
	}
	return nil, nil
}

type Config struct {
	AutoIntegrated     bool
	CurrentNamespace   string
	CurrentStatefulSet string
	Host               string
	GrpcPort           int
	RestPort           int
	ApiPort            int
	IsSecured          bool
	CertFile           string
	KubeConfigPath     string
	ConnectionType     string
}

func (c *Config) Save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(".kubetools.yaml", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) GetGRPCHostPort() (string, int) {
	return c.Host, c.GrpcPort
}

func (c *Config) GetRestHttpURI() string {
	if c.IsSecured {
		return fmt.Sprintf("https://%s:%s", c.Host, c.RestPort)
	} else {
		return fmt.Sprintf("http://%s:%s", c.Host, c.RestPort)
	}

}
func (c *Config) GetRestWsURI() string {
	if c.IsSecured {
		return fmt.Sprintf("wss://%s:%d", c.Host, c.RestPort)
	} else {
		return fmt.Sprintf("ws://%s:%d", c.Host, c.RestPort)
	}
}

func (c *Config) GetApiHttpURI() string {
	if c.IsSecured {
		return fmt.Sprintf("https://%s:%d", c.Host, c.ApiPort)
	} else {
		return fmt.Sprintf("http://%s:%d", c.Host, c.ApiPort)
	}

}
func (c *Config) GetApiWsURI() string {
	if c.IsSecured {
		return fmt.Sprintf("wss://%s:%d", c.Host, c.ApiPort)
	} else {
		return fmt.Sprintf("ws://%s:%d", c.Host, c.ApiPort)
	}
}
