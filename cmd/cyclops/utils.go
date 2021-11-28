package cyclops

import (
	"errors"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
)

func (config *Configuration) ExtractServerURL() (string, error) {
	var serverPath string

	if config.Server {
		serverPath = DEFAULT_URL
	} else {
		if config.ServerURL == "" {
			return "", errors.New("Server url is not specified")
		}
		serverPath = config.ServerURL
	}

	return serverPath, nil
}

func CheckPath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0744)
		return err
	}
	return nil
}

/*
  Clean server urls
  1. To allow users even define localhost:3000 and still work
  2. Remove trailing slash from URL
*/
func ConstructServerURL(rawURL string) (string, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)

	if err != nil {
		return "", err
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "http://" + strings.TrimSuffix(rawURL, "/"), nil
	}

	return strings.TrimSuffix(rawURL, "/"), nil
}

func SaveFile(image []byte, path string, filename string) error {

	if err := CheckPath(path); err != nil {
		FatalPrint(err)
	}

	if err := ioutil.WriteFile(strings.TrimSuffix(path, "/")+"/"+filename+".png", image, 0o644); err != nil {
		log.Fatal(err)
	}

	return nil
}
