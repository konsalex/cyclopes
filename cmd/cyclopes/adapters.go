package cyclopes

import "github.com/konsalex/cyclopes/cmd/cyclopes/adapters"

type Adapter struct {
	Yaml *[]byte
}

type AdapterInterface interface {
	/** Checks if the yaml file have all the necessary fields before starting testing */
	Preflight() error
	/** Execute the adapter */
	Execute(imagePath string) error
}

func LoadAdapter(adapter AdapterInterface) error {
	return nil
}

var AdaptersMap = map[string]interface{}{
	"slack":  adapters.SlackAdapter{},
	"trello": adapters.SlackAdapter{},
}
