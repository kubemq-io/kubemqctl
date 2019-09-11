package config

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	apiv1 "k8s.io/api/core/v1"
)

type ConfigMapEntry struct {
	Name           string
	Description    string
	ClusterName    string
	EnvVarName     string
	FilePath       string
	FileName       string
	ConfigMapName  string
	ConfigMapValue string
	VolumeName     string
	useExisted     bool
}

func (cm *ConfigMapEntry) Execute() error {
	cm.ConfigMapName = fmt.Sprintf("%s-%s", cm.ClusterName, cm.Name)
	cm.VolumeName = fmt.Sprintf("%s-%s-vol", cm.ClusterName, cm.Name)

	existed := fmt.Sprintf("Use Existing ConfigMap - %s", cm.Name)
	new := "New ConfigMap"
	selected := ""
	selectOptions := &survey.Select{
		Renderer:      survey.Renderer{},
		Message:       fmt.Sprintf("Please Select %s source:", cm.Description),
		Options:       []string{existed, new},
		Default:       existed,
		Help:          "",
		PageSize:      0,
		VimMode:       false,
		FilterMessage: "",
		Filter:        nil,
	}

	err := survey.AskOne(selectOptions, &selected)
	if err != nil {
		return err
	}
	if selected == existed {
		cm.useExisted = true
		return nil
	}

	editor := &Editor{
		Message:    fmt.Sprintf("Enter or Copy/Paste %s data:", cm.Description),
		Validators: nil,
		Default:    "",
		Help:       "",
	}
	err = editor.Ask(&cm.ConfigMapValue)
	if err != nil {
		return err
	}

	return nil
}

func (cm *ConfigMapEntry) EnvVar() *EnvVar {
	return &EnvVar{
		Name:  cm.EnvVarName,
		Value: cm.FilePath,
	}
}

func (cm *ConfigMapEntry) Volume() *Volume {
	return &Volume{
		Volume: &apiv1.Volume{
			Name: cm.VolumeName,
			VolumeSource: apiv1.VolumeSource{
				ConfigMap: &apiv1.ConfigMapVolumeSource{
					LocalObjectReference: apiv1.LocalObjectReference{
						Name: cm.ConfigMapName},
					Items:       nil,
					DefaultMode: nil,
					Optional:    nil,
				},
			}},
		Mount: &apiv1.VolumeMount{
			Name:             cm.VolumeName,
			ReadOnly:         false,
			MountPath:        cm.FilePath,
			SubPath:          cm.FileName,
			MountPropagation: nil,
			SubPathExpr:      "",
		},
	}
}
func (cm *ConfigMapEntry) ConfigMap() *ConfigMap {
	if cm.useExisted {
		return nil
	}
	return &ConfigMap{
		Name:     cm.ConfigMapName,
		Value:    cm.ConfigMapValue,
		FileName: cm.FileName,
	}
}
func (cm *ConfigMapEntry) Secret() *Secret {
	return nil
}
