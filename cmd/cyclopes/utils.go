package cyclopes

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// ExtractServerURL is returning:
// 1: the default server url if the user has not defined a remoteURL
// 2: the user's remoteURL if it is defined
func (config *Configuration) ExtractServerURL() string {
	var serverPath string
	if config.Visual.RemoteURL == "" {
		serverPath = DEFAULT_URL
	} else {
		serverPath = config.Visual.RemoteURL
	}

	return serverPath
}

// CheckPath checks if path exists, else creates it
func CheckPath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// ConstructServerURL is cleaning server urls so:
// 1. To allow users even define localhost:3000 and still work
//  2. Remove trailing slash from URL
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

// SaveFile save an image to local storage
func SaveFile(image []byte, path string, filename string) error {

	fullpathFilename := fmt.Sprintf("%s/%s.jpeg", strings.TrimSuffix(path, "/"), filename)

	if err := os.WriteFile(fullpathFilename, image, 0644); err != nil {
		return err
	}

	return nil
}
