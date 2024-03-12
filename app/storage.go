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
	Region        int
	Responsible   int
	Verified      int
	VkUrl         int
	OkUrl         int
	TgUrl         int
	Reason        int
	CommentaryNpa int
	FullName      int
	Ogrn          int
	Status        int
	Commentary    int
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
		cell, err := excel.GetCellValue("Томская область", request)

		if err != nil {
			return nil, err
		}

		switch {
		case cell == region:
			hs.Region = i
		case cell == responsible:
			hs.Responsible = i
		case cell == verified:
			hs.Verified = i
		case cell == vkUrl:
			hs.VkUrl = i
		case cell == okUrl:
			hs.OkUrl = i
		case cell == tgUrl:
			hs.TgUrl = i
		case cell == reason:
			hs.Reason = i
		case cell == commentaryNpa:
			hs.CommentaryNpa = i
		case cell == fullName:
			hs.FullName = i
		case cell == ogrn:
			hs.Ogrn = i
		case cell == status:
			hs.Status = i
		case cell == commentary:
			hs.Commentary = i
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

	cell, err := excel.GetCellValue("Томская область", "A1")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, cell)
	hs, err := s.excelToHeaders(excel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, hs)
}
