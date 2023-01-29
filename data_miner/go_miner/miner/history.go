package miner

import (
	"context"
	"miner/db_wrapper"
	"miner/helper"
	"miner/logger"
	"miner/tink_wrapper"
	"strings"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
)

func SyncHistoryCandles(dbCli *db_wrapper.DbCli, investCli *tink_wrapper.TinkCli, ctx *context.Context, onDone func()) {
	logger.Info("started")
	syncHistoryCandlesImpl(dbCli, investCli, ctx)
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			syncHistoryCandlesImpl(dbCli, investCli, ctx)
		case <-(*ctx).Done():
			logger.Info("done by context interruption")
			onDone()
			return
		}
	}
}

func syncHistoryCandlesImpl(dbCli *db_wrapper.DbCli, investCli *tink_wrapper.TinkCli, ctx *context.Context) {
	logger.Info("started")
	defer logger.Info("done")

	dbInstruments, err := dbCli.GetDbInstruments()
	if err != nil {
		logger.Error("Fail to get instruments from db: '%s'", err.Error())
	}

	dbCurrencies, err := dbCli.GetDbCurrencies()
	if err != nil {
		logger.Error("Fail to get instruments from db: '%s'", err.Error())
	}

	figies := make([]string, 0)
	for _, inst := range *dbInstruments {
		figies = append(figies, inst.FIGI)
	}
	for _, cur := range *dbCurrencies {
		figies = append(figies, cur.FIGI)
	}

	for pos, figi := range figies {
		if figi != "BBG000B9XRY4" {
			continue
		}
		logger.Info("run sync for figi=%s (%d/%d)", figi, pos+1, len(figies))
		if err := syncHistoryCandlesForFigi(dbCli, investCli, figi, ctx); err != nil {
			logger.Error("Error during sync candle. figi='%s'", figi)
			return
		}
		return
	}
}

func syncHistoryCandlesForFigi(dbCli *db_wrapper.DbCli, investCli *tink_wrapper.TinkCli, figi string, ctx *context.Context) error {
	// sync period is 1 day
	// at first, receive the last existing date in DB
	lastCandleTime, err := dbCli.GetLastMinuteCandleTime(figi)
	if err != nil {
		return err
	}

	lastCandleTime.Truncate(24 * time.Hour)

	currentTime := time.Now().UTC()
	currentTime.Truncate(24 * time.Hour)

	logger.Info("figi='%s', last_date='%s', current_date='%s', days_to_parse='%d'", figi, helper.FormatDate(lastCandleTime), helper.FormatDate(currentTime), int(currentTime.Sub(lastCandleTime).Hours()/24+0.5)+1)
	time.Sleep(5 * time.Second)
	for currentTime.After(lastCandleTime) {
		if (*ctx).Err() != nil {
			logger.Warn("Context was cancelled")
			return err
		}
		// first of all, check candles are exist in current year
		startOfYear := time.Date(lastCandleTime.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		endOfYear := time.Date(lastCandleTime.Year()+1, 1, 1, 0, 0, 0, 0, time.UTC)
		candles, yearCandleErr := investCli.GetHistoricalCandles(figi, startOfYear, endOfYear, sdk.CandleInterval1Day)
		if yearCandleErr != nil {
			return yearCandleErr
		}
		logger.Debug("figi='%s': got %d day-candles for year='%d'", figi, len(candles), startOfYear.Year())
		if len(candles) == 0 {
			lastCandleTime = endOfYear
			continue
		}
		// year is not-empty. Save (or update) received day-candles 
		createdCount := 0
		for _, cndl := range candles {
			if err := dbCli.CreateHistoricalCandle(&cndl); err != nil {
				if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					logger.Error("figi='%s': fail to save candle with ts='%s', date='%s', err='%s'", cndl.FIGI, helper.FormatTime(cndl.TS), helper.FormatTime(time.Time(cndl.Date)), err.Error())
				} else {
					logger.Warn("figi='%s': candle was not safe because of duplicate key: ts='%s', date='%s'", cndl.FIGI, helper.FormatTime(cndl.TS), helper.FormatTime(time.Time(cndl.Date)))

				}
			} else {
				createdCount += 1
			}
		}
		logger.Info("figi='%s', year='%s': handled %d year-candles", figi, helper.FormatDate(startOfYear), createdCount)
		
		// and then iterate over min-candles
	}
	return nil
}
