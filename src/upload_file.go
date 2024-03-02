package server

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

type headers struct {
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

func excelToHeaders(excel *excelize.File) (*headers, error) {
	const columnsCount = 12

	hs := new(headers)

	for i := 0; i < columnsCount; i++ {
		request := fmt.Sprintf("%c1", 'A'+i)
		cell, err := excel.GetCellValue("Томская область", request)

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
