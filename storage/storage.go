//Кладём интерфейс
//За счёт этого, storage может работать и с файловой системой, и с БД и с чем угодно

package storage

import (
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/jedyEvgeny/simple_telegram_bot/lib/e"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (hash string, err error) {
	defer func() { err = e.Wrap("не смогли рассчитать хеш", err) }()
	h := sha1.New()
	_, err = io.WriteString(h, p.URL)
	if err != nil {
		return "", err
	}
	_, err = io.WriteString(h, p.UserName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
