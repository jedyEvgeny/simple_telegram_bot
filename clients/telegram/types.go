package telegram

// Структура с API телеграм-бота https://core.telegram.org/bots/api#getting-updates
type Update struct {
	ID      int    `json:"update_id"`
	Message string `json:"message"`
}

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}
