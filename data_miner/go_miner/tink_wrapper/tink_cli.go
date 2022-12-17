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
	return cli.getActualStocks()
}

func (cli *TinkCli) getActualStocks() ([]db_wrapper.Instrument, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stocks, err := cli.impl.Stocks(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]db_wrapper.Instrument, 0)
	for _, st := range stocks {
		result = append(result, db_wrapper.Instrument{Instrument: st, ISIN: st.ISIN})
	}
	return result, nil
}
