package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"

	"github.com/glebarez/sqlite"
	"github.com/xuri/excelize/v2"

	"gorm.io/gorm"
)

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

	db.Create(NewUser(PrivilegeAdmin, "admin", "minda"))

	return &storage{db: db}, nil
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

	c := &columns{}

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

// Получить всю таблицу
func (s *storage) getTable(w http.ResponseWriter, r *http.Request) {
	var rows []tableRow

	result := s.db.Find(&rows)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusBadRequest)
	}

	arr := rowsToString(rows)
	data, _ := json.Marshal(arr)
	w.Write(data)
}

// Удалить запись
func (s *storage) deleteRow(w http.ResponseWriter, r *http.Request) {
	idData := r.FormValue("id")
	if idData == "" {
		http.Error(w, "no id received", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	row := tableRow{Id: Id(id)}

	result := s.db.Delete(&row).Where("id")
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusBadRequest)
	}
}

func rowsToString(rows []tableRow) [][]string {
	var rowsArr [][]string

	for _, row := range rows {
		// Преобразуем каждую структуру в массив строк
		data := rowToString(row)

		// Добавляем массив строк в итоговый массив
		rowsArr = append(rowsArr, data)
	}

	return rowsArr
}

// Строка в массив
func rowToString(s interface{}) []string {
	// Получаем тип структуры
	t := reflect.TypeOf(s)

	// Получаем значение структуры
	v := reflect.ValueOf(s)

	// Создаем массив строк
	var paramArr []string

	// Перебираем все поля структуры
	for i := 0; i < t.NumField(); i++ {
		// Получаем значение поля
		fieldValue := v.Field(i)

		// typeName := fieldValue.Type().Name()
		// if typeName == "Id" {
		// 	continue
		// }

		value := fieldValue.Interface()
		data := fmt.Sprintf("%v", value)

		// Преобразуем значение поля в строку и добавляем его в массив
		paramArr = append(paramArr, fmt.Sprintf("%v", data))
	}

	return paramArr
}
