package builder

import (
	"github.com/Masterminds/sprig"
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
	fmap := sprig.TxtFuncMap()
	tmpl := template.Must(template.New("tmpl").Funcs(fmap).Parse(t.Structure))
	err := tmpl.Execute(t, t.Data)
	if err != nil {
		return nil, err
	}
	return t.output, nil
}
