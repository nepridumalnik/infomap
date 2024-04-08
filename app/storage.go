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

// Индексы столбцов таблицы
type columns struct {
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

type tableRows []*tableRow

// Строка таблицы
type tableRow struct {
	Id            Id `gorm:"primaryKey"`
	Region        string
	Responsible   string
	Verified      string
	VkUrl         string
	OkUrl         string
	TgUrl         string
	Reason        string
	CommentaryNpa string
	FullName      string
	Ogrn          string
	Status        string
	Commentary    string
}

// Хранилище данных
type storage struct {
	db *gorm.DB
}

// Создать хранилище
func NewStorage() (*storage, error) {
	db, err := gorm.Open(sqlite.Open(sqliteDbFile))

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&tableRow{}, &User{}, &session{})
	if err != nil {
		return nil, err
	}

	db.Create(NewUser(Admin, "admin", "minda"))

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
func (s *storage) excelToColumns(excel *excelize.File) (*columns, error) {
	const columnsCount = 12

	c := new(columns)

	for i := 0; i < columnsCount; i++ {
		request := fmt.Sprintf("%c1", 'A'+i)
		cell, err := excel.GetCellValue(mainList, request)

		if err != nil {
			return nil, err
		}

		switch {
		case cell == region:
			c.region = i
		case cell == responsible:
			c.responsible = i
		case cell == verified:
			c.verified = i
		case cell == vkUrl:
			c.vkUrl = i
		case cell == okUrl:
			c.okUrl = i
		case cell == tgUrl:
			c.tgUrl = i
		case cell == reason:
			c.reason = i
		case cell == commentaryNpa:
			c.commentaryNpa = i
		case cell == fullName:
			c.fullName = i
		case cell == ogrn:
			c.ogrn = i
		case cell == status:
			c.status = i
		case cell == commentary:
			c.commentary = i
		default:
			return nil, errors.New("unknown cell name")
		}
	}

	return c, nil
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

	c, err := s.excelToColumns(excel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.moveRowsToDb(c, excel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// Получить строки из excel
func (s *storage) extractRows(c *columns, excel *excelize.File) (tableRows, error) {
	rows := make([]*tableRow, 0)

	for i := 2; ; i++ {
		request := fmt.Sprintf("%c%d", 'A'+c.fullName, i)
		cell, err := excel.GetCellValue(mainList, request)

		if err != nil || len(cell) == 0 {
			break
		}

		row, err := s.getRowByIdx(c, excel, i)

		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, nil
}

// Обходит построчно таблицу и заносит данные в БД
func (s *storage) moveRowsToDb(c *columns, excel *excelize.File) error {
	rows, err := s.extractRows(c, excel)

	if err != nil {
		return err
	}

	result := s.db.CreateInBatches(rows, len(rows))

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Получить строку из excel
func (s *storage) getRowByIdx(c *columns, excel *excelize.File, idx int) (*tableRow, error) {
	row := &tableRow{}

	var err error

	row.Region, err = excel.GetCellValue(mainList, fmt.Sprintf("%c%d", 'A'+c.region, idx))
	if err != nil {
		return nil, err
	}
	row.Responsible, err = excel.GetCellValue(mainList, fmt.Sprintf("%c%d", 'A'+c.responsible, idx))
	if err != nil {
		return nil, err
	}
	row.Verified, err = excel.GetCellValue(mainList, fmt.Sprintf("%c%d", 'A'+c.verified, idx))
	if err != nil {
		return nil, err
	}
	row.VkUrl, err = excel.GetCellValue(mainList, fmt.Sprintf("%c%d", 'A'+c.vkUrl, idx))
	if err != nil {
		return nil, err
	}
	row.OkUrl, err = excel.GetCellValue(mainList, fmt.Sprintf("%c%d", 'A'+c.okUrl, idx))
	if err != nil {
		return nil, err
	}
	row.TgUrl, err = excel.GetCellValue(mainList, fmt.Sprintf("%c%d", 'A'+c.tgUrl, idx))
	if err != nil {
		return nil, err
	}
	row.Reason, err = excel.GetCellValue(mainList, fmt.Sprintf("%c%d", 'A'+c.reason, idx))
	if err != nil {
		return nil, err
	}
	row.CommentaryNpa, err = excel.GetCellValue(mainList, fmt.Sprintf("%c%d", 'A'+c.commentaryNpa, idx))
	if err != nil {
		return nil, err
	}
	row.FullName, err = excel.GetCellValue(mainList, fmt.Sprintf("%c%d", 'A'+c.fullName, idx))
	if err != nil {
		return nil, err
	}
	row.Ogrn, err = excel.GetCellValue(mainList, fmt.Sprintf("%c%d", 'A'+c.ogrn, idx))
	if err != nil {
		return nil, err
	}
	row.Status, err = excel.GetCellValue(mainList, fmt.Sprintf("%c%d", 'A'+c.status, idx))
	if err != nil {
		return nil, err
	}
	row.Commentary, err = excel.GetCellValue(mainList, fmt.Sprintf("%c%d", 'A'+c.commentary, idx))
	if err != nil {
		return nil, err
	}

	return row, nil
}
