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

// Configure the slog logger using config values in viper.
// This method should only be called after LoadConfig.
//
// Returns the os.File pointer that slog writes to, so it may be gracefully shut:
// ```
// logFilePointer := config.ConfigureLogger()
//
//	if logFilePointer != nil{
//		defer logFilePointer.Close()
//	}
//
// ```
func ConfigureLogger() *os.File {
	logLevel := viper.GetString("LogLevel")
	slogHandlerOptions := slog.HandlerOptions{
		AddSource: true,
	}

	// --------------------------------------------------------------------------------

	switch logLevel {
	case "none":
		// No logging is required, disable the logger and return
		slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
		return nil
	case "error":
		slogHandlerOptions.Level = slog.LevelError
	case "warn":
		slogHandlerOptions.Level = slog.LevelWarn
	case "info":
		slogHandlerOptions.Level = slog.LevelInfo
	case "debug":
		slogHandlerOptions.Level = slog.LevelDebug
	default:
		slog.Error("error when decoding unexpected log level in ConfigureLogger", "LogLevel", logLevel)
		panic("unexpected log level encountered in config")
	}

	// --------------------------------------------------------------------------------

	logFile := viper.GetString("LogFile")
	var logFilePointer *os.File
	var slogHandler slog.Handler
	if logFile == "" {
		logFilePointer = nil
		slogHandler = slog.NewTextHandler(os.Stdout, &slogHandlerOptions)
	} else {
		logFilePointer, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			slog.Error("error while creating log file", "LogFile", logFile, "err", err)
			panic(err)
		}
		slogHandler = slog.NewJSONHandler(logFilePointer, &slogHandlerOptions)
	}

	// --------------------------------------------------------------------------------

	slog.SetDefault(slog.New(slogHandler))
	return logFilePointer
}
