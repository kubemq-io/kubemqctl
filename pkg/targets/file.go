package targets

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"io/ioutil"
)

func BuildFile() ([]byte, error) {
	loadFileOptions := ""
	loadFileSelectionPrompt := &survey.Select{
		Renderer: survey.Renderer{},
		Message:  "Select Load data from file:",
		Options:  []string{"Open an Editor", "Open a File"},
		Default:  "Open an Editor",
		Help:     "Select Load data from file:",
	}
	err := survey.AskOne(loadFileSelectionPrompt, &loadFileOptions)
	if err != nil {
		return nil, err
	}
	switch loadFileOptions {
	case "Open an Editor":
		dataFromEditor := ""
		promptDataFromEditor := &survey.Editor{
			Message: "Copy & Paste data to editor:",
			Default: "",
			Help:    "Copy & Paste data to editor:",
		}
		err := survey.AskOne(promptDataFromEditor, &dataFromEditor)
		if err != nil {
			return nil, err
		}
		return []byte(dataFromEditor), nil
	case "Open a File":
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
		data, err := ioutil.ReadFile(dataFileName)
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		return nil, fmt.Errorf("invalid selection")
	}
}
