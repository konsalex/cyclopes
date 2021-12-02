package adapters

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/slack-go/slack"
	"gopkg.in/yaml.v2"
)

type SlackAdapter struct {
	Yaml *[]byte
}

type SlackConfiguration struct {
	OAUTH_TOKEN string `yaml:"OAUTH_TOKEN"`
	CHANNEL_ID  string `yaml:"CHANNEL_ID"`
}

type SlackRoot struct {
	SlackConfiguration `yaml:",inline"`
	Slack              SlackConfiguration `yaml:"slack"`
}

func (s *SlackAdapter) Preflight() error {

	conf := SlackRoot{}
	err := yaml.Unmarshal(*s.Yaml, &conf)

	if err != nil {
		return err
	}

	if conf.Slack.CHANNEL_ID == "" || conf.Slack.OAUTH_TOKEN == "" {
		return errors.New("slack configuration is not set properly")
	}

	api := slack.New(conf.Slack.OAUTH_TOKEN)

	// Validate that user can make authenticated requests
	_, err = api.AuthTest()
	if err != nil {
		return err
	}

	return nil
}

func (s *SlackAdapter) Execute(imagePath string) error {
	conf := SlackRoot{}
	err := yaml.Unmarshal(*s.Yaml, &conf)

	if err != nil {
		return err
	}

	api := slack.New(conf.Slack.OAUTH_TOKEN)

	/*
		Source: https://stackoverflow.com/a/63391026/3247715
		First we need to:
		1. Upload files one-by-one and get their URLs (unpublished)
		2. Compose them in a markdown and post them (Gallery mode not supported)
		3. Create a new button message with the URL of the previous created message to go on top (avoid scrolling up)
	*/
	// var files []string

	files, err := ioutil.ReadDir(imagePath)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	var filesSortedCleaned []string

	for _, f := range files {
		if !f.IsDir() {
			parts := strings.Split(f.Name(), ".")
			if parts[len(parts)-1] == "jpeg" {
				filesSortedCleaned = append(filesSortedCleaned, f.Name())
			}
		}
	}

	if err != nil {
		return err
	}

	images := make(map[string]string)
	// https://nathanleclaire.com/blog/2014/04/27/a-surprising-feature-of-golang-that-colored-me-impressed/
	imagesKeys := make([]string, len(filesSortedCleaned))

	for _, file := range filesSortedCleaned {
		dat, err := os.ReadFile(fmt.Sprintf("%s/%s", imagePath, file))
		if err != nil {
			return err
		}

		params := slack.FileUploadParameters{
			Filetype: "png",
			Title:    file,
			Filename: file,
			Reader:   strings.NewReader(string(dat)),
		}
		imageUploaded, err := api.UploadFile(params)

		if err != nil {
			fmt.Println(err)
		}

		images[file] = imageUploaded.Permalink
		imagesKeys = append(imagesKeys, file)
	}

	var mrkdwnBody = "*Cyclopes Testing* \n\n"
	for _, value := range imagesKeys {
		if value != "" {
			mrkdwnBody += fmt.Sprintf("<%s | %s, >   ", images[value], value)
		}
	}

	channelId, timestamp, err := api.PostMessage(
		conf.Slack.CHANNEL_ID,
		slack.MsgOptionText(mrkdwnBody, false),
	)

	if err != nil {
		log.Fatalf("%s\n", err)
	}

	fmt.Sprintln("Message successfully sent to Channel %s at %s\n", channelId, timestamp)

	return nil
}
