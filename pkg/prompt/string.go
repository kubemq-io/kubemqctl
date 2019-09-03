package prompt

import survey "github.com/AlecAivazis/survey/v2"

type ValueType string

const (
	ValueTypeString ValueType = "string"
	ValueTypeInt ValueType ="int"
	ValueTypeBool ValueType="bool"
	ValueTypeFloat ValueType="float"

)



type EntryString struct {
	Name string
	Value string
	ValueType ValueType
	Question string
	Description string
	Options []string
	Default string
}

func (e EntryString) Ask () error {
	survey:=survey.
	if e.Options== nil {

	}
}