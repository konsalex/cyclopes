package adapters

import (
	"errors"
	"io/ioutil"
	"sort"
	"strings"
)

type Adapter struct {
	Yaml *[]byte
}

type AdapterInterface interface {
	/** Checks if the yaml file have all the necessary fields before starting testing */
	Preflight() error
	/** Execute the adapter */
	Execute(imagePath string) error
}

// NewAdapter is an Adapter factory
func NewAdapter(name string) (AdapterInterface, error) {
	switch name {
	case "slack":
		return &SlackAdapter{}, nil
	case "trello":
		return &TrelloAdapter{}, nil
	default:
		return nil, errors.New("Adapter not found")
	}
}

// SortedImageFileNames returns all the image file names (jpeg), sorted by date added
func SortedImageFileNames(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return filesSortedCleaned, nil
}
