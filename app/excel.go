package app

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

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
