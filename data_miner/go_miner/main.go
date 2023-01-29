package main

import (
	"context"
	"flag"
	"miner/db_wrapper"
	"miner/logger"
	"miner/miner"
	"miner/tink_wrapper"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var token = os.Getenv("TINKOFF_FULL_TOKEN")
var dbHost = "localhost"
var dbPort = 5432
var dbName = "imt"
var dbUser = "postgres"
var dbPassword = os.Getenv("PG_PASS")

type SdkCli interface{}

func validateSettings() {
	if token == "" {
		logger.Fatal("token must be specified")
	}
	if dbPassword == "" {
		logger.Fatal("db_pass must be specified")
	}
}

func signalHandler(signalChannel chan os.Signal, cancelFunc func()) {
	for sig := range signalChannel {
		logger.Info("Handle signal: %s", sig.String())
		cancelFunc()
	}
}

func main() {
	logger.Info("app started")
	defer logger.Info("app done")

	flag.Parse()
	validateSettings()

	investCli := tink_wrapper.NewInvestCli(token)
	dbCli, err := db_wrapper.NewDbCli(dbHost, dbPort, dbUser, dbPassword, dbName)

	if err != nil {
		logger.Fatal("Fail to create db connection: '%s'", err.Error())
	}

	if err := dbCli.UpdateSchema(); err != nil {
		logger.Fatal("Fail to update schema: '%s'", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())

	osChannel := make(chan os.Signal, 1)
	signal.Notify(osChannel, os.Interrupt, syscall.SIGTERM)
	go signalHandler(osChannel, cancel)

	var wg sync.WaitGroup
	wg.Add(1)
	// go miner.SyncInstruments(dbCli, investCli, &ctx, func() { wg.Done() })
	go miner.SyncHistoryCandles(dbCli, investCli, &ctx, func() { wg.Done() })
	wg.Wait()
}
