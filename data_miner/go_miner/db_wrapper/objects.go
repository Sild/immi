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

type Candle struct {
	gorm.Model
	sdk.Instrument
}
