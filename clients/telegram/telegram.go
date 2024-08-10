// Клиент для общения с API Телеграм
// Клиент выполняет две вещи: 1. Получение новых сообщений (updates) и 2. Отправка собственных сообщений пользователям
package telegram

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/jedyEvgeny/simple_telegram_bot/lib/e"
)

const (
	getUpdatesMethod  = "getUpdates"  //из документации API-Телеграма
	sendMessageMethod = "sendMessage" //из документации API-Телеграма
)

type Client struct {
	host     string      //хост API-сервиса Телеграм
	basePath string      //Базовый путь, префикс с которого начинаются все запросы, например tg-bot.com/bot<token>
	client   http.Client //Чтобы не создавать для каждого запроса отдельно
}

// New создаёт клиент
func New(h string, token string) Client { //разные типы данных, т.к. токен в теории может быть интеджером
	apiPath := newBasePath(token)
	return Client{
		host:     h,
		basePath: apiPath,
		client:   http.Client{},
	}
}

// Update получает сообщения от пользователей бота
func (c *Client) Updates(offset, limit int) (updates []Update, err error) {
	defer func() { err = e.WrapIfErr("не смогли получить обновления", err) }()
	//Формируем параметры запроса
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset)) //Добавляем параметры к запросу
	q.Add("limit", strconv.Itoa(limit))   //Добавляем параметры к запросу
	log.Println("Создали запрос: ", q)

	//Запрос будет одинаковым для всех методов, поэтому выводим в отдельный метод
	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}
	var res UpdatesResponse
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

// SendMessage отправляет сообщения пользователям бота
func (c *Client) SendMessage(chatID int, text string) error {
	//Подготавливаем параметры запроса
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID)) //Первый аргумент из документации TG-API: https://core.telegram.org/bots/api#sendmessage
	q.Add("text", text)                    //Также

	//Тело ответа нам не нужно
	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("не смогли отправить сообщение", err)
	}
	return nil
}

// newBasePath формирует базовый путь запроса
func newBasePath(token string) string {
	return "bot" + token
}

// doRequest формирует запрос для отправки, отправляет его и получает ответ
func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() {
		err = e.WrapIfErr("не выполнен запрос", err)
	}()
	fullPath := path.Join(c.basePath, method) //убирает лишние слешы или добавляет недостающие
	//формируем URL на который будет отправляться запрос
	u := url.URL{ //Результат выглядит примерно так: https://api.telegram.org/bot1234567890/getUpdates
		Scheme: "https", //протокол
		Host:   c.host,
		Path:   fullPath,
	}

	//формируем объект запроса
	req, err := http.NewRequest(http.MethodGet, u.String(), nil) //Тело запроса пустое, т.к. метод GET обычно без тела,
	// и всё необходимое указано в виде параметров
	// Выглядит примерно так: &{GET https://api.telegram.org/bot1234567890/getUpdates HTTP/1.1 1 1 map[] <nil> <nil> 0 [] false api.telegram.org map[] map[] <nil> map[]   <nil> <nil> <nil> {{}}}

	if err != nil {
		return nil, err
	}

	//Добавляем к объекту запроса req параметры запроса из сигнатуры функции
	req.URL.RawQuery = query.Encode()
	//Пример результата: &{GET https://api.telegram.org/bot1234567890/getUpdates?chat_id=1234567890&text=Hello%2C+world%21 HTTP/1.1 1 1 map[] <nil> <nil> 0 [] false api.telegram.org map[] map[] <nil> map[]   <nil> <nil> <nil> {{}}}

	//отправляем получившийся запрос в ожидании ответа
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Получаем содержимое ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
