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

	updatePeriod := time.Duration(viper.GetInt64("UpdatePeriod")) * time.Nanosecond
	var tickerWaitGroup sync.WaitGroup
	for _, ticker := range tickers {
		tickerWaitGroup.Go(func() {
			// An unfortunate name, time.Ticker is a timing device to count a certain time before updating.
			// We will refer to this as the timer throughout to avoid confusion with a stock ticker.
			timer := time.NewTicker(updatePeriod)

			for updateTimerTimestamp := range timer.C {
				updateStartTime := time.Now()
				ticker.Update()
				updateDuration := time.Since(updateStartTime)

				slog.Debug("ticker updated",
					slog.Group("ticker",
						"name", ticker.String(),
						"value", ticker.GetValue(),
					),
					slog.Group("timing",
						"timestamp", updateTimerTimestamp,
						"updateDuration", updateDuration,
						"expectedUpdatePeriod", updatePeriod,
					),
				)

				if updateDuration > updatePeriod {
					slog.Warn("timer update is lagging behind update period",
						slog.Group("ticker",
							"name", ticker.String(),
							"value", ticker.GetValue(),
						),
						slog.Group("timing",
							"timestamp", updateTimerTimestamp,
							"updateDuration", updateDuration,
							"expectedUpdatePeriod", updatePeriod,
						),
					)
				}
			}
		})
	}
	tickerWaitGroup.Wait()
}
