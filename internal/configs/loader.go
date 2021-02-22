package configs

import "github.com/spf13/viper"

// Loader ...
type Loader interface {
	LoadFile(filename string) (*Configs, error)
}

// NewLoader ...
func NewLoader() Loader {
	return &ViperLoader{
		viper.New(),
	}
}
