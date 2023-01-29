package miner

import (
	"context"
	"miner/db_wrapper"
	"miner/logger"
	"miner/tink_wrapper"
	"time"
)

func SyncInstruments(dbCli *db_wrapper.DbCli, investCli *tink_wrapper.TinkCli, ctx *context.Context, onDone func()) {
	logger.Info("started")
	syncInstrumentsImpl(dbCli, investCli)
	syncCurrenciesImpl(dbCli, investCli)
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-ticker.C:
			syncInstrumentsImpl(dbCli, investCli)
			syncCurrenciesImpl(dbCli, investCli)
		case <-(*ctx).Done():
			logger.Info("done by context interruption")
			onDone()
			return
		}
	}
}

func syncInstrumentsImpl(dbCli *db_wrapper.DbCli, investCli *tink_wrapper.TinkCli) {
	logger.Info("started")
	defer logger.Info("done")

	dbInstruments, err := dbCli.GetDbInstruments()
	if err != nil {
		logger.Error("Fail to get instruments from db: '%s'", err.Error())
		return
	}
	actualInstruments, err := investCli.GetActualInstruments()
	if err != nil {
		logger.Error("Fail to get actual instruments: '%s'", err.Error())
		return
	}

	logger.Info("db_count: %d, actual_count: %d", len(*dbInstruments), len(actualInstruments))

	created := make([]*db_wrapper.Instrument, 0)
	for _, v := range actualInstruments {
		if _, found := (*dbInstruments)[v.ISIN]; !found {
			if err := dbCli.CreateInstrument(&v); err != nil {
				logger.Warn("Fail to create instrument ticker='%s', name='%s', err='%s'", v.Ticker, v.Name, err.Error())
			} else {
				logger.Debug("Instrument created: %s (%s)", v.Ticker, v.Name)
				created = append(created, &v)
			}
		}
	}
	logger.Info("Created instruments count: %d", len(created))
}

func syncCurrenciesImpl(dbCli *db_wrapper.DbCli, investCli *tink_wrapper.TinkCli) {
	logger.Info("started")
	defer logger.Info("done")

	dbCurrencies, err := dbCli.GetDbCurrencies()
	if err != nil {
		logger.Error("Fail to get currencies from db: '%s'", err.Error())
		return
	}
	actualCurrencies, err := investCli.GetActualCurrencies()
	if err != nil {
		logger.Error("Fail to get actual currencies: '%s'", err.Error())
		return
	}

	logger.Info("db_count: %d, actual_count: %d", len(*dbCurrencies), len(actualCurrencies))

	created := make([]*db_wrapper.Currency, 0)
	for _, v := range actualCurrencies {
		if _, found := (*dbCurrencies)[v.FIGI]; !found {
			if err := dbCli.CreateCurrency(&v); err != nil {
				logger.Warn("Fail to create currency with ticker='%s', name='%s', err='%s'", v.Ticker, v.Name, err.Error())
			} else {
				logger.Debug("Currency created: %s (%s)", v.Ticker, v.Name)
				created = append(created, &v)
			}

		}
	}
	logger.Info("Created currencies count: %d", len(created))
}
