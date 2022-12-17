package db_wrapper

import (
	"fmt"
	"miner/logger"

	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbCli struct {
	impl *gorm.DB
	mtx  sync.Mutex
}

func NewDbCli(host string, port int, user string, pass string, dbname string) (*DbCli, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, pass, host, port, dbname)

	cli, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	logger.Info("db connection created")
	return &DbCli{
		impl: cli,
	}, nil

}

func (dbCli *DbCli) UpdateSchema() error {
	dbCli.mtx.Lock()
	defer dbCli.mtx.Unlock()

	models := []interface{}{
		(*Instrument)(nil),
		(*Candle)(nil),
	}

	for _, model := range models {
		dbCli.impl.AutoMigrate(model)
	}
	return nil
}

func (dbCli *DbCli) GetDbInstruments() (*map[string]Instrument, error) {
	dbCli.mtx.Lock()
	defer dbCli.mtx.Unlock()

	res_data := make([]Instrument, 0)
	result := dbCli.impl.Find(&res_data)
	if result.Error != nil {
		return nil, result.Error
	}

	instruments := make(map[string]Instrument)
	for _, v := range res_data {
		instruments[v.ISIN] = v
	}
	return &instruments, nil
}

func (dbCli *DbCli) CreateInstrument(obj *Instrument) error {
	dbCli.mtx.Lock()
	defer dbCli.mtx.Unlock()

	return dbCli.impl.Create(&obj).Error
}

func (dbCli *DbCli) getInstrumentByFigi() *Instrument {
	dbCli.mtx.Lock()
	defer dbCli.mtx.Unlock()

	// res := ms([]Instrument, 0)
	return nil
}

func (dbCli *DbCli) getInstrumentByTicker() *Instrument {
	dbCli.mtx.Lock()
	defer dbCli.mtx.Unlock()

	return &Instrument{}
}
