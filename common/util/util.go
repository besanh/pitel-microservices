package util

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/xuri/excelize/v2"
)

const MAX_LIMIT = 50_000

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func ParseLimit(limit any) int {
	var val = 10
	var err error
	if reflect.TypeOf(limit).Kind() == reflect.String {
		str := ParseString(limit)
		val, err = strconv.Atoi(str)
		if err != nil {
			val = 10
		}
	} else {
		val = ParseInt(fmt.Sprintf("%d", limit))
	}
	if val < 0 {
		val, _ = strconv.Atoi(fmt.Sprintf("%s", limit))
	} else if val > MAX_LIMIT {
		val = MAX_LIMIT
	}
	return val
}

func ParseOffset(offset any) int {
	var val = 0
	var err error
	if reflect.TypeOf(offset).Kind() == reflect.String {
		str := ParseString(offset)
		val, err = strconv.Atoi(str)
		if err != nil {
			val = 0
		}
	} else {
		val = ParseInt(fmt.Sprintf("%d", offset))
	}
	if val < 0 {
		val = 0
	}
	return val
}

func ParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func ParseAnyToString(value any) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ParseStringToAny(value string, dest any) error {
	if err := json.Unmarshal([]byte(value), dest); err != nil {
		return err
	}
	return nil
}

func ParseAnyToAny(value any, dest any) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, dest); err != nil {
		return err
	}
	return nil
}

func ParseString(value any) string {
	str, ok := value.(string)
	if !ok {
		return str
	}
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Trim(str, "\r\n")
	str = strings.TrimSpace(str)
	return str
}

func InArray(item any, array any) bool {
	arr := reflect.ValueOf(array)
	if arr.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}
	return false
}

func GenerateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func ParseToAnyArray(value []string) []any {
	result := make([]any, 0)
	for _, v := range value {
		result = append(result, v)
	}
	return result
}

func TimeToString(valueTime time.Time) string {
	return TimeToStringLayout(valueTime, "2006-01-02 15:04:05")
}

func TimeToStringLayout(valueTime time.Time, layout string) string {
	return valueTime.Format(layout)
}

func InArrayContains(item string, array []string) bool {
	for _, v := range array {
		if strings.Contains(item, v) {
			return true
		}
	}
	return false
}

func ParseStructToByte(data any) ([]byte, error) {
	value, err := ParseAnyToString(data)
	if err != nil {
		return nil, err
	}
	byteSlice := []byte(value)
	return byteSlice, nil
}

func ParseQueryArray(slice []string) []string {
	result := make([]string, 0)
	for _, v := range slice {
		if len(v) > 0 {
			result = append(result, v)
		}
	}
	return result
}

func ContainKeywords(content string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(content, keyword) {
			return true
		}
	}
	return false
}

func SliceToMap[T comparable](slice []T) map[T]bool {
	set := make(map[T]bool)
	for _, v := range slice {
		set[v] = true
	}
	return set
}

func HandleExcelStreamWriter(headers [][]string, rows [][]string, mergedColumns ...string) (buffer *bytes.Buffer, err error) {
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
		Alignment: &excelize.Alignment{WrapText: true, Horizontal: "center", Vertical: "center"},
	})

	if err = streamWriter.SetColWidth(1, 16, 12); err != nil {
		log.Error(err)
		return
	}
	// write headers row
	for _, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(1, index)
		values := make([]any, 0)
		for _, cellValue := range header {
			values = append(values, excelize.Cell{
				Value:   cellValue,
				StyleID: styleID,
			})
		}

		if err = streamWriter.SetRow(cell, values); err != nil {
			log.Error(err)
			return
		}
		index++
	}

	// merge columns
	for i := 0; i+1 < len(mergedColumns); i += 2 {
		if err = streamWriter.MergeCell(mergedColumns[i], mergedColumns[i+1]); err != nil {
			log.Error(err)
			return
		}
	}

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

func HandleCSVStreamWriter(fileName string, headers [][]string, rows [][]string, mergedColumns ...string) (buffer *bytes.Buffer, err error) {
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

	for _, header := range headers {
		if err = w.Write(header); err != nil {
			log.Error(err)
			return
		}
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
