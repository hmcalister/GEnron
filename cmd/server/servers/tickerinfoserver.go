package servers

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/hmcalister/genron/cmd/server/ticker"
	tickerv1 "github.com/hmcalister/genron/gen/api/ticker/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrorTickerDoesNotExist = connect.NewError(connect.CodeNotFound, errors.New("no ticker exists with requested name"))
)

type TickerInfoServer struct {
	// A map from ticker name to ticker structs
	Tickers     map[string]ticker.Ticker
	TickerNames []string
}

func (serv *TickerInfoServer) GetAllTickerNames(
	ctx context.Context,
	req *connect.Request[emptypb.Empty],
) (*connect.Response[tickerv1.GetAllTickerNamesResponse], error) {
	slog.Info("new get all ticker names request",
		"ctx", ctx,
		"reqHeader", req.Header(),
		"reqMsg", req.Msg,
	)

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	res := connect.NewResponse(&tickerv1.GetAllTickerNamesResponse{
		TickerName: serv.TickerNames,
	})
	return res, nil
}

func (serv *TickerInfoServer) GetTickerValue(
	ctx context.Context,
	req *connect.Request[tickerv1.GetTickerValueRequest],
) (*connect.Response[tickerv1.GetTickerValueResponse], error) {
	slog.Info("new get ticker value request",
		"ctx", ctx,
		"reqHeader", req.Header(),
		"reqMsg", req.Msg,
	)

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
