package main

import (
	"flag"
	"log"

	tgClient "github.com/jedyEvgeny/simple_telegram_bot/clients/telegram"
	event_consumer "github.com/jedyEvgeny/simple_telegram_bot/consumer/event-consumer"
	"github.com/jedyEvgeny/simple_telegram_bot/events/telegram"
	"github.com/jedyEvgeny/simple_telegram_bot/storage/files"
)

const (
	storagePath = "files_storage" //вынести в конфиг
	host        = "api.telegram.org"
	batchSize   = 100 //Размер пачки
)

func main() {
	//Структура проекта:
	//+ tgClient = telegram.New(tocken) //Для общения с API Телеграма фетчеру и процессору потребуется клиент
	//+ tocken - строка, которую получаем от Телеграма, и которую передаём клиенту
	//fetcher = fetcher.New(tgClient) - нужен для получения событий - отправляет запросы в API телеграма для получения новых событий
	//processor = processor.New(tgClient) - нужен для обработки событий и будет отправлять нам новые сообщения (в боте)
	//consumer.Start(fetcher, processor) //для получения и обработки событий
	//Фетчер и процессор будут общаться с API телеграма

	t := mustToken() //Получаем токен бота через консоль
	// h := mustHost()  //для гибкости приложения хост не константный

	eventsProseccor := telegram.New(
		tgClient.New(host, t),
		files.New(storagePath),
	)

	log.Println("сервис запущен")

	//Запускаем консьюмера
	consumer := event_consumer.New(eventsProseccor, eventsProseccor, batchSize)

	err := consumer.Start()
	//Ошибка только если косньюмер аварийно остановился
	if err != nil {
		log.Fatal()
	}
}

// приставка must делается для функций, которые вместо возвращения ошибки,
// аварийно завершают программу
// В основном применяется для запуска программы и парсинга конфигов
func mustToken() string {
	token := flag.String( //ссылка на функции
		"tg-bot-token",
		"", //значение по-умолчанию задаём пустое, т.к. токен - обязательный
		"токен для доступа в телеграм", //подсказка к флагу, видимая после компиляции
	)
	flag.Parse() //Кладём значение в token
	if *token == "" {
		log.Fatal("Токен не указан")
	}
	return *token
}

func mustHost() string {
	host := flag.String(
		"tg-bot-host",
		"api.telegram.org",
		"хост API-сервиса Телеграм",
	)
	log.Println("Выбран хост: ", host)
	if *host == "" {
		log.Fatal("Токен не указан")
	}
	return *host
}

// func mustHost() string {
// 	var host string
// 	flag.StringVar(
// 		&host,
// 		"tg-bot-host",
// 		"api.telegram.org",
// 		"хост API-сервиса Телеграм",
// 	)
// 	log.Println("Выбран хост: ", host)
// 	return host
// }
