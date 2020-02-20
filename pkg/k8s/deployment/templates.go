package deployment

import (
	"text/template"
)

type Template struct {
	Structure string
	Data      interface{}
	output    []byte
}

func NewTemplate(str string, data interface{}) *Template {
	return &Template{
		Structure: str,
		Data:      data,
	}
}

func (t *Template) Write(p []byte) (n int, err error) {
	t.output = append(t.output, p...)
	return len(t.output), nil
}
func (t *Template) Get() ([]byte, error) {
	tmpl, err := template.New("tmpl").Parse(t.Structure)
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(t, t.Data)
	if err != nil {
		return nil, err
	}
	return t.output, nil
}
