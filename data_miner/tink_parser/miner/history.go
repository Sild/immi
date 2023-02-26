package miner

import (
	"context"
	"miner/db_wrapper"
	"miner/helper"
	"miner/logger"
	"miner/tink_wrapper"
	"sort"
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

	sort.Slice(figies, func(i, j int) bool {
		return figies[j] < figies[i]
	})

	for pos, figi := range figies {
		if err := (*ctx).Err(); err != nil {
			logger.Warn("Context was cancelled")
			return
		}
		logger.Info("run sync for figi=%s (%d/%d)", figi, pos+1, len(figies))
		if err := syncHistoryCandlesForFigi(dbCli, investCli, figi, ctx); err != nil {
			logger.Error("Error during sync candle. figi='%s', err='%s'", figi, err.Error())
			return
		}
	}
}

func syncHistoryCandlesForFigi(dbCli *db_wrapper.DbCli, investCli *tink_wrapper.TinkCli, figi string, ctx *context.Context) error {
	// sync period is 1 day
	// at first, receive the last existing date in DB
	lastCandleTime, err := dbCli.GetLastMinuteCandleTime(figi)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			lastCandleTime = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
		} else {
			return err
		}
	}

	lastCandleTime.Truncate(24 * time.Hour)

	currentHour := time.Now().UTC().Truncate(time.Hour)

	logger.Info("figi='%s', last_time='%s', current_hour='%s', days_to_parse~='%d'", figi, helper.FormatTime(lastCandleTime), helper.FormatTime(currentHour), int(currentHour.Sub(lastCandleTime).Hours()/24+0.5)+1)

	reqStartTime := lastCandleTime
	for currentHour.After(reqStartTime) {
		if err := (*ctx).Err(); err != nil {
			logger.Warn("Context was cancelled")
			return err
		}
		// first of all, check candles are exist in current year
		startOfYear := time.Date(reqStartTime.Year(), 1, 1, 0, 0, 0, 0, time.Local)
		endOfYear := time.Date(reqStartTime.Year()+1, 1, 1, 0, 0, 0, 0, time.Local)
		candles, err := investCli.GetHistoricalCandles(figi, startOfYear, endOfYear, sdk.CandleInterval1Day)
		if err != nil {
			return err
		}
		logger.Debug("figi='%s': got %d day-candles for year='%d'", figi, len(candles), startOfYear.Year())
		if len(candles) <= 150 {
			reqStartTime = endOfYear
			continue
		}

		handled := 0
		dates := make(map[string]bool)
		for _, cndl := range candles {
			dates[helper.FormatDate(time.Time(cndl.Date))] = true
			err := dbCli.CreateHistoricalCandle(&cndl)
			if err != nil {
				if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					logger.Error("figi='%s': fail to save candle with ts='%s', date='%s', err='%s'", cndl.FIGI, helper.FormatTime(cndl.TS), helper.FormatTime(time.Time(cndl.Date)), err.Error())
				} else {
					logger.Warn("figi='%s': candle was not safe because of duplicate key: ts='%s', date='%s'", cndl.FIGI, helper.FormatTime(cndl.TS), helper.FormatTime(time.Time(cndl.Date)))
				}
			} else {
				handled++
			}
		}
		if handled > 0 {
			logger.Debug("figi='%s': handled %d %s-candles for year='%d'", figi, handled, sdk.CandleInterval1Day, startOfYear.Year())
		}


		// and then iterate over min-candles
		for endOfYear.After(reqStartTime) {
			if (*ctx).Err() != nil {
				logger.Warn("Context was cancelled")
				return err
			}

			reqEndTime := reqStartTime.Add(12 * time.Hour)
			if reqStartTime.Weekday() == time.Saturday || reqStartTime.Weekday() == time.Sunday {
				logger.Debug("ignoring date='%s' as weekend", helper.FormatDate(reqStartTime))
				reqStartTime = reqEndTime
				continue
			}
			if _, found := dates[helper.FormatDate(reqStartTime)]; !found {
				reqStartTime = reqEndTime
				continue
			}

			candles, err := investCli.GetHistoricalCandles(figi, reqStartTime, reqEndTime, sdk.CandleInterval1Min)
			if err != nil {
				return err
			}
			logger.Debug("figi='%s': got %d %s-candles for time='%s'-'%s'", figi, len(candles), sdk.CandleInterval1Min, helper.FormatTime(reqStartTime), helper.FormatTime(reqEndTime))
			handled = 0
			for _, cndl := range candles {
				if err := dbCli.CreateHistoricalCandle(&cndl); err != nil {
					if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
						logger.Error("figi='%s': fail to save candle with ts='%s', date='%s', err='%s'", cndl.FIGI, helper.FormatTime(cndl.TS), helper.FormatTime(time.Time(cndl.Date)), err.Error())
					} else {
						handled++
					}
				} else {
					handled++
				}
			}
			logger.Debug("figi='%s': handled %d %s-candles for time='%s'-'%s'", figi, handled, sdk.CandleInterval1Min, helper.FormatTime(reqStartTime), helper.FormatTime(reqEndTime))
			reqStartTime = reqEndTime
		}

		reqStartTime = endOfYear

	}
	return nil
}
