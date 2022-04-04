package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	_ "github.com/mattn/go-sqlite3"
)

var token = flag.String("token", "", "your token")
var isProd = flag.Bool("is_prod", false, "is sandbox env")
var dbPath = flag.String("db_path", "", "db path")

type SdkCli interface{}

func makeInvestCli() SdkCli {

	if *isProd {
		log.Println("cli type: production")
		return sdk.NewRestClient(*token)
	}
	log.Println("cli type: sandbox")
	return sdk.NewSandboxRestClient(*token)

}

func validate() {
	if *token == "" {
		log.Println("token must be specified")
	}
	if *dbPath == "" {
		log.Println("db_path must be specified")
	}
}

func getExistsStocks(dbCli *DbCli) map[string]Instrument {
	var dbResult []Instrument
	dbCli.Find(&dbResult)

	dbStocks := make(map[string]Instrument)

	for _, item := range dbResult {
		dbStocks[item.Ticker] = item
	}
	return dbStocks
}

func getActualStocks(client *sdk.RestClient) []sdk.Instrument {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	actualStocks, err := client.Stocks(ctx)
	if err != nil {
		log.Fatalln(errorHandle(err))
	}
	return actualStocks
}

func actualizeStocks(client *sdk.RestClient, dbCli *DbCli) map[string]Instrument {
	actualStocks := getActualStocks(client)
	existsStocks := getExistsStocks(dbCli)

	oldCounter := 0
	newCounter := 0
	for _, stock := range actualStocks {
		_, found := existsStocks[stock.Ticker]
		if !found {
			log.Println("new instrument with ticker", stock.Ticker)
			newInstrument := Instrument{Instrument: stock}
			existsStocks[stock.Ticker] = newInstrument
			dbCli.Create(newInstrument)
			newCounter++
		} else {
			oldCounter++
		}
	}
	log.Printf("instruments actualized: old=%d, new=%d", oldCounter, newCounter)
	return existsStocks
}

func saveStockCandles(client *sdk.RestClient, dbCli *DbCli, instrument Instrument) {
	log.Println("handle ticker:", instrument.Ticker)
	if instrument.Ticker != "AAPL" {
		return
	}
}

func getData(investCliIface *SdkCli, dbCli *DbCli) {
	var sdkCli *sdk.RestClient
	if *isProd {
		sdkCli = (*investCliIface).(*sdk.RestClient)
	} else {
		sdkCli = (*investCliIface).(*sdk.SandboxRestClient).RestClient
	}

	stocks := actualizeStocks(sdkCli, dbCli)

	//to := time.Now().Add(-time.Hour * 24)
	//from := to.Add(-time.Hour * 12)
	for _, item := range stocks {
		saveStockCandles(sdkCli, dbCli, item)
		//log.Println("[", idx, "]", st)
		//ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		//defer cancel()
		//candles, err := client.Candles(ctx, from, to, sdk.CandleInterval5Min, st.FIGI)
		//if err != nil {
		//	log.Fatalln(errorHandle(err))
		//}
		//for idx, cnd := range candles {
		//	log.Println("[", idx, "]", cnd)
		//}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // инициируем Seed рандома для функции requestID
	flag.Parse()
	validate()

	dbCli, closeDb := makeDbCli()
	defer closeDb()

	initDb(dbCli)

	// Создание
	//db.Create(&Instrument{})

	investCli := makeInvestCli()
	getData(&investCli, dbCli)

	//stream()
}

//func stream() {
//	logger := log.New(os.Stdout, "[invest-openapi-go-sdk]", log.LstdFlags)
//
//	client, err := sdk.NewStreamingClient(logger, *token)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	defer client.Close()
//
//	// Запускаем цикл обработки входящих событий. Запускаем асинхронно
//	// Сюда будут приходить сообщения по подпискам после вызова соответствующих методов
//	// SubscribeInstrumentInfo, SubscribeCandle, SubscribeOrderbook
//	go func() {
//		err = client.RunReadLoop(func(event interface{}) error {
//			logger.Printf("Got event %+v", event)
//			return nil
//		})
//		if err != nil {
//			log.Fatalln(err)
//		}
//	}()
//
//	log.Println("Подписка на получение событий по инструменту BBG005DXJS36 (TCS)")
//	err = client.SubscribeInstrumentInfo("BBG005DXJS36", requestID())
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	log.Println("Подписка на получение свечей по инструменту BBG005DXJS36 (TCS)")
//	err = client.SubscribeCandle("BBG005DXJS36", sdk.CandleInterval5Min, requestID())
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	log.Println("Подписка на получения стакана по инструменту BBG005DXJS36 (TCS)")
//	err = client.SubscribeOrderbook("BBG005DXJS36", 10, requestID())
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	// Приложение завершится через 10секунд.
//	// Hint: В боевом приложении лучше обрабатывать сигналы завершения и работать в бесконечном цикле
//	time.Sleep(10 * time.Second)
//
//	log.Println("Отписка от получения событий по инструменту BBG005DXJS36 (TCS)")
//	err = client.UnsubscribeInstrumentInfo("BBG005DXJS36", requestID())
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	log.Println("Отписка от получения свечей по инструменту BBG005DXJS36 (TCS)")
//	err = client.UnsubscribeCandle("BBG005DXJS36", sdk.CandleInterval5Min, requestID())
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	log.Println("Отписка от получения стакана по инструменту BBG005DXJS36 (TCS)")
//	err = client.UnsubscribeOrderbook("BBG005DXJS36", 10, requestID())
//	if err != nil {
//		log.Fatalln(err)
//	}
//}

//func sandboxRest() {
//	// Особенность работы в песочнице состоит в том что перед началом работы надо вызвать метод Register
//	// и выставить установить(нарисовать) себе активов в портфеле методами SetCurrencyBalance (валютные активы) и
//	// SetPositionsBalance (НЕ валютные активы)
//	// Обнулить портфель можно методом Clear
//	// Все остальные методы rest клиента так же доступны в песочнице
//	client := sdk.NewSandboxRestClient(*token)
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	stocks, err := client.Stocks(ctx)
//	if err != nil {
//		log.Fatalln(errorHandle(err))
//	}
//	for idx, st := range stocks {
//		log.Println("[", idx, "]", st)
//	}
//
//	to := time.Now()
//	from := to.Add(-time.Hour * 2)
//	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	candles, err := client.Candles(ctx, from, to, sdk.CandleInterval5Min, "BBG000HLJ7M4")
//	if err != nil {
//		log.Fatalln(errorHandle(err))
//	}
//	for idx, cnd := range candles {
//		log.Println("[", idx, "]", cnd)
//	}
//
//	return
//	//
//	//log.Println("Регистрация обычного счета в песочнице")
//	//account, err := client.Register(ctx, sdk.AccountTinkoff)
//	//if err != nil {
//	//	log.Fatalln(errorHandle(err))
//	//}
//	//log.Printf("%+v\n", account)
//	//
//	//defer func() {
//	//	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	//	defer cancel()
//	//	log.Println("Удаляем счет в песочнице")
//	//	err = client.Remove(ctx, account.ID)
//	//	if err != nil {
//	//		log.Fatalln(err)
//	//	}
//	//}()
//	//
//	//ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	//defer cancel()
//	//
//	//log.Println("Рисуем себе 100500 рублей в портфеле песочницы")
//	//err = client.SetCurrencyBalance(ctx, account.ID, sdk.RUB, 100500)
//	//if err != nil {
//	//	log.Fatalln(err)
//	//}
//	//
//	//ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	//defer cancel()
//	//
//	//log.Println("Получение списка валютных и НЕ валютных активов портфеля для счета по-умолчанию")
//	//// Метод является совмещеним PositionsPortfolio и CurrenciesPortfolio
//	//portfolio, err := client.Portfolio(ctx, sdk.DefaultAccount)
//	//if err != nil {
//	//	log.Fatalln(err)
//	//}
//	//log.Printf("%+v\n", portfolio)
//	//
//	//ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	//defer cancel()
//	//
//	//log.Println("Получение всех брокерских счетов")
//	//accounts, err := client.Accounts(ctx)
//	//if err != nil {
//	//	log.Fatalln(errorHandle(err))
//	//}
//	//log.Printf("%+v\n", accounts)
//	//
//	//ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	//defer cancel()
//	//
//	//log.Println("Рисуем себе 100 акций TCS в портфеле песочницы")
//	//err = client.SetPositionsBalance(ctx, account.ID, "BBG005DXJS36", 100)
//	//if err != nil {
//	//	log.Fatalln(err)
//	//}
//	//
//	//ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	//defer cancel()
//	//
//	//log.Println("Очищаем состояние портфеля в песочнице")
//	//err = client.Clear(ctx, account.ID)
//	//if err != nil {
//	//	log.Fatalln(err)
//	//}
//}

//func rest() {
//	client := sdk.NewRestClient(*token)
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Println("Получение всех брокерских счетов")
//	accounts, err := client.Accounts(ctx)
//	if err != nil {
//		log.Fatalln(errorHandle(err))
//	}
//	log.Printf("%+v\n", accounts)
//
//	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Println("Получение валютных инструментов")
//	// Например: USD000UTSTOM - USD, EUR_RUB__TOM - EUR
//	currencies, err := client.Currencies(ctx)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Printf("%+v\n", currencies)
//
//	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Println("Получение фондовых инструментов")
//	// Например: FXMM - Казначейские облигации США, FXGD - золото
//	etfs, err := client.ETFs(ctx)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Printf("%+v\n", etfs)
//
//	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Println("Получение облигационных инструментов")
//	// Например: SU24019RMFS0 - ОФЗ 24019
//	bonds, err := client.Bonds(ctx)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Printf("%+v\n", bonds)
//
//	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Println("Получение акционных инструментов")
//	// Например: SBUX - Starbucks Corporation
//	stocks, err := client.Stocks(ctx)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Printf("%+v\n", stocks)
//
//	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Println("Получение инструменов по тикеру TCS")
//	// Получение инструмента по тикеру, возвращает массив инструментов потому что тикер уникален только в рамках одной биржи
//	// но может совпадать на разных биржах у разных кампаний
//	// Например: https://www.moex.com/ru/issue.aspx?code=FIVE и https://www.nasdaq.com/market-activity/stocks/FIVE
//	// В этом примере получить нужную компанию можно проверив поле Currency
//	instruments, err := client.InstrumentByTicker(ctx, "TCS")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Printf("%+v\n", instruments)
//
//	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Println("Получение инструмента по FIGI BBG005DXJS36 (TCS)")
//	// Получение инструмента по FIGI(https://en.wikipedia.org/wiki/Financial_Instrument_Global_Identifier)
//	// Узнать FIGI нужного инструмента можно методами указанными выше
//	// Например: BBG000B9XRY4 - Apple, BBG005DXJS36 - Tinkoff
//	instrument, err := client.InstrumentByFIGI(ctx, "BBG005DXJS36")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Printf("%+v\n", instrument)
//
//	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Println("Получение списка операций для счета по-умолчанию за последнюю неделю по инструменту(FIGI) BBG000BJSBJ0")
//	// Получение списка операций за период по конкретному инструменту(FIGI)
//	// Например: ниже запрашиваются операции за последнюю неделю по инструменту NEE
//	operations, err := client.Operations(ctx, sdk.DefaultAccount, time.Now().AddDate(0, 0, -7), time.Now(), "BBG000BJSBJ0")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Printf("%+v\n", operations)
//
//	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Println("Получение списка НЕ валютных активов портфеля для счета по-умолчанию")
//	positions, err := client.PositionsPortfolio(ctx, sdk.DefaultAccount)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Printf("%+v\n", positions)
//
//	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Println("Получение списка валютных активов портфеля для счета по-умолчанию")
//	positionCurrencies, err := client.CurrenciesPortfolio(ctx, sdk.DefaultAccount)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Printf("%+v\n", positionCurrencies)
//
//	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Println("Получение списка валютных и НЕ валютных активов портфеля для счета по-умолчанию")
//	// Метод является совмещеним PositionsPortfolio и CurrenciesPortfolio
//	portfolio, err := client.Portfolio(ctx, sdk.DefaultAccount)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Printf("%+v\n", portfolio)
//
//	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Println("Получение списка выставленных заявок(ордеров) для счета по-умолчанию")
//	orders, err := client.Orders(ctx, sdk.DefaultAccount)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Printf("%+v\n", orders)
//
//	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Println("Получение часовых свечей за последние 24 часа по инструменту BBG005DXJS36 (TCS)")
//	// Получение свечей(ордеров)
//	// Внимание! Действуют ограничения на промежуток и доступный размер свечей за него
//	// Интервал свечи и допустимый промежуток запроса:
//	// - 1min [1 minute, 1 day]
//	// - 2min [2 minutes, 1 day]
//	// - 3min [3 minutes, 1 day]
//	// - 5min [5 minutes, 1 day]
//	// - 10min [10 minutes, 1 day]
//	// - 15min [15 minutes, 1 day]
//	// - 30min [30 minutes, 1 day]
//	// - hour [1 hour, 7 days]
//	// - day [1 day, 1 year]
//	// - week [7 days, 2 years]
//	// - month [1 month, 10 years]
//	// Например получение часовых свечей за последние 24 часа по инструменту BBG005DXJS36 (TCS)
//	candles, err := client.Candles(ctx, time.Now().AddDate(0, 0, -1), time.Now(), sdk.CandleInterval1Hour, "BBG005DXJS36")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Printf("%+v\n", candles)
//
//	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Println("Получение ордербука(он же стакан) глубиной 10 по инструменту BBG005DXJS36")
//	// Получение ордербука(он же стакан) по инструменту
//	orderbook, err := client.Orderbook(ctx, 10, "BBG005DXJS36")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Printf("%+v\n", orderbook)
//
//	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Println("Выставление рыночной заявки для счета по-умолчанию на покупку ОДНОЙ акции BBG005DXJS36 (TCS)")
//	// Выставление рыночной заявки для счета по-умолчанию
//	// В примере ниже выставляется заявка на покупку ОДНОЙ акции BBG005DXJS36 (TCS)
//	placedOrder, err := client.MarketOrder(ctx, sdk.DefaultAccount, "BBG005DXJS36", 1, sdk.BUY)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Printf("%+v\n", placedOrder)
//
//	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Println("Выставление лимитной заявки для счета по-умолчанию на покупку ОДНОЙ акции BBG005DXJS36 (TCS) по цене не выше 20$")
//	// Выставление лимитной заявки для счета по-умолчанию
//	// В примере ниже выставляется заявка на покупку ОДНОЙ акции BBG005DXJS36 (TCS) по цене не выше 20$
//	placedOrder, err = client.LimitOrder(ctx, sdk.DefaultAccount, "BBG005DXJS36", 1, sdk.BUY, 20)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Printf("%+v\n", placedOrder)
//
//	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Printf("Отмена ранее выставленной заявки для счета по-умолчанию. %+v\n", placedOrder)
//	// Отмена ранее выставленной заявки для счета по-умолчанию.
//	// ID заявки возвращается в структуре PlacedLimitOrder в поле ID в запросе выставления заявки client.LimitOrder
//	// или в структуре Order в поле ID в запросе получения заявок client.Orders
//	err = client.OrderCancel(ctx, sdk.DefaultAccount, placedOrder.ID)
//	if err != nil {
//		log.Fatalln(err)
//	}
//}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Генерируем уникальный ID для запроса
func requestID() string {
	b := make([]rune, 12)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func errorHandle(err error) error {
	if err == nil {
		return nil
	}

	if tradingErr, ok := err.(sdk.TradingError); ok {
		if tradingErr.InvalidTokenSpace() {
			tradingErr.Hint = "Do you use sandbox token in production environment or vise verse?"
			return tradingErr
		}
	}

	return err
}
