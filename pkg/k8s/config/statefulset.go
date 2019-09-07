package config

import "github.com/AlecAivazis/survey/v2"

var qs = []*survey.Question{
	{
		Name: "namespace",
		Prompt: &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "Enter namespace of KubeMQ cluster creation:",
			Default:  "default",
			Help:     "",
		},
		Validate:  survey.Validator(IsRequired()),
		Transform: nil,
	},
	{
		Name: "name",
		Prompt: &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "Enter KubeMQ cluster name:",
			Default:  "kubemq-cluster",
			Help:     "",
		},
		Validate:  survey.Validator(IsRequired()),
		Transform: nil,
	},
	{
		Name: "version",
		Prompt: &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "Set KubeMQ image version:",
			Default:  "latest",
			Help:     "",
		},
		Validate:  survey.Validator(IsRequired()),
		Transform: nil,
	},
	{
		Name: "replicas",
		Prompt: &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "Set KubeMQ cluster nodes:",
			Default:  "3",
			Help:     "",
		},
		Validate:  survey.Validator(IsUint()),
		Transform: nil,
	},
	{
		Name: "volume",
		Prompt: &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "Set KubeMQ cluster persistence volume claim size (0 - no persistence claims):",
			Default:  "0",
			Help:     "",
		},
		Validate:  survey.Validator(IsUint()),
		Transform: nil,
	},
}

type StatefulSetConfig struct {
	ApiVersion string
	Name       string
	Namespace  string
	Replicas   int
	Token      string
	Version    string
	Volume     int
}

func NewDefaultStatefulSetConfig(token string) StatefulSetConfig {
	return StatefulSetConfig{
		ApiVersion: "apps/v1",
		Name:       "kubemq-cluster",
		Namespace:  "default",
		Replicas:   3,
		Token:      "kubemq-cluster",
		Version:    "latest",
		Volume:     0,
	}
}

func NewStatefulSetConfigFromQuestions(token string) StatefulSetConfig {
	config := StatefulSetConfig{
		ApiVersion: "",
		Name:       "",
		Namespace:  "",
		Replicas:   0,
		Token:      token,
		Version:    "",
		Volume:     0,
	}
	err := survey.Ask(qs, &config)
	if err != nil {
		return NewDefaultStatefulSetConfig(token)
	}
	return config
}

func (s StatefulSetConfig) Spec() ([]byte, error) {
	t := NewTemplate(defaultStsTemplate, s)
	return t.Get()
}
