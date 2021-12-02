package cyclopes

import (
	"errors"
	"fmt"
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

/** Check if path exists, else create it */
func CheckPath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
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

	fullpathFilename := fmt.Sprintf("%s/%s.jpeg", strings.TrimSuffix(path, "/"), filename)

	if err := os.WriteFile(fullpathFilename, image, 0644); err != nil {
		return err
	}

	return nil
}
