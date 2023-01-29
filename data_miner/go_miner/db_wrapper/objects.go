package db_wrapper

import (
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Instrument struct {
	gorm.Model
	sdk.Instrument
	ISIN string `gorm:"index:ISIN,unique"`
}

type Currency struct {
	gorm.Model
	sdk.Instrument
	FIGI     string `gorm:"index:FIGI,unique"`
	ISIN     string `sql:"-:all"`
	Currency string `sql:"-:all"`
	Type     string `sql:"-:all"`
}

type HistoricalCandle struct {
	gorm.Model
	sdk.Candle
	FIGI     string             `gorm:"index:FIGI,uniqueIndex:FIGI_TS_INTERVAL"`
	Date     datatypes.Date     `gorm:"index:DATE"`
	TS       time.Time          `gorm:"uniqueIndex:FIGI_TS_INTERVAL"`
	Interval sdk.CandleInterval `gorm:"index:INTERVAL,uniqueIndex:FIGI_TS_INTERVAL"`
}
