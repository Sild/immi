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
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			syncInstrumentsImpl(dbCli, investCli)
		case <-(*ctx).Done():
			logger.Info("done by context")
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
			dbCli.CreateInstrument(&v)
			logger.Debug("Instrument created: %s (%s)", v.Ticker, v.Name)
			created = append(created, &v)
		}
	}
	logger.Info("Created instruments count: %d", len(created))
}
