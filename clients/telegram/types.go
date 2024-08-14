package telegram

// Структура с API телеграм-бота https://core.telegram.org/bots/api#getting-updates
type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"` //Поле Message может отсутствовать, и мы получим nil. Поэтому ссылка на структуру
}

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

//Для понимания от кого сообщение
type From struct {
	Username string `json:"username"`
}

//Для понимания, как отправить ответное сообщение
type Chat struct {
	ID int `json:"id"`
}
