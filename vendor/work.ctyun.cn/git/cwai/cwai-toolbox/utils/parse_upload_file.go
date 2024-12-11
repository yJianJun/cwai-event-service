package utils

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/extrame/xls"
	"github.com/xuri/excelize/v2"
)

func ParseUploadFile(upload *multipart.FileHeader) ([][]string, error) {
	file, err := upload.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ext := strings.ToLower(upload.Filename[strings.LastIndex(upload.Filename, ".")+1:])
	switch ext {
	case "xls":
		return parseXLS(file)
	case "xlsx":
		return parseXLSX(file)
	case "csv":
		return parseCSV(file)
	default:
		return nil, fmt.Errorf("不支持的文件格式: %s", ext)
	}
}

func parseXLS(file multipart.File) ([][]string, error) {
	xlsFile, err := xls.OpenReader(file, "utf-8")
	if err != nil {
		return nil, fmt.Errorf("无法解析 Excel 文件: %v", err)
	}

	var records [][]string
	sheetCount := xlsFile.NumSheets()
	for sheetIndex := 0; sheetIndex < sheetCount; sheetIndex++ {
		sheet := xlsFile.GetSheet(sheetIndex)
		if sheet == nil {
			continue
		}

		for rowIndex := 0; rowIndex <= int(sheet.MaxRow); rowIndex++ {
			row := sheet.Row(rowIndex)
			if row == nil {
				continue
			}

			var rowData []string
			for colIndex := 0; colIndex < row.LastCol(); colIndex++ {
				cell := row.Col(colIndex)
				rowData = append(rowData, cell)
			}
			records = append(records, rowData)
		}
	}

	return records, nil
}

func parseXLSX(file multipart.File) ([][]string, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("无法解析 XLSX 文件: %v", err)
	}

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("XLSX 文件中没有工作表")
	}
	// 选择第一个工作表
	sheet := sheets[0]
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("无法读取工作表 %s: %v", sheet, err)
	}

	return rows, nil
}

func parseCSV(file multipart.File) ([][]string, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("无法解析 CSV 文件: %v", err)
	}
	return records, nil
}
