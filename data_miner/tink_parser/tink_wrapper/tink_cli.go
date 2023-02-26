package tink_wrapper

import (
	"context"
	"miner/db_wrapper"
	"miner/logger"
	"strings"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"gorm.io/datatypes"
)

type TinkCli struct {
	impl *sdk.RestClient
}

func NewInvestCli(token string) *TinkCli {
	return &TinkCli{
		impl: sdk.NewRestClient(token),
	}
}

func (cli *TinkCli) GetActualInstruments() ([]db_wrapper.Instrument, error) {
	result := make([]db_wrapper.Instrument, 0)
	for _, getInstrumentFunc := range []func(context.Context) ([]sdk.Instrument, error){
		cli.impl.Stocks,
		cli.impl.Bonds,
		cli.impl.ETFs,
	} {
		instruments, err := cli.getInstruments(getInstrumentFunc)
		if err != nil {
			return nil, err
		}
		for _, inst := range instruments {
			result = append(result, db_wrapper.Instrument{Instrument: inst, ISIN: inst.ISIN})
		}
	}
	return result, nil
}

func (cli *TinkCli) GetActualCurrencies() ([]db_wrapper.Currency, error) {
	result := make([]db_wrapper.Currency, 0)
	instruments, err := cli.getInstruments(cli.impl.Currencies)
	if err != nil {
		return nil, err
	}
	for _, inst := range instruments {
		result = append(result, db_wrapper.Currency{Instrument: inst, FIGI: inst.FIGI})
	}
	return result, nil
}

func (cli *TinkCli) getInstruments(getFunc func(context.Context) ([]sdk.Instrument, error)) ([]sdk.Instrument, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return getFunc(ctx)
}

func (cli *TinkCli) GetHistoricalCandles(figi string, from time.Time, to time.Time, interval sdk.CandleInterval) ([]db_wrapper.HistoricalCandle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := make([]db_wrapper.HistoricalCandle, 0)
	candles, err := cli.impl.Candles(ctx, from, to, interval, figi)
	if err != nil {
		if strings.Contains(err.Error(), "code=429") {
			logger.Warn("Rate limit achieved. Sleep 60 seconds...")
			time.Sleep(60 * time.Second)
			return cli.GetHistoricalCandles(figi, from, to, interval)
		}
		return nil, err
	}
	for _, cndl := range candles {
		result = append(
			result,
			db_wrapper.HistoricalCandle{
				Candle:   cndl,
				FIGI:     cndl.FIGI,
				Date:     datatypes.Date(cndl.TS),
				TS:       cndl.TS,
				Interval: cndl.Interval,
			},
		)
	}

	return result, nil

}
