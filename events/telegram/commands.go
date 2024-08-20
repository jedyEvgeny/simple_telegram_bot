package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/jedyEvgeny/simple_telegram_bot/lib/e"
	"github.com/jedyEvgeny/simple_telegram_bot/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

// doCmd смотрит на текст сообщения и по его формату и содержанию определяет, какая это команда
func (p *Proceccor) doCmd(text string, chatID int, username string) error {
	//Удаляем пробелы из текста сообщения, т.к. они будут нам мешать
	text = strings.TrimSpace(text)

	log.Printf("получена новая комманда '%s' от '%s'\n", text, username)

	//Структура ожидаемых команд:
	//1.Сохранить страницу - в сообщении что-то вроде https://...
	//2.Получить страницу - в сообщении команда /rnd
	//3.Получить помощь по работе с  ботом - в сообщении написать /help
	//4.Команда /start отправляется боту автоматически, когда человек начинает с ним общение
	//При старте мы отправляем человеку приветствие и тот же текст, что отправляем в /help

	if isAddCmd(text) {
		//TODO: AddPage()
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}

}

func (p *Proceccor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() {
		err = e.WrapIfErr("не смогли выполнить команду 'сохранить страницу'", err)
	}()

	//подготавливаем страницу, которую собираемся сохранить
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	//Проверяем, не существует ли такая страница
	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExicts)
	}

	//Сохраняем страницу
	err = p.storage.Save(page)
	if err != nil {
		return err
	}

	//Если страница корректно сохранился, сообщаем об этом пользователю
	err = p.tg.SendMessage(chatID, msgSave)
	if err != nil {
		return err
	}

	return nil
}

func (p *Proceccor) sendRandom(chatID int, username string) (err error) {
	defer func() {
		err = e.WrapIfErr("не смогли выполнить команду 'отправка случайной страницы'", err)
	}()
	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavePages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavePages) {
		return p.tg.SendMessage(chatID, msgNoSavePages)
	}

	err = p.tg.SendMessage(chatID, page.URL)
	if err != nil {
		return err
	}

	//Если отправили ссылку, то её нужно удалить из хранилища
	return p.storage.Remove(page)
}

func (p *Proceccor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Proceccor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

// isAddCmd проверяет, является ли текст сообщения ожидаемым для сохранения параметром (ссылкой на сайт)
func isAddCmd(text string) bool {
	return isUrl(text)
}

// isUrl проверяет, является ли текст сообщения ссылкой на веб-страницу
func isUrl(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
