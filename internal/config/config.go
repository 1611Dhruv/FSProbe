package config

import (
	fmt "fmt"
	"github.com/spf13/viper"
)

const (
	BlockSize   = 4096
	VSFMagic    = "VSF\x00\x00\x00\x00\x00"
	MagicLength = 8
	Alignment   = 32
)

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	return nil
}
