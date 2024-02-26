package server

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/xuri/excelize/v2"
)

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
}
