package tink_wrapper

import (
	"context"
	"miner/db_wrapper"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
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
			return nil, err;
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
		return nil, err;
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
