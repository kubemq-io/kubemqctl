package config

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	apiv1 "k8s.io/api/core/v1"
)

type SecretEntry struct {
	Name        string
	Description string
	ClusterName string
	EnvVarName  string
	FilePath    string
	FileName    string
	SecretName  string
	SecretValue string
	VolumeName  string
	SecretType  string
	useExisted  bool
}

func (sec *SecretEntry) Execute() error {
	sec.SecretName = fmt.Sprintf("%s-%s", sec.ClusterName, sec.Name)
	sec.VolumeName = fmt.Sprintf("%s-%s-vol", sec.ClusterName, sec.Name)
	existed := fmt.Sprintf("Use Existing Secret - %s", sec.Name)
	new := "New Secret"
	selected := ""
	selectOptions := &survey.Select{
		Renderer:      survey.Renderer{},
		Message:       fmt.Sprintf("Please Select %s source:", sec.Description),
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
		sec.useExisted = true
		return nil
	}
	editor := &Editor{
		Message:    fmt.Sprintf("Enter or Copy/Paste %s data:", sec.Description),
		Validators: nil,
		Default:    "",
		Help:       "",
	}
	err = editor.Ask(&sec.SecretValue)
	if err != nil {
		return err
	}

	return nil
}

func (sec *SecretEntry) EnvVar() *EnvVar {
	return &EnvVar{
		Name:  sec.EnvVarName,
		Value: sec.FilePath,
	}
}

func (sec *SecretEntry) Volume() *Volume {
	return &Volume{
		Volume: &apiv1.Volume{
			Name: sec.VolumeName,
			VolumeSource: apiv1.VolumeSource{
				Secret: &apiv1.SecretVolumeSource{
					SecretName:  sec.SecretName,
					Items:       nil,
					DefaultMode: nil,
					Optional:    nil,
				},
			}},
		Mount: &apiv1.VolumeMount{
			Name:             sec.VolumeName,
			ReadOnly:         false,
			MountPath:        sec.FilePath,
			SubPath:          sec.FileName,
			MountPropagation: nil,
			SubPathExpr:      "",
		},
	}
}

func (sec *SecretEntry) ConfigMap() *ConfigMap {
	return nil
}

func (sec *SecretEntry) Secret() *Secret {
	if sec.useExisted {
		return nil
	}

	return &Secret{
		Name:       sec.SecretName,
		Value:      sec.SecretValue,
		FileName:   sec.FileName,
		SecretType: sec.SecretType,
	}
}
