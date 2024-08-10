package telegram

import "github.com/jedyEvgeny/simple_telegram_bot/clients/telegram"

type Proceccor struct {
	tg     *telegram.Client
	offset int
	//storage для хранения сообщений (ссылок)
}

func New(client *telegram.Client) {

}
