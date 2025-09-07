package ticker

import (
	"errors"
	"log/slog"
	"time"

	"github.com/spf13/viper"
)

var (
	ErrorUnknownTickerType = errors.New("ticker type is not known")
)

// The default ticker interface.
// All tickers must implement the below methods, but may use vastly different implementations.
// This allows for, say, different simulation methods.
//
// To create a new Ticker implementation, make a new file (aptly named), define a struct,
// include a BaseTicker as an embedded component, then define the Initialize and Update methods.
// Then, in the NewTickerFromConfig method below, add the relevant case to the switch statement.
// Don't forget to update documentation in the README.
type Ticker interface {
	// Return the name of the ticker.
	// This method is implemented by the BaseTicker struct.
	String() string

	// Get the current value of a ticker.
	//
	// Note that no ticker value may be below zero, as a rule of business logic.
	// This method is implemented by the BaseTicker struct.
	GetValue() float64

	// Initialize the ticker using the passed viper config map.
	// Ensures that all tickers have their initialization code in one file,
	// with the rest of their implementation.
	Initialize(*viper.Viper) error

	// Update the value of a ticker.
	//
	// All updates should finish by ensuring the ticker value is non-negative,
	// that is, clamp the updated value to be zero or larger.
	Update()
}

// Factory pattern to initialize the new ticker.
// Returns:
//   - (Ticker, nil) if all checks pass and the ticker is initialized
//   - (nil, ErrorUnknownTickerType) if the ticker type in tickerConfig is unknown
//   - (nil, error) if the ticker is initialized with invalid parameters (e.g. a negative initial value)
//     The specifics of the error are determined by the initialization method of the selected ticker
func NewTickerFromConfig(name string, tickerConfig *viper.Viper) (Ticker, error) {
	tickerType := tickerConfig.GetString("Type")

	var t Ticker
	switch tickerType {
	case "UniformRandom":
		t = &UniformRandomTicker{}
	default:
		return nil, ErrorUnknownTickerType
	}

	if err := t.Initialize(tickerConfig); err != nil {
		return nil, err
	}
	return t, nil
}

// Parse tickers from the viper config, looking into the `Tickers` array to find definitions.
//
// Returns a map from tickerName to the initialized ticker. Note that tickers have not yet been started!
// After parsing tickers, `StartTicker` must be called on each ticker (in a goroutine)
func ParseTickers() map[string]Ticker {
	allTickers := make(map[string]Ticker, 0)

	// viper.GetStringMap("Tickers") gives a map of all tickers specified under
	// the `tickers` key in the config file. Loop over these and call viper.Sub
	// to get the config of each ticker.
	for tickerName := range viper.GetStringMap("Tickers") {
		tickerConfig := viper.Sub("Tickers." + tickerName)
		tickerConfig.Set("name", tickerName) // Add the ticker name to the config as a way to easily pass this along to initializations

		t, err := NewTickerFromConfig(tickerName, tickerConfig)
		if err != nil {
			slog.Error("error when parsing ticker",
				"err", err,
				slog.Group(
					"ticker",
					"name", tickerName,
					"config", tickerConfig.AllSettings(),
				),
			)
			continue
		}

		allTickers[tickerName] = t
		slog.Debug("parsed new ticker",
			slog.Group(
				"ticker",
				"name", tickerName,
				"config", tickerConfig.AllSettings(),
			),
		)
	}

	return allTickers
}
