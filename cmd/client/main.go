package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"connectrpc.com/connect"
	tickerv1 "github.com/hmcalister/genron/gen/api/ticker/v1"
	"github.com/hmcalister/genron/gen/api/ticker/v1/tickerv1connect"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	LOG_LEVEL           slog.Level    = slog.LevelInfo
	SERVER_ADDR         string        = "http://localhost:8080"
	SERVER_POLLING_RATE time.Duration = 1_000_000 * time.Nanosecond
	MEASURE_DURATION    time.Duration = 1 * time.Second
	DATA_FILE           string        = "polledData.json"
)

type TickerData struct {
	// The name of this ticker
	TickerName string

	// The values of this ticker over time
	// One-to-one with tickerTimestampHistory
	TickerValueHistory []float64

	// The timestamps at which tickerValueHistory were updated
	TickerTimestampHistory []int64
}

func main() {
	slogHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     LOG_LEVEL,
	})
	slog.SetDefault(slog.New(slogHandler))

	// Step Zero: Make sure we can actually open the data file!
	dataFile, err := os.Create(DATA_FILE)
	if err != nil {
		slog.Error("cannot open data file",
			"dataFile", dataFile,
			"err", err,
		)
		panic(err)
	}
	defer dataFile.Close()

	// Step One: Get all the ticker names for polling later.
	client := tickerv1connect.NewTickerInfoServiceClient(
		http.DefaultClient,
		SERVER_ADDR,
		// connect.WithGRPC(),
	)
	res, err := client.GetAllTickerNames(
		context.Background(),
		&connect.Request[emptypb.Empty]{},
	)
	if err != nil {
		slog.Error("error when requesting all ticker names", "err", err)
		panic(err)
	}
	slog.Debug("all ticker names received", "responseHeader", res.Header(), "responseMsg", res.Msg)

	if len(res.Msg.TickerName) == 0 {
		slog.Error("no tickers offered by server")
		return
	}

	tickerData := make([]*TickerData, 0)
	for _, tickerName := range res.Msg.TickerName {
		tickerData = append(tickerData, &TickerData{
			TickerName:             tickerName,
			TickerValueHistory:     make([]float64, 0),
			TickerTimestampHistory: make([]int64, 0),
		})
	}

	// Step Two: Start polling each ticker asynchronously
	var serverPollingWaitgroup sync.WaitGroup
	cancelMap := make(map[string]context.CancelFunc)
	for _, tickerData := range tickerData {
		slog.Debug("start polling server for ticker", "tickerName", tickerData.TickerName)
		ctx, cancel := context.WithCancel(context.Background())
		cancelMap[tickerData.TickerName] = cancel
		client := tickerv1connect.NewTickerInfoServiceClient(
			http.DefaultClient,
			SERVER_ADDR,
			// connect.WithGRPC(),
		)
		serverPollingWaitgroup.Go(func() {
			pollTicker(ctx, client, tickerData)
		})
	}

	// Step Three: Wait for some time to collect data
	time.Sleep(MEASURE_DURATION)

	// Step Four: Cancel the polling functions and wait for the polling to finish
	for _, cancel := range cancelMap {
		cancel()
	}
	serverPollingWaitgroup.Wait()

	// Step Five: Write the data to disk
	encoder := json.NewEncoder(dataFile)
	if err := encoder.Encode(tickerData); err != nil {
		slog.Error("error when trying to encoder ticker data", "err", err)
		panic(err)
	}
}

func pollTicker(ctx context.Context, client tickerv1connect.TickerInfoServiceClient, tickerData *TickerData) {
	timer := time.NewTicker(SERVER_POLLING_RATE)

	for {
		select {
		case <-ctx.Done():
			slog.Debug("cancel called on polling function", "tickerName", tickerData.TickerName)
			return
		case <-timer.C:
			pollStartTime := time.Now()
			res, err := client.GetTickerValue(ctx, connect.NewRequest(&tickerv1.GetTickerValueRequest{
				TickerName: tickerData.TickerName,
			}))
			if err != nil {
				slog.Error("error when polling server",
					"tickerName", tickerData.TickerName,
					"error", err)
				continue
			}

			slog.Debug("polled data from server",
				"tickerName", tickerData.TickerName,
				"responseHeader", res.Header(),
				"responseMsg", res.Msg)
			if len(tickerData.TickerTimestampHistory) > 0 && tickerData.TickerTimestampHistory[len(tickerData.TickerTimestampHistory)-1] == res.Msg.LastUpdatedTimestamp {
				continue
			}

			tickerData.TickerValueHistory = append(tickerData.TickerValueHistory, res.Msg.TickerValue)
			tickerData.TickerTimestampHistory = append(tickerData.TickerTimestampHistory, res.Msg.LastUpdatedTimestamp)
			pollTotalDuration := time.Since(pollStartTime)
			if pollTotalDuration > SERVER_POLLING_RATE {
				slog.Warn("server polling function longer than SERVER_POLLING_RATE",
					"pollTotalDuration", pollTotalDuration,
					"SERVER_POLLING_RATE", SERVER_POLLING_RATE)
			}
		}
	}
}
