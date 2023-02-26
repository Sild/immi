package tsignal

import (
	"miner/db_wrapper"
	"miner/helper"
	"miner/logger"
	"miner/tink_wrapper"
	"sort"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
)

var (
// trackFigi = []string{"ATVI", "SPG"}
)

func Run(dbCli *db_wrapper.DbCli, investCli *tink_wrapper.TinkCli) {
	ticker := "ATVI"

	instr, err := dbCli.GetDbInstrumentByTicker(ticker)
	if err != nil {
		logger.Error("Fail to get instrument for ticker='%s': err='%s'", ticker, err.Error())
		return
	}
	candlesStorage := make([]db_wrapper.HistoricalCandle, 0)

	currentHour := time.Now().UTC().Truncate(time.Hour)
	currentHour = time.Date(2022, 11, 9, 0, 0, 0, 0, time.UTC)
	// receive 1min-candles for last week
	startReceiveDate := currentHour.Add(-1 * time.Hour * 24 * 7)
	for currentHour.After(startReceiveDate) {
		endReceiveDate := startReceiveDate.Add(time.Hour * 24)
		candles, err := investCli.GetHistoricalCandles(instr.FIGI, startReceiveDate, endReceiveDate, sdk.CandleInterval1Min)
		if err != nil {
			logger.Error("'%s': Fail to get candles for '%s'-'%s', err: '%s'", instr.FIGI, helper.FormatTime(startReceiveDate), helper.FormatTime(endReceiveDate), err.Error())
			return
		}
		candlesStorage = append(candlesStorage, candles...)
		startReceiveDate = endReceiveDate
	}
	sort.SliceStable(candlesStorage, func(i int, j int) bool {
		return candlesStorage[i].TS.Before(candlesStorage[j].TS)
	})
	logger.Info("figi='%s', candles_cnt='%d'", instr.FIGI, len(candlesStorage))
	for idx, cndl := range candlesStorage {
		err := dbCli.CreateHistoricalCandle(&cndl)
		logger.Info("%d: ts='%s', o='%.2f', c='%.2f', l='%.2f', h='%.2f', v='%.1f', svd=%t",
			idx,
			helper.FormatTime(cndl.TS),
			cndl.OpenPrice,
			cndl.ClosePrice,
			cndl.LowPrice,
			cndl.HighPrice,
			cndl.Volume,
			err == nil,
		)
	}

}
