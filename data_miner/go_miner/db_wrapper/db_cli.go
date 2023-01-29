package db_wrapper

import (
	"fmt"
	"miner/helper"
	"miner/logger"
	"time"

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
		(*Currency)(nil),
		(*HistoricalCandle)(nil),
	}

	for _, model := range models {
		if err := dbCli.impl.AutoMigrate(model); err != nil {
			return err
		}
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

func (dbCli *DbCli) GetDbCurrencies() (*map[string]Currency, error) {
	dbCli.mtx.Lock()
	defer dbCli.mtx.Unlock()

	res_data := make([]Currency, 0)
	result := dbCli.impl.Find(&res_data)
	if result.Error != nil {
		return nil, result.Error
	}

	currencies := make(map[string]Currency)
	for _, v := range res_data {
		currencies[v.FIGI] = v
	}
	return &currencies, nil
}

func (dbCli *DbCli) GetLastMinuteCandleTime(figi string) (time.Time, error) {
	dbCli.mtx.Lock()
	defer dbCli.mtx.Unlock()

	rows, err := dbCli.impl.Model(&HistoricalCandle{}).Select("max(date) as max_date").Where("figi = ?", figi).Rows()
	if err != nil {
		return time.Time{}, err
	}

	defer rows.Close()

	type result struct {
		Date    time.Time
		MaxDate string
	}
	var res result

	for rows.Next() {
		// ScanRows is a method of `gorm.DB`, it can be used to scan a row into a struct
		if err := dbCli.impl.ScanRows(rows, &res); err != nil {
			return time.Time{}, err
		}
		if res.MaxDate == "" {
			return time.Unix(0, 0).UTC(), nil
		}
		logger.Info("%s", res.Date)
		logger.Info("%s", res.MaxDate)
	}

	return helper.DBDateFromStr(res.MaxDate), nil
}

func (dbCli *DbCli) CreateInstrument(obj *Instrument) error {
	dbCli.mtx.Lock()
	defer dbCli.mtx.Unlock()
	return dbCli.impl.Create(&obj).Error
}

func (dbCli *DbCli) CreateCurrency(obj *Currency) error {
	dbCli.mtx.Lock()
	defer dbCli.mtx.Unlock()
	return dbCli.impl.Create(&obj).Error
}

func (dbCli *DbCli) CreateHistoricalCandle(obj *HistoricalCandle) error {
	dbCli.mtx.Lock()
	defer dbCli.mtx.Unlock()
	if dbCli.impl.Model(&obj).Where("figi = ? and ts = ? and interval = ?", obj.FIGI, obj.TS, obj.Interval).Updates(&obj).RowsAffected == 0 {
		return dbCli.impl.Create(&obj).Error
	} else {
		logger.Debug("figi='%s' for ts='%s' updated instead of created", obj.FIGI, obj.TS)
	}
	return nil
}

// func (dbCli *DbCli) getInstrumentByFigi() *Instrument {
// 	dbCli.mtx.Lock()
// 	defer dbCli.mtx.Unlock()

// 	// res := ms([]Instrument, 0)
// 	return nil
// }

// func (dbCli *DbCli) getInstrumentByTicker() *Instrument {
// 	dbCli.mtx.Lock()
// 	defer dbCli.mtx.Unlock()

// 	return &Instrument{}
// }
