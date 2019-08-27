package logs

import (
	"regexp"
	"text/template"
	"time"

	"k8s.io/apimachinery/pkg/labels"
)

// Config contains the config for logs
type Config struct {
	Namespace      string
	PodQuery       *regexp.Regexp
	Timestamps     bool
	ContainerQuery *regexp.Regexp
	Exclude        []*regexp.Regexp
	Include        []*regexp.Regexp
	Since          time.Duration
	AllNamespaces  bool
	LabelSelector  labels.Selector
	TailLines      *int64
	Template       *template.Template
}
