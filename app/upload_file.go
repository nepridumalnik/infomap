package app

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/xuri/excelize/v2"
)

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

type TableRow struct {
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

// Извлечь excel файл из multipart данных
func extractExcel(file *multipart.File) (*excelize.File, error) {
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

func excelToHeaders(excel *excelize.File) (*TableRow, error) {
	const columnsCount = 12

	hs := new(TableRow)

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
func upload(w http.ResponseWriter, r *http.Request) {
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

	excel, err := extractExcel(&file)
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
	hs, err := excelToHeaders(excel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, hs)
}
