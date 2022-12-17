package miner

import (
	"context"
	"miner/db_wrapper"
	"miner/logger"
	"miner/tink_wrapper"
	"time"
)

func SyncHistory(dbCli *db_wrapper.DbCli, investCli *tink_wrapper.TinkCli, ctx *context.Context, onDone func()) {
	logger.Info("started")
	syncHistoryImpl(dbCli, investCli)
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

func syncHistoryImpl(dbCli *db_wrapper.DbCli, investCli *tink_wrapper.TinkCli) {
	logger.Info("started")
	defer logger.Info("done")

	dbInstruments, err := dbCli.GetDbInstruments()
	if err != nil {
		logger.Error("Fail to get instruments from db: '%s'", err.Error())
	}
	actualInstruments, err := investCli.GetActualInstruments()
	if err != nil {
		logger.Error("Fail to get actual instruments: '%s'", err.Error())
	}

	logger.Info("db_count: %d, actual_count: %d", len(*dbInstruments), len(actualInstruments))

	createCounter := 0
	for _, v := range actualInstruments {
		if _, found := (*dbInstruments)[v.ISIN]; !found {
			dbCli.CreateInstrument(&v)
			createCounter++
		}
	}
	logger.Info("Created instruments count: %d", createCounter)
}
