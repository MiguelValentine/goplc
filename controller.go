package goplc

type Controller struct {
	config *Config

	Name string
}

func New(config *Config) {
	controller := Controller{}

	if config != nil {
		controller.config = config
	}
}
