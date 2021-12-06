package adapters

import "errors"

type Adapter struct {
	Yaml *[]byte
}

type AdapterInterface interface {
	/** Checks if the yaml file have all the necessary fields before starting testing */
	Preflight() error
	/** Execute the adapter */
	Execute(imagePath string) error
}

func NewAdapter(name string, yaml *[]byte) (AdapterInterface, error) {
	switch name {
	case "slack":
		return &SlackAdapter{Yaml: yaml}, nil
	default:
		return nil, errors.New("Adapter not found")
	}
}
