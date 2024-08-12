package telegram

import (
	"github.com/jedyEvgeny/simple_telegram_bot/clients/telegram"
	"github.com/jedyEvgeny/simple_telegram_bot/events"
	"github.com/jedyEvgeny/simple_telegram_bot/lib/e"
	"github.com/jedyEvgeny/simple_telegram_bot/storage"
)

type Proceccor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

func New(client *telegram.Client, storage storage.Storage) *Proceccor {
	return &Proceccor{
		tg:      client,
		storage: storage,
	}
}

func (p *Proceccor) Fetch(limit int) ([]events.Event, error) {
	//С помощью клиента получаем все апдейты
	update, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("не смогли получить события", err)
	}

	//Аллоцируем память под результат
	res := make([]events.Event, 0, len(update))

	//Обходим все апдейты и преобразуем их в эвенты
	for _, u := range update {
		e := event(u)
		res = append(res, e)
	}
}

func event(upd telegram.Update) events.Event {
	res := events.Event{
		Type: fetchType(upd),
		Text: fetchText(upd),
	}
}

func fetchType(upd telegram.Update) string {

}

func fetchText(upd telegram.Update) string {

}
