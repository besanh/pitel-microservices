package common

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/xuri/excelize/v2"
)

func HandleExcelStreamWriter(headers []string, rows [][]string) (buffer *bytes.Buffer, err error) {
	excelFile := excelize.NewFile()
	SHEET1 := "Sheet1"
	index := 1
	streamWriter, err := excelFile.NewStreamWriter(SHEET1)
	if err != nil {
		log.Error(err)
		return
	}
	styleID, _ := excelFile.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FCD5B4"}, Pattern: 1},
		Alignment: &excelize.Alignment{WrapText: true},
	})

	// write headers row
	cell, _ := excelize.CoordinatesToCellName(1, index)
	values := []any{}
	for _, header := range headers {
		values = append(values, excelize.Cell{
			Value:   header,
			StyleID: styleID,
		})
	}
	if err = streamWriter.SetRow(cell, values); err != nil {
		log.Error(err)
		return
	}
	index++
	styleID, _ = excelFile.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#ffffff"}, Pattern: 1},
		Alignment: &excelize.Alignment{WrapText: false, Horizontal: "left"},
	})
	for _, row := range rows {
		cell, _ := excelize.CoordinatesToCellName(1, index)
		values := []any{}
		for _, cellValue := range row {
			values = append(values, excelize.Cell{
				Value:   cellValue,
				StyleID: styleID,
			})
		}
		if err := streamWriter.SetRow(cell, values); err != nil {
			log.Error(err)
			break
		}
		index++
	}
	if err = streamWriter.Flush(); err != nil {
		log.Error(err)
		return
	}

	buffer, err = excelFile.WriteToBuffer()
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func HandleCSVStreamWriter(fileName string, headers []string, rows [][]string) (buffer *bytes.Buffer, err error) {
	// Create a temporary file to store the CSV content
	tmpFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return
	}
	defer os.Remove(tmpFile.Name())
	// Write the CSV content to the temporary file
	w := csv.NewWriter(tmpFile)
	w.UseCRLF = true
	bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
	_, err = tmpFile.Write(bomUtf8)
	if err != nil {
		log.Error(err)
		return
	}

	if err = w.Write(headers); err != nil {
		log.Error(err)
		return
	}

	for _, row := range rows {
		records := make([]string, len(row))
		for k, v := range row {
			records[k] = fmt.Sprintf("%v", v)
		}
		if err = w.Write(records); err != nil {
			log.Error(err)
			return
		}
	}
	w.Flush()
	if err = tmpFile.Close(); err != nil {
		log.Error(err)
		return
	}

	// Convert the temporary file content to a *bytes.Buffer
	fileBytes, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		log.Error(err)
		return
	}

	buffer = bytes.NewBuffer(fileBytes)
	return
}
