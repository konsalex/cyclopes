package cyclopes

type Adapter interface {
	/** Checks if the yaml file have all the necessary fields before starting testing */
	Preflight() error
	/** Execute the adapter */
	Execute(imagePath string) error
}

func LoadAdapter(adapter Adapter) error {
	return nil
}
