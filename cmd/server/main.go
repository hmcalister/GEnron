package main

import (
	"flag"
	"log/slog"
	"net/http"

	"github.com/hmcalister/genron/cmd/server/config"
	"github.com/hmcalister/genron/cmd/server/servers"
	"github.com/hmcalister/genron/cmd/server/ticker"
	"github.com/hmcalister/genron/gen/api/ticker/v1/tickerv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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

	for n, t := range tickers {
		go func() {
			slog.Debug("starting ticker", "tickerName", n)
			ticker.StartTicker(t)
		}()
	}

	// --------------------------------------------------------------------------------
	mux := http.NewServeMux()

	tickerNames := make([]string, 0, len(tickers))
	for k := range tickers {
		tickerNames = append(tickerNames, k)
	}
	tickerInfoServer := &servers.TickerInfoServer{
		Tickers:     tickers,
		TickerNames: tickerNames,
	}
	tickerInfoServerPath, tickerInfoServerHandler := tickerv1connect.NewTickerInfoServiceHandler(tickerInfoServer)
	mux.Handle(tickerInfoServerPath, tickerInfoServerHandler)

	if err := http.ListenAndServe(
		"localhost:8080",
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		slog.Error("error during listen and serve of http mux", "err", err)
		panic(err)
	}
}
