package builder

import (
	"github.com/ghodss/yaml"
	"strings"
	"time"
)

type Manifest struct {
	Id        int64         `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	Type      string        `json:"type"`
	Items     []interface{} `json:"items"`
	Prepend   []string      `json:"prepend"`
}

func NewManifest() *Manifest {
	return &Manifest{
		Id:        0,
		CreatedAt: time.Now(),
		Type:      "",
		Items:     nil,
		Prepend:   nil,
	}
}

func (m *Manifest) SetID(value int64) *Manifest {
	m.Id = value
	return m
}

func (m *Manifest) SetType(value string) *Manifest {
	m.Type = value
	return m
}

func (m *Manifest) SetItems(values ...interface{}) *Manifest {
	m.Items = append(m.Items, values...)
	return m
}
func (m *Manifest) AddPrepend(value string) *Manifest {
	m.Prepend = append(m.Prepend, value)
	return m
}
func (m *Manifest) String() string {
	var list []string
	list = append(list, m.Prepend...)
	for _, item := range m.Items {
		data, err := yaml.Marshal(item)
		if err != nil {
			continue
		}
		list = append(list, string(data))
	}
	return strings.Join(list, "\n---\n")
}
