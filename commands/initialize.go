package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
)

var moonguardConfig = "moonguard.json"

func initAction(c *cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	destination := path.Join(cwd, moonguardConfig)
	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists in current directory, please delete prior to re-initializing", moonguardConfig)
	}

	var qs = []*survey.Question{
		{
			Name:     "serviceName",
			Prompt:   &survey.Input{Message: "What's your service's name?"},
			Validate: survey.Required,
		},
		{
			Name:     "protoDefinitions",
			Prompt:   &survey.Input{Message: "Where are your protobuf definitions? (e.x. **/*.proto)"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		ServiceName      string `json:"serviceName"`
		ProtoDefinitions string `json:"protoDefinitions"`
	}{}

	err = survey.Ask(qs, &answers)
	if err != nil {
		return fmt.Errorf("there was a problem with your responses: %s", err)
	}

	data, _ := json.MarshalIndent(answers, "", "    ")

	err = ioutil.WriteFile(destination, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to %s: %s", destination, err)
	}

	return nil
}

// GetInitializeCommand command for initializing a moonguard config
func GetInitializeCommand() *cli.Command {
	return &cli.Command{
		Name:   "init",
		Action: initAction,
		Usage:  "initiate a new moonguard service definition in your current directory",
	}
}
