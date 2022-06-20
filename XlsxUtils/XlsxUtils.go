package XlsxUtils

import (
	"bytes"
	"os"
	"strconv"

	"github.com/zhaobisheng/bstool/FileUtils"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func GenerateCommonXlsxStream(sheetName string, title []string, dataMap [][]string) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	f.SetSheetName("Sheet1", sheetName)
	mainRowStart := 2
	titleStart := 'A'
	for k, v := range title {
		f.SetCellValue(sheetName, string(titleStart+rune(k))+"1", v)
	}
	for index, dataRow := range dataMap {
		for colIndex, data := range dataRow {
			f.SetCellValue(sheetName, string(titleStart+rune(colIndex))+strconv.Itoa(mainRowStart+index), data)
		}
	}
	return f.WriteToBuffer()
}

func GenerateCommonXlsx(newFilePath, sheetName string, title []string, dataMap [][]interface{}) error {
	f := excelize.NewFile()
	f.SetSheetName("Sheet1", sheetName)
	mainRowStart := 2
	titleStart := 'A'
	for k, v := range title {
		f.SetCellValue(sheetName, string(titleStart+rune(k))+"1", v)
	}
	for index, dataRow := range dataMap {
		for colIndex, data := range dataRow {
			f.SetCellValue(sheetName, string(titleStart+rune(colIndex))+strconv.Itoa(mainRowStart+index), data)
		}
	}
	fileDir := FileUtils.GetFileDir(newFilePath)
	if !FileUtils.PathExists(fileDir) {
		os.MkdirAll(fileDir, os.ModePerm)
	}
	err := f.SaveAs(newFilePath)
	if err != nil {
		return err
	}
	return nil
}

func GenerateCommonXlsxByString(newFilePath, sheetName string, title []string, dataMap [][]string) error {
	f := excelize.NewFile()
	f.SetSheetName("Sheet1", sheetName)
	mainRowStart := 2
	titleStart := 'A'
	for k, v := range title {
		f.SetCellValue(sheetName, string(titleStart+rune(k))+"1", v)
	}
	for index, dataRow := range dataMap {
		for colIndex, data := range dataRow {
			f.SetCellValue(sheetName, string(titleStart+rune(colIndex))+strconv.Itoa(mainRowStart+index), data)
		}
	}
	fileDir := FileUtils.GetFileDir(newFilePath)
	if !FileUtils.PathExists(fileDir) {
		os.MkdirAll(fileDir, os.ModePerm)
	}
	err := f.SaveAs(newFilePath)
	if err != nil {
		return err
	}
	return nil
}

func GenerateComplexXlsx(newFilePath string, title []string, dataMap map[string][][]interface{}) error {
	f := excelize.NewFile()
	mainRowStart := 2
	for sheetName, val := range dataMap {
		f.NewSheet(sheetName)
		titleStart := 'A'
		for k, v := range title {
			//titleStart = titleStart + 1
			f.SetCellValue(sheetName, string(titleStart+rune(k))+"1", v)
		}
		for index, dataRow := range val {
			for colIndex, data := range dataRow {
				f.SetCellValue(sheetName, string(titleStart+rune(colIndex))+strconv.Itoa(mainRowStart+index), data)
			}
		}
	}
	f.DeleteSheet("Sheet1")
	fileDir := FileUtils.GetFileDir(newFilePath)
	if !FileUtils.PathExists(fileDir) {
		os.MkdirAll(fileDir, os.ModePerm)
	}
	err := f.SaveAs(newFilePath)
	if err != nil {
		return err
	}
	return nil
}

func GenerateFileTreeXlsx(newFilePath string, title []string, dataMap map[string][][]interface{}) error {
	f := excelize.NewFile()
	mainRowStart := 2
	for sheetName, val := range dataMap {
		if sheetName == "未知" {
			f.SetSheetName("Sheet1", sheetName)
		} else {
			f.NewSheet(sheetName)
		}
		titleStart := 'A'
		for k, v := range title {
			//titleStart = titleStart + 1
			f.SetCellValue(sheetName, string(titleStart+rune(k))+"1", v)
		}
		for index, dataRow := range val {
			for colIndex, data := range dataRow {
				f.SetCellValue(sheetName, string(titleStart+rune(colIndex))+strconv.Itoa(mainRowStart+index), data)
			}
		}
	}
	f.DeleteSheet("Sheet1")
	fileDir := FileUtils.GetFileDir(newFilePath)
	if !FileUtils.PathExists(fileDir) {
		os.MkdirAll(fileDir, os.ModePerm)
	}
	err := f.SaveAs(newFilePath)
	if err != nil {
		return err
	}
	return nil
}
