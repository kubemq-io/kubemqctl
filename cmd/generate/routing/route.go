package routing

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"strings"
)

type Route struct {
	Key    string
	Routes string
}

func getRoute() (*Route, error) {
	r := &Route{}
	qs := []*survey.Question{
		{
			Name: "key",
			Prompt: &survey.Input{
				Message: "Set Route Key: (key1 ...)",
				Default: "",
				Help:    "",
				Suggest: nil,
			},
			Validate: survey.Required,
		},
		{
			Name: "routes",
			Prompt: &survey.Input{
				Message: "Set Routes: (queues:foo.bar;events:baz.foo...)",
				Default: "",
				Help:    "",
				Suggest: nil,
			},
			Validate: func(ans interface{}) error {
				val := ans.(string)
				routes := strings.Split(val, ";")
				if len(routes) == 0 {
					return fmt.Errorf("routes must have at least one route")
				}
				for _, route := range routes {
					kv := strings.Split(route, ":")
					if len(kv) != 2 {
						return fmt.Errorf("single route must have a channel value")
					}
					switch kv[0] {
					case "queues", "events", "events_store", "routes":
					default:
						return fmt.Errorf("route must start with 'queues' or 'events' or 'events_store")
					}
				}
				return nil
			},
		},
	}
	err := survey.Ask(qs, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
