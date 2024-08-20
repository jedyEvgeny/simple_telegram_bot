package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/jedyEvgeny/simple_telegram_bot/lib/e"
	"github.com/jedyEvgeny/simple_telegram_bot/storage"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774

func New(b string) Storage {
	return Storage{basePath: b}
}

func (s Storage) Save(page *storage.Page) (err error) {
	//Определяем способ обработки ошибок
	defer func() { err = e.WrapIfErr("не смогли сохранить страницу", err) }()

	//Определяем путь до директории, куда будет сохраняться файл
	fPath := filepath.Join(s.basePath, page.UserName)

	//Создаём все нужные директории в определённом пути
	err = os.MkdirAll(fPath, defaultPerm)
	if err != nil {
		return err
	}

	//Формируем имя файла
	fName, err := fileName(page)
	if err != nil {
		return err
	}

	//Дописываем имя файла к полному пути
	fPath = filepath.Join(fPath, fName)

	//Создаём файл
	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	//Преобразуем страницу в формат gob и записываем в указанный файл
	err = gob.NewEncoder(file).Encode(page)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	//Определяем способ обработки ошибок
	defer func() { err = e.WrapIfErr("не смогли выбрать случайную страницу", err) }()
	fPath := filepath.Join(s.basePath, userName)
	
	files, err := os.ReadDir(fPath)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, storage.ErrNoSavePages
	}

	// rand.Seed(time.Now().UnixNano())
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	n := r.Intn(len(files))

	file := files[n]
	fullfPath := filepath.Join(fPath, file.Name())
	return s.decodePage(fullfPath)
}

func (s Storage) Remove(p *storage.Page) error {
	fName, err := fileName(p)
	if err != nil {
		return e.Wrap("не смогли удалить файл", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fName)
	err = os.Remove(path)
	if err != nil {
		msg := fmt.Sprintf("не смогли удалить файл %s", path) //прописываем путь к файлу, который не смогли удалить
		return e.Wrap(msg, err)
	}
	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("не смогли проверить существование файла", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fName)

	//Проверяем существование файла
	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintln("не смогли проверить существование файла ", path)
		return false, e.Wrap(msg, err)
	}
	return true, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

func (s Storage) decodePage(filePath string) (p *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("не смогли декодировать страницу", err) }()
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	//Декодируем файл
	err = gob.NewDecoder(f).Decode(&p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
