package config

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	apiv1 "k8s.io/api/core/v1"
)

type MenuItem struct {
	Label   string
	Action  *EntryGroup
	SubMenu *Menu
}

func (mi *MenuItem) Run() error {
	if mi.SubMenu != nil {
		return mi.SubMenu.Run()
	}
	if mi.Action != nil {
		return mi.Action.Execute()
	}
	return nil
}

type Menu struct {
	Prefix string
	Items  []*MenuItem
}

func (m *Menu) Run() error {
	items := []string{}
	itemsMap := make(map[string]*MenuItem)
	for _, item := range m.Items {
		items = append(items, item.Label)
		itemsMap[item.Label] = item
	}
	selected := []string{}
	prompt := &survey.MultiSelect{
		Renderer:      survey.Renderer{},
		Message:       fmt.Sprintf("(%s) Please select items", m.Prefix),
		Options:       items,
		Help:          "",
		PageSize:      len(m.Items),
		VimMode:       false,
		FilterMessage: "",
		Filter:        nil,
	}
	err := survey.AskOne(prompt, &selected)
	if err != nil {
		return err
	}

	for _, selectedItem := range selected {
		menuItem := itemsMap[selectedItem]
		err := menuItem.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Menu) ExportEnvVar() []apiv1.EnvVar {
	var envVars []apiv1.EnvVar
	for _, item := range m.Items {
		if item.SubMenu != nil {
			envVars = append(envVars, item.SubMenu.ExportEnvVar()...)
		}
		if item.Action != nil {
			envVars = append(envVars, item.Action.ExportEnvVar()...)
		}
	}
	return envVars
}

func (m *Menu) ExportVolumeMounts() []apiv1.VolumeMount {
	var volMounts []apiv1.VolumeMount
	for _, item := range m.Items {
		if item.SubMenu != nil {
			volMounts = append(volMounts, item.SubMenu.ExportVolumeMounts()...)
		}
		if item.Action != nil {
			volMounts = append(volMounts, item.Action.ExportVolumeMounts()...)
		}
	}
	return volMounts
}

func (m *Menu) ExportVolumes() []apiv1.Volume {
	var vol []apiv1.Volume
	for _, item := range m.Items {
		if item.SubMenu != nil {
			vol = append(vol, item.SubMenu.ExportVolumes()...)
		}
		if item.Action != nil {
			vol = append(vol, item.Action.ExportVolumes()...)
		}
	}
	return vol
}

func (m *Menu) ExportConfigMaps() []ConfigMap {
	var cm []ConfigMap
	for _, item := range m.Items {
		if item.SubMenu != nil {
			cm = append(cm, item.SubMenu.ExportConfigMaps()...)
		}
		if item.Action != nil {
			cm = append(cm, item.Action.ExportConfigMaps()...)
		}
	}
	return cm
}

func (m *Menu) ExportSecrets() []Secret {
	var cm []Secret
	for _, item := range m.Items {
		if item.SubMenu != nil {
			cm = append(cm, item.SubMenu.ExportSecrets()...)
		}
		if item.Action != nil {
			cm = append(cm, item.Action.ExportSecrets()...)
		}
	}
	return cm
}
