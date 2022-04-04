package main

import (
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Instrument struct {
	gorm.Model
	sdk.Instrument
}

type Candle struct {
	gorm.Model
	sdk.Instrument
}

type DbCli struct {
	*gorm.DB
}

func makeDbCli() (*DbCli, func()) {
	db, err := gorm.Open(sqlite.Open(*dbPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &DbCli{db}, func() {
		dbImpl, _ := db.DB()
		dbImpl.Close()
	}
}

func initDb(dbCli *DbCli) {
	dbCli.AutoMigrate(&Instrument{})
	dbCli.AutoMigrate(&Candle{})
}
