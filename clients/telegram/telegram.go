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
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host     string      //хост API-сервиса Телеграм
	basePath string      //Базовый путь, префикс с которого начинаются все запросы, например tg-bot.com/bot<token>
	client   http.Client //Чтобы не создавать для каждого запроса отдельно
}

// New создаёт клиент
func New(h string, token string) Client { //разные типы данных, т.к. токен в теории может быть интеджером
	url := newBasePath(token)
	return Client{
		host:     h,
		basePath: url,
		client:   http.Client{},
	}
}

// Update получает сообщения
func (c *Client) Updates(offset, limit int) (updates []Update, err error) {
	defer func() { err = e.WrapIfErr("не смогли получить обновления", err) }()
	//Формируем параметры запроса
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset)) //Добавляем параметры к запросу
	q.Add("limit", strconv.Itoa(limit))   //Добавляем параметры к запросу
	log.Println("Создали запрос: ", q)
	//Запрос будет одинаковым для всех методов
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

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() {
		err = e.WrapIfErr("не выполнен запрос", err)
	}()
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil) //Тело запроса пустое, т.к. метод GET обычно без тела,
	// и всё необходимое указано в виде параметров
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = query.Encode() //Передаём объект req (реквест) в параметры запроса, полученные в аргументе сигнатуры

	//отправляем получившийся запрос
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
