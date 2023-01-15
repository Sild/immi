package db_wrapper

import (
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
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
	FIGI string `gorm:"index:FIGI,unique"`
	ISIN string `sql:"-:all"`
	Currency string `sql:"-:all"`
	Type string `sql:"-:all"`
}

type Candle struct {
	gorm.Model
	sdk.Instrument
}
