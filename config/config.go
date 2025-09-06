package config

import (
	"io"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

func LoadConfig(configFilePath string) {
	slog.Debug("loading config")

	viper.SetDefault("LogLevel", "info")
	viper.SetDefault("LogFile", "")

	// Read explicitly from the config file give.
	// This may ignore other config paths (e.g. environment variables), worth testing.
	//
	// TODO: Test if other configurations are ignored when using viper.SetConfigFile,
	// or if this only ignores other *files*.
	// See: https://github.com/spf13/viper?tab=readme-ov-file#reading-config-files
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Info("no config file found", "configFilePath", configFilePath)
		} else {
			slog.Error("error during config read", "err", err)
			panic(err)
		}
	}

	// Now config is loaded, all keys should be available through viper.Get
}

