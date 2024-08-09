package main

import (
	"flag"
	"log"

	"github.com/jedyEvgeny/simple_telegram_bot/clients/telegram"
)

func main() {
	//Структура проекта:
	//tgClient = telegram.New(tocken) //Для общения с API Телеграма фетчеру и процессору потребуется клиент
	//tocken - строка, которую получаем от Телеграма, и которую передаём клиенту
	//fetcher = fetcher.New(tgClient) - нужен для получения событий - отправляет запросы в API телеграма для получения новых событий
	//processor = processor.New(tgClient) - нужен для обработки событий и будет отправлять нам новые сообщения (в боте)
	//consumer.Start(fetcher, processor) //для получения и обработки событий
	//Фетчер и процессор будут общаться с API телеграма
	t := mustToken()
	h := mustHost() //для гибкости приложения хост не константный
	tgClient := telegram.New(h, t)
}

// приставка must делается для функций, которые вместо возвращения ошибки,
// аварийно завершают программу
// В основном применяется для запуска программы и парсинга конфигов
func mustToken() string {
	token := flag.String( //ссылка на функции
		"token-bot-token",
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
	var host string
	flag.StringVar(
		&host,
		"host-bot-host",
		"api.telegram.org",
		"хост API-сервиса Телеграм",
	)
	log.Println("Выбран хост: ", host)
	return host
}
