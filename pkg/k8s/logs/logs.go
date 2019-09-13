package logs

import (
	"context"
	"encoding/json"
	"github.com/fatih/color"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
	"regexp"
	"text/template"
	"time"
)

type Options struct {
	PodQuery       string
	ContainerQuery string
	Timestamps     bool
	Since          time.Duration
	Namespace      string
	Exclude        []string
	Include        []string
	AllNamespaces  bool
	Selector       string
	Tail           int64
	Color          string
}

var defaultOpts = &Options{
	PodQuery:       ".*",
	ContainerQuery: ".*",
	Timestamps:     false,
	Since:          0,
	Namespace:      "",
	Exclude:        nil,
	Include:        nil,
	AllNamespaces:  true,
	Selector:       "",
	Tail:           0,
	Color:          "auto",
}

// Run starts the main run loop
func Run(ctx context.Context, client *client.Client, o *Options) error {
	if o == nil {
		o = defaultOpts
	}
	var err error
	var namespace string
	// A specific namespace is ignored if all-namespaces is provided
	if o.AllNamespaces {
		namespace = ""
	} else {
		namespace = o.Namespace
		if namespace == "" {
			namespace, _, err = client.ClientConfig.Namespace()
			if err != nil {
				return errors.Wrap(err, "unable to get default namespace")
			}
		}
	}
	config, err := parseConfig(o)
	if err != nil {
		return errors.Wrap(err, "failed to parse config")
	}
	added, removed, err := Watch(ctx, client.ClientSet.CoreV1().Pods(namespace), config.PodQuery, config.ContainerQuery, config.LabelSelector)
	if err != nil {
		return errors.Wrap(err, "failed to set up attach")
	}

	tails := make(map[string]*Tail)

	go func() {
		for p := range added {
			id := p.GetID()
			if tails[id] != nil {
				continue
			}

			tail := NewTail(p.Namespace, p.Pod, p.Container, config.Template, &TailOptions{
				Timestamps:   config.Timestamps,
				SinceSeconds: int64(config.Since.Seconds()),
				Exclude:      config.Exclude,
				Include:      config.Include,
				Namespace:    config.AllNamespaces,
				TailLines:    config.TailLines,
			})
			tails[id] = tail

			tail.Start(ctx, client.ClientSet.CoreV1().Pods(p.Namespace))
		}
	}()

	go func() {
		for p := range removed {
			id := p.GetID()
			if tails[id] == nil {
				continue
			}
			tails[id].Close()
			delete(tails, id)
		}
	}()

	<-ctx.Done()

	return nil
}

func parseConfig(o *Options) (*Config, error) {

	var podQuery string
	if o.PodQuery == "" {
		podQuery = ".*"
	} else {
		podQuery = o.PodQuery
	}
	pod, err := regexp.Compile(podQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile regular expression from query")
	}

	var containerQuery string
	if o.ContainerQuery == "" {
		containerQuery = ".*"
	} else {
		containerQuery = o.ContainerQuery
	}

	container, err := regexp.Compile(containerQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile regular expression for container query")
	}

	var exclude []*regexp.Regexp
	for _, ex := range o.Exclude {
		rex, err := regexp.Compile(ex)
		if err != nil {
			return nil, errors.Wrap(err, "failed to compile regular expression for exclusion filter")
		}

		exclude = append(exclude, rex)
	}

	var include []*regexp.Regexp
	for _, inc := range o.Include {
		rin, err := regexp.Compile(inc)
		if err != nil {
			return nil, errors.Wrap(err, "failed to compile regular expression for inclusion filter")
		}

		include = append(include, rin)
	}

	var labelSelector labels.Selector
	selector := o.Selector
	if selector == "" {
		labelSelector = labels.Everything()
	} else {
		labelSelector, err = labels.Parse(selector)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse selector as label selector")
		}
	}

	var tailLines *int64
	if o.Tail != 0 {
		tailLines = &o.Tail
	}
	var t string
	colorFlag := o.Color
	if colorFlag == "auto" {
		t = "[{{color .NamespaceColor .Namespace}}] -> [{{color .PodColor .PodName}}] -> [{{color .ContainerColor .ContainerName}}] -> {{.Message}}"
	} else {
		color.NoColor = true
		t = "[{{.Namespace}}] -> [{{.PodName}}] -> [{{.ContainerName}}] -> {{.Message}}"
	}

	funs := map[string]interface{}{
		"json": func(in interface{}) (string, error) {
			b, err := json.Marshal(in)
			if err != nil {
				return "", err
			}
			return string(b), nil
		},
		"color": func(color color.Color, text string) string {
			return color.SprintFunc()(text)
		},
	}
	template, err := template.New("log").Funcs(funs).Parse(t)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse template")
	}

	if o.Since == 0 {
		o.Since = 172800000000000 // 48h
	}

	return &Config{
		PodQuery:       pod,
		ContainerQuery: container,
		Exclude:        exclude,
		Include:        include,
		Timestamps:     o.Timestamps,
		Since:          o.Since,
		Namespace:      o.Namespace,
		AllNamespaces:  o.AllNamespaces,
		LabelSelector:  labelSelector,
		TailLines:      tailLines,
		Template:       template,
	}, nil
}
