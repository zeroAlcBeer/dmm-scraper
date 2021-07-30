package config

import (
	"github.com/spf13/viper"
)

// ViperLoader ...
type ViperLoader struct {
	*viper.Viper
}

// LoadFile ...
func (v *ViperLoader) LoadFile(filename string) (*Configs, error) {
	c := &Configs{}
	viper.SetConfigName(filename)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	err := viper.Unmarshal(c)
	return c, err
}
