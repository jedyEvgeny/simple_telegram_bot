package event_consumer

import (
	"log"
	"sync"
	"time"

	"github.com/jedyEvgeny/simple_telegram_bot/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int //Размер пачки говорит о том, сколько событий будем обрабатывать за раз
}

func New(f events.Fetcher, p events.Processor, b int) Consumer {
	return Consumer{
		fetcher:   f,
		processor: p,
		batchSize: b,
	}
}

var (
	wg sync.WaitGroup
)

func (c Consumer) Start() error {
	//Постоянно ждём новые события и обрабатываем их
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		//Возможная причина ошибки - проблема с сетью.
		//Если попытка получить данные раз в час или около того
		//имеет смысл выполнить повторную попытку достучаться до сервера, а не ждать следующий час
		if err != nil {
			log.Printf("[ERR] консьюмер %s", err.Error())
			continue
		}

		//Если события не получили, то пропускаем итерацию
		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		err = c.handleEvents(gotEvents)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

//Проблема 1: потеря событий. Решения: ретраи - не страхует на 100%; возвращение в хранилище, фоллбеки, подтверждение для фетчера
// Проблема 2: обработка всей пачки. Если проблема будет повторяться, то скорее всего проблема с сетью
//Решаем либо остановкой после первой ошибки, либо остановкой после определённого количества ошибок
// Проблема 3: параллельная обработка

func (c *Consumer) handleEvents(events []events.Event) error {
	wg.Add(len(events))
	for _, event := range events {
		go c.handleOneEvent(event)
	}
	wg.Wait()
	return nil
}

func (c *Consumer) handleOneEvent(event events.Event) {
	defer wg.Done()
	log.Printf("получено новое событие: %s\n", event.Text)

	//Обрабатываем событие с помощью процессора
	err := c.processor.Process(event)
	if err != nil {
		log.Println("не смогли обработать событие: ", err.Error())
		return
	}
}
