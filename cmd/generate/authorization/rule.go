package authorization

import (
	"github.com/AlecAivazis/survey/v2"
)

type Rule struct {
	ClientID    string
	Events      bool
	EventsStore bool
	Queues      bool
	Commands    bool
	Queries     bool
	Channel     string
	Read        bool
	Write       bool
}

func getRule() (*Rule, error) {
	answers := struct {
		ClientID  string
		Channel   string
		Resources []string
		Actions   []string
	}{}
	qs := []*survey.Question{
		{
			Name: "clientID",
			Prompt: &survey.Input{
				Message: "Set Access for ClientId (value/regular): ",
				Default: ".*",
				Help:    "",
				Suggest: nil,
			},
			Validate: survey.Required,
		},
		{
			Name: "channel",
			Prompt: &survey.Input{
				Message: "Set Access for Channel (value/regular):",
				Default: ".*",
				Help:    "",
				Suggest: nil,
			},
			Validate: survey.Required,
		},
		{
			Name: "resources",
			Prompt: &survey.MultiSelect{
				Message: "Set Access to Channel's Resources:",
				Options: []string{"Queues", "Events", "Events Store", "Queries", "Commands"},
				Default: nil,
				Help:    "Set Access to resources",
			},
		},
		{
			Name: "actions",
			Prompt: &survey.MultiSelect{
				Message: "Set Actions on Resources:",
				Options: []string{"Read", "Write"},
				Default: nil,
				Help:    "Set Actions on resources",
			},
		},
	}
	err := survey.Ask(qs, &answers)
	if err != nil {
		return nil, err
	}
	r := &Rule{
		ClientID:    answers.ClientID,
		Events:      false,
		EventsStore: false,
		Queues:      false,
		Commands:    false,
		Queries:     false,
		Channel:     answers.Channel,
		Read:        false,
		Write:       false,
	}
	for _, resource := range answers.Resources {
		switch resource {
		case "Queues":
			r.Queues = true

		case "Events":
			r.Events = true

		case "Events Store":
			r.EventsStore = true

		case "Queries":
			r.Queries = true

		case "Commands":
			r.Commands = true
		}
	}
	for _, action := range answers.Actions {
		switch action {
		case "Read":
			r.Read = true
		case "Write":
			r.Write = true
		}
	}
	return r, nil
}
