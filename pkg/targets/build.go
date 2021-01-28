package targets

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"io/ioutil"
	"strings"
)

func BuildRequest() ([]byte, error) {
	fmt.Println("KubeMQ Target Request Builder:")
	mdInput := ""
	promptMetadata := &survey.Input{
		Message: "Set Metadata Key/Value (key1=value1,key2=value2....)",
		Default: "",
		Help:    "Set Metadata Key/Value (key1=value1,key2=value2....)",
	}
	err := survey.AskOne(promptMetadata, &mdInput)
	if err != nil {
		return nil, err
	}
	params := strings.Split(mdInput, ",")
	if len(params) == 0 {
		return nil, fmt.Errorf("kubemq target request must have a metadata key/value pairs")
	}
	md := NewMetadata()
	for _, param := range params {
		kv := strings.Split(param, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("kubemq target request key Value must have key=value format")
		}
		md.Set(kv[0], kv[1])
	}
	loadDataOptions := ""
	loadDataSelectionPrompt := &survey.Select{
		Renderer: survey.Renderer{},
		Message:  "Set Request Data loading options:",
		Options:  []string{"Empty Request data", "Enter Request Data", "Load from file"},
		Default:  "Empty Request data",
		Help:     "Set KubeMQ Target Request data loading options",
	}
	err = survey.AskOne(loadDataSelectionPrompt, &loadDataOptions)
	if err != nil {
		return nil, err
	}
	var data []byte
	switch loadDataOptions {
	case "Empty Request data":

	case "Enter Request Data":
		dataInput := ""
		promptDataInput := &survey.Input{
			Message: "Enter Request Data:",
			Default: "",
			Help:    "Enter Request Data",
		}
		err := survey.AskOne(promptDataInput, &dataInput)
		if err != nil {
			return nil, err
		}
		data = []byte(dataInput)
	case "Load from file":
		dataFileName := ""
		promptDataFileName := &survey.Input{
			Message: "Enter Filename to load:",
			Default: "",
			Help:    "Enter Filename to load:",
		}
		err := survey.AskOne(promptDataFileName, &dataFileName)
		if err != nil {
			return nil, err
		}
		data, err = ioutil.ReadFile(dataFileName)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid request data selection input")

	}

	return NewRequest().
		SetMetadata(md).
		SetData(data).
		MarshalBinary(), nil
}
