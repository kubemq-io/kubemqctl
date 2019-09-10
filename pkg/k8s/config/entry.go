package config

import (
	"github.com/AlecAivazis/survey/v2"
	apiv1 "k8s.io/api/core/v1"
)

type EnvVar struct {
	Name  string
	Value string
}

type ConfigMap struct {
	Name     string
	Value    string
	FileName string
}

type Secret struct {
	Name       string
	Value      string
	FileName   string
	SecretType string
}

type Volume struct {
	Volume *apiv1.Volume
	Mount  *apiv1.VolumeMount
}

type Entry interface {
	Execute() error
	EnvVar() *EnvVar
	Volume() *Volume
	ConfigMap() *ConfigMap
	Secret() *Secret
}

type EntryGroup struct {
	Name      string
	Entries   []Entry
	SubGroups []*EntryGroup
	Result    map[string]Entry
}

func (eg *EntryGroup) Execute() error {
	eg.Result = map[string]Entry{}
	var err error

	for _, entry := range eg.Entries {
		err := entry.Execute()
		if err != nil {
			return err
		}
		if entry.EnvVar() != nil && entry.EnvVar().Value != "" {
			eg.Result[entry.EnvVar().Value] = entry
		}
	}

	var multiSelectOptions []string

	for _, group := range eg.SubGroups {
		multiSelectOptions = append(multiSelectOptions, group.Name)
	}
	if multiSelectOptions == nil {
		return nil
	}
	multiSelect := &survey.MultiSelect{
		Renderer:      survey.Renderer{},
		Message:       "Select what would you like to configure ?:",
		Options:       multiSelectOptions,
		Default:       false,
		Help:          "",
		PageSize:      0,
		VimMode:       false,
		FilterMessage: "",
		Filter:        nil,
	}
	multiSelectEntries := []string{}
	err = survey.AskOne(multiSelect, &multiSelectEntries)
	if err != nil {
		return err
	}
	if multiSelectEntries == nil {
		return nil
	}
	multiSelectMap := make(map[string]struct{})
	for _, entry := range multiSelectEntries {
		multiSelectMap[entry] = struct{}{}
	}
	for _, group := range eg.SubGroups {
		_, ok := multiSelectMap[group.Name]
		if !ok {
			continue
		}
		err := group.Execute()
		if err != nil {
			return err
		}
		for name, entry := range group.Result {
			eg.Result[name] = entry
		}
	}
	return nil
}

func (eg *EntryGroup) ExportEnvVar() []apiv1.EnvVar {
	var envVars []apiv1.EnvVar
	for _, entry := range eg.Result {
		envVars = append(envVars, apiv1.EnvVar{
			Name:      entry.EnvVar().Name,
			Value:     entry.EnvVar().Value,
			ValueFrom: nil,
		})
	}
	return envVars
}

func (eg *EntryGroup) ExportVolumeMounts() []apiv1.VolumeMount {
	var volMounts []apiv1.VolumeMount
	for _, entry := range eg.Result {

		if entry.Volume() != nil {
			volMounts = append(volMounts, *entry.Volume().Mount)
		}
	}
	return volMounts
}

func (eg *EntryGroup) ExportVolumes() []apiv1.Volume {
	var vol []apiv1.Volume
	for _, entry := range eg.Result {
		if entry.Volume() != nil {
			vol = append(vol, *entry.Volume().Volume)
		}
	}
	return vol
}

func (eg *EntryGroup) ExportConfigMaps() []ConfigMap {
	var cm []ConfigMap
	for _, entry := range eg.Result {
		if entry.ConfigMap() != nil {
			cm = append(cm, *entry.ConfigMap())
		}

	}
	return cm
}

func (eg *EntryGroup) ExportSecrets() []Secret {
	var cm []Secret
	for _, entry := range eg.Result {
		if entry.Secret() != nil {
			cm = append(cm, *entry.Secret())
		}
	}
	return cm
}
