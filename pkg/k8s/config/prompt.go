package config

import (
	"github.com/AlecAivazis/survey/v2"
)

type Prompt interface {
	Ask(answer interface{}) error
	SetDefault(def string)
}
type Input struct {
	Message    string
	Validators []Validator
	Default    string
	Help       string
}

func (i *Input) Ask(answer interface{}) error {
	prompt := &survey.Input{
		Renderer: survey.Renderer{},
		Message:  i.Message,
		Default:  i.Default,
		Help:     i.Help,
	}
	opts := []survey.AskOpt{}
	for _, validator := range i.Validators {
		opts = append(opts, survey.WithValidator(survey.Validator(validator)))
	}
	return survey.AskOne(prompt, answer, opts...)
}

func (i *Input) SetDefault(def string) {
	i.Default = def
}

//
//type SelectionIndex struct {
//	Message    string
//	Options    []string
//	Validators []Validator
//	Default    string
//	Help       string
//}
//
//func (si *SelectionIndex) Ask(answer interface{}) error {
//	prompt := &survey.Select{
//		Renderer: survey.Renderer{},
//		Message:  si.Message,
//		Options:  si.Options,
//		Default:  si.Default,
//		Help:     si.Help,
//	}
//	opts := []survey.AskOpt{}
//	for _, validator := range si.Validators {
//		opts = append(opts, survey.WithValidator(survey.Validator(validator)))
//	}
//	selectAnswer := ""
//	err := survey.AskOne(prompt, &selectAnswer, opts...)
//	if err != nil {
//		return err
//	}
//	for key, value := range si.Options {
//		if value == selectAnswer {
//			answer = fmt.Sprintf("%d", key+1)
//			return nil
//		}
//	}
//	return nil
//}
//func (si *SelectionIndex) SetDefault(def string) {
//	si.Default = def
//}

type Selection struct {
	Message    string
	Options    []string
	Validators []Validator
	Default    string
	Help       string
}

func (s *Selection) Ask(answer interface{}) error {
	prompt := &survey.Select{
		Renderer: survey.Renderer{},
		Message:  s.Message,
		Options:  s.Options,
		Default:  s.Default,
		Help:     s.Help,
	}
	opts := []survey.AskOpt{}
	for _, validator := range s.Validators {
		opts = append(opts, survey.WithValidator(survey.Validator(validator)))
	}
	return survey.AskOne(prompt, answer, opts...)
}

func (s *Selection) SetDefault(def string) {
	s.Default = def
}
