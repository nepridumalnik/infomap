package app

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/glebarez/sqlite"
	"github.com/xuri/excelize/v2"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Столбцы таблицы
const (
	mainList      = "Томская область"
	region        = "Регион"
	responsible   = "Назначен ответственный"
	verified      = "Страница подтверждена"
	vkUrl         = "Ссылка на официальную страницу Вконтакте"
	okUrl         = "Ссылка на официальную страницу Одноклассники"
	tgUrl         = "Ссылка на официальную страницу Telegram"
	reason        = "Официальная страница не ведется на основании"
	commentaryNpa = "Комментарий по НПА"
	fullName      = "Полное наименование объекта"
	ogrn          = "ОГРН"
	status        = "Статус"
	commentary    = "Комментарий"
)

// Строка таблицы
type tableRow struct {
	region        int
	responsible   int
	verified      int
	vkUrl         int
	okUrl         int
	tgUrl         int
	reason        int
	commentaryNpa int
	fullName      int
	ogrn          int
	status        int
	commentary    int
}

// Хранилище данных
type storage struct {
	db *gorm.DB
}

// Создать хранилище
func NewStorage() (*storage, error) {
	db, err := gorm.Open(sqlite.Open("data.db"))

	if err != nil {
		return nil, err
	}

	return &storage{db: db}, nil
}

// Регистрация обработчиков
func (s *storage) RegisterHandlers(r *mux.Route) {
	r.HandlerFunc(s.upload)
}

// Извлечь excel файл из multipart данных
func (s *storage) extractExcel(file *multipart.File) (*excelize.File, error) {
	data, err := io.ReadAll(*file)

	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(data)
	excel, err := excelize.OpenReader(reader)

	if err != nil {
		return nil, err
	}

	return excel, nil
}

// Конвертировать excel файл в заголовки
func (s *storage) excelToHeaders(excel *excelize.File) (*tableRow, error) {
	const columnsCount = 12

	hs := new(tableRow)

	for i := 0; i < columnsCount; i++ {
		request := fmt.Sprintf("%c1", 'A'+i)
		cell, err := excel.GetCellValue(mainList, request)

		if err != nil {
			return nil, err
		}

		switch {
		case cell == region:
			hs.region = i
		case cell == responsible:
			hs.responsible = i
		case cell == verified:
			hs.verified = i
		case cell == vkUrl:
			hs.vkUrl = i
		case cell == okUrl:
			hs.okUrl = i
		case cell == tgUrl:
			hs.tgUrl = i
		case cell == reason:
			hs.reason = i
		case cell == commentaryNpa:
			hs.commentaryNpa = i
		case cell == fullName:
			hs.fullName = i
		case cell == ogrn:
			hs.ogrn = i
		case cell == status:
			hs.status = i
		case cell == commentary:
			hs.commentary = i
		default:
			return nil, errors.New("unknown cell name")
		}
	}

	return hs, nil
}

// Загрузить excel файл на сервер
func (s *storage) upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(limitation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	excel, err := s.extractExcel(&file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hs, err := s.excelToHeaders(excel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.moveRowsToDb(hs, excel)
}

// Обходит построчно таблицу и заносит данные в бд
func (s *storage) moveRowsToDb(hs *tableRow, excel *excelize.File) {
	for i := 2; ; i++ {
		request := fmt.Sprintf("%c%d", 'A'+hs.fullName, i)
		cell, err := excel.GetCellValue(mainList, request)

		if err != nil || len(cell) == 0 {
			return
		}

		fmt.Println(cell)
	}
}