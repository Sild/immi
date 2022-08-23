package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
)

var token = "t.t6ZrsKMfcqXQs2HchmhMxmKBrO_DS32dHkZPjj6Ul3U1PfrDqoOg_CgY5VGS7a9cxxd__OwiluqHVKlm2LQL1g" //flag.String("token", "t.t.t.wtfwGhQZpQ9Kh37FNSXwSUVqS2L8rxm_lZIljbUo5ff9Trljjd6GLYZkp5sYRw_8D7dcRQqTkBgP5DiZ0_Qxow", "your token")

func main() {
	rand.Seed(time.Now().UnixNano()) // инициируем Seed рандома для функции requestID
	flag.Parse()

	stream()
}

func stream() {
	logger := log.New(os.Stdout, "[invest-openapi-go-sdk]", log.LstdFlags)

	client, err := sdk.NewStreamingClient(logger, token)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// Запускаем цикл обработки входящих событий. Запускаем асинхронно
	// Сюда будут приходить сообщения по подпискам после вызова соответствующих методов
	// SubscribeInstrumentInfo, SubscribeCandle, SubscribeOrderbook
	go func() {
		err = client.RunReadLoop(func(event interface{}) error {
			logger.Printf("Got event: '%+v'", event)
			return nil
		})
		if err != nil {
			log.Fatalln(err)
		}
	}()

	log.Println("Подписка на получение событий по инструменту BBG005DXJS36 (TCS)")
	err = client.SubscribeInstrumentInfo("BBG005DXJS36", requestID())
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Подписка на получение свечей по инструменту BBG005DXJS36 (TCS)")
	err = client.SubscribeCandle("BBG005DXJS36", sdk.CandleInterval5Min, requestID())
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Подписка на получения стакана по инструменту BBG005DXJS36 (TCS)")
	err = client.SubscribeOrderbook("BBG005DXJS36", 10, requestID())
	if err != nil {
		log.Fatalln(err)
	}

	// Приложение завершится через 10секунд.
	// Hint: В боевом приложении лучше обрабатывать сигналы завершения и работать в бесконечном цикле
	time.Sleep(10 * time.Second)

	log.Println("Отписка от получения событий по инструменту BBG005DXJS36 (TCS)")
	err = client.UnsubscribeInstrumentInfo("BBG005DXJS36", requestID())
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Отписка от получения свечей по инструменту BBG005DXJS36 (TCS)")
	err = client.UnsubscribeCandle("BBG005DXJS36", sdk.CandleInterval5Min, requestID())
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Отписка от получения стакана по инструменту BBG005DXJS36 (TCS)")
	err = client.UnsubscribeOrderbook("BBG005DXJS36", 10, requestID())
	if err != nil {
		log.Fatalln(err)
	}
}


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
