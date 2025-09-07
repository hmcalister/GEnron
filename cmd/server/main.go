package main

import (
	"flag"
	"log/slog"
	"sync"
	"time"

	"github.com/hmcalister/genron/config"
	"github.com/hmcalister/genron/ticker"
	"github.com/spf13/viper"
)

func main() {
	configFilePath := flag.String("configFilePath", "config.yaml", "Set the file path to the config file. Accepts JSON, YAML, TOML, and envfiles. See README for config specifications.")
	flag.Parse()

	config.LoadConfig(*configFilePath)
	logFilePointer := config.ConfigureLogger()
	if logFilePointer != nil {
		defer logFilePointer.Close()
	}
	slog.Debug("logger configured")

	tickers := ticker.ParseTickers()
	slog.Debug("parsed tickers", "tickers", tickers)

	// --------------------------------------------------------------------------------

	updatePeriod := time.Duration(viper.GetInt64("updateperiod")) * time.Nanosecond
	var tickerWaitGroup sync.WaitGroup
	for _, t := range tickers {
		tickerWaitGroup.Go(func() {
			ticker.StartTicker(t, updatePeriod)
		})
	}
	tickerWaitGroup.Wait()
}
