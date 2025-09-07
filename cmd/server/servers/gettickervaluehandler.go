package servers

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/hmcalister/genron/cmd/server/ticker"
	tickerv1 "github.com/hmcalister/genron/gen/api/ticker/v1"
)

var (
	ErrorTickerDoesNotExist = connect.NewError(connect.CodeNotFound, errors.New("no ticker exists with requested name"))
)

type GetTickerValueServer struct {
	// A map from ticker name to ticker structs
	Tickers map[string]ticker.Ticker
}

func (serv *GetTickerValueServer) GetTickerValue(
	ctx context.Context,
	req *connect.Request[tickerv1.GetTickerValueRequest],
) (*connect.Response[tickerv1.GetTickerValueResponse], error) {
	slog.Info("new get ticker value request", "ctx", ctx, "reqHeader", req.Header(), "reqMsg", req.Msg)

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	requestedTicker, ok := serv.Tickers[req.Msg.TickerName]
	if !ok {
		slog.Info("requested ticker name does not exist",
			"reqMsg", req.Msg,
			"requestedTickerName", req.Msg.TickerName,
		)
		return nil, ErrorTickerDoesNotExist
	}

	// TODO: Rate limiting, exhausting, etc...
	tickerName, newValue, lastUpdatedTimestamp := requestedTicker.GetInfo()

	res := connect.NewResponse(&tickerv1.GetTickerValueResponse{
		TickerName:           tickerName,
		TickerValue:          newValue,
		LastUpdatedTimestamp: lastUpdatedTimestamp.UnixNano(),
	})
	return res, nil
}
