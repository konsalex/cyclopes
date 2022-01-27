package cyclopes

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/manifoldco/promptui"
	"github.com/pterm/pterm"
	"gopkg.in/yaml.v2"
)

type boolChoice string

const (
	yes boolChoice = "Yes"
	no  boolChoice = "No"
)

func getBooleanQuestion(label string) bool {
	prompt := promptui.Select{
		Label:    label,
		Items:    []string{string(yes), string(no)},
		HideHelp: false,
	}

	_, result, err := prompt.Run()

	if err != nil {
		panic(err)
	}

	if result == string(yes) {
		return true
	}
	return false
}

func getStringValue(label string) string {
	prompt := promptui.Prompt{
		Label: label,
		Validate: func(input string) error {
			if len(input) == 0 {
				return errors.New("input cannot be empty")
			}
			return nil

		},
	}
	result, err := prompt.Run()

	if err != nil {
		panic(err)
	}

	return result
}

// GeneratorCLI is a helper
// CLI function to generate a cyclopes.yml file
func GeneratorCLI() {

	imagePath := getStringValue("What is the path the screenshots are saved? (ex.: `./images`)")
	conf := Configuration{ImagesDir: imagePath}

	visualTesting := getBooleanQuestion("Do you want to perform visual testing?")
	if visualTesting {
		conf.Visual = &VisualTesting{}
		urlMessage := "Do you have a URL to grab the screenshots from?\n" +
			"Either a remote URL like `https://www.example.com` or localhost like `localhost:3000`"

		remoteURL := getBooleanQuestion(urlMessage)
		if remoteURL {
			url := getStringValue("Provide the URL")
			conf.Visual.RemoteURL = url

		} else {
			buildDir := getStringValue("What is the path of your local build?")
			conf.Visual.BuildDir = buildDir
		}
	}

	adapterPosting := getBooleanQuestion("Are you using any adapter (ex. Slack) to post the results?")
	if adapterPosting {
		slack := make(map[string]string)
		slack["OATH_TOKEN"] = "example-token"
		slack["CHANNEL_ID"] = "example-channel-id"
		conf.Adapters = make(map[string]interface{})
		conf.Adapters["slack"] = slack
	}

	filename := getStringValue("What is the desired filename to save the configuration?")

	data, err := yaml.Marshal(conf)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filename, data, 0o644)
	if err != nil {
		panic(err)
	}
	pterm.Success.Println("Configuration file generated 🎉")
}
