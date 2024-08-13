package telegram

import (
	"errors"

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

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUncknownEventType = errors.New("неизвестный тип события")
	ErrEmptyUpdates      = errors.New("внутренняя ошибка")
	ErrUncknownMeta      = errors.New("неизвестный тип Meta")
)

func New(client *telegram.Client, storage storage.Storage) *Proceccor {
	return &Proceccor{
		tg:      client,
		storage: storage,
	}
}

func (p *Proceccor) Fetch(limit int) ([]events.Event, error) {
	//С помощью клиента получаем все апдейты (апдейты характерны именно для мессенджера Telegram)
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("не смогли получить события", err)
	}

	//Возвращаем нулевой результат, если не нашли апдейтов
	if len(updates) == 0 {
		return nil, e.Wrap("список апдейтов пустой:", ErrEmptyUpdates)
	}

	//Аллоцируем память под результат. Тип Event - более общая сущность, чем апдейты
	//В Event мы можем преобразовывать всё что получаем от потенциальных других мессенджеров
	//в каком бы формате они не предоставляли информацию
	res := make([]events.Event, 0, len(updates))

	//Перебираем все апдейты и преобразуем их в тип Event пакета events
	for _, u := range updates {
		e := event(u)
		res = append(res, e)
	}

	//Обновляем значение внутреннего поля offset, чтобы в следующий раз получить следующую пачку изменений
	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

// Process выполняет различные действия в зависимости от типа эвента
func (p *Proceccor) Process(event events.Event) error {
	switch event.Type {
	case events.Message: //когда работаем с сообщением
		p.ProcessMessage(event)
	default: //Когда не знаем с чем работаем
		return e.Wrap("не смогли выполнить сообщение", ErrUncknownEventType)
	}
}

func (p *Proceccor) ProcessMessage(event events.Event) {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("не смогли выполнить сообщение", err)
	}
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("не смогли получить Meta")
	}
	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: fetchType(upd),
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown //Если сообщение нулевое, то тип неизвестен
	}
	//Если сообщение не нулевое, то событие - это сообщение
	return events.Message
}
