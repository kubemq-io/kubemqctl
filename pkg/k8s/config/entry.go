package config

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	apiv1 "k8s.io/api/core/v1"
)

var entryTemplate = `
            - name: %s
              value: %s
`

type Entry struct {
	VarName  string
	VarValue string
	Prompt   Prompt
}

func (e *Entry) String() string {
	if e.VarValue == "" {
		return ""
	}
	return fmt.Sprintf(entryTemplate, e.VarName, e.VarValue)
}

func (e *Entry) Execute() error {
	if e.Prompt == nil {
		return nil
	}
	for {
		err := e.Prompt.Ask(&e.VarValue)
		if err == nil {
			return nil
		}
	}
}

type EntryGroup struct {
	Name      string
	Entries   []*Entry
	SubGroups []*EntryGroup
	Result    map[string]*Entry
}

func (eg *EntryGroup) Execute() error {
	confirm := &survey.Confirm{
		Renderer: survey.Renderer{},
		Message:  fmt.Sprintf("Would you like to configure %s ?", eg.Name),
		Default:  false,
		Help:     "",
	}
	resp := false
	err := survey.AskOne(confirm, &resp)
	if err != nil {
		return err
	}
	if !resp {
		return nil
	}
	eg.Result = map[string]*Entry{}

	for _, entry := range eg.Entries {
		err := entry.Execute()
		if err != nil {
			return err
		}
		if entry.VarValue != "" {
			eg.Result[entry.VarName] = entry
		}
	}

	for _, group := range eg.SubGroups {
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
			Name:      entry.VarName,
			Value:     entry.VarValue,
			ValueFrom: nil,
		})
	}
	return envVars
}
func (eg *EntryGroup) getEntry(name string) *Entry {
	for _, entry := range eg.Entries {
		if entry.VarName == name {
			return entry
		}
	}
	for _, group := range eg.SubGroups {
		entry := group.getEntry(name)
		if entry != nil {
			return entry
		}

	}
	return nil
}

func (eg *EntryGroup) LoadEnvVar(envVars []apiv1.EnvVar) {
	for _, env := range envVars {
		if entry := eg.getEntry(env.Name); entry != nil {
			entry.VarValue = env.Value
			entry.Prompt.SetDefault(env.Value)
		}
	}

}
