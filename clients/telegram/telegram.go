// Клиент для общения с API Телеграм
// Клиент выполняет две вещи: 1. Получение новых сообщений (updates) и 2. Отправка собственных сообщений пользователям
package telegram

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/jedyEvgeny/simple_telegram_bot/lib/e"
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

func newBasePath(token string) string {
	return "bot" + token
}

// Update получает сообщения
func (c *Client) Updates(offset, limit int) ([]Update, error) {
	//Формируем параметры запроса
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset)) //Добавляем параметры к запросу
	q.Add("limit", strconv.Itoa(limit))   //Добавляем параметры к запросу
	log.Println("Создали запрос: ", q)
	//Запрос будет одинаковым для всех методов
}

// SendMessage отправляет сообщения пользователям бота
func (c *Client) SendMessage() {

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
