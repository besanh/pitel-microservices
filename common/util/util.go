package util

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
)

const MAX_LIMIT = 50_000

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
	if val <= 0 {
		val = 10
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

func ParseId(id any) (result string, ok bool) {
	if reflect.TypeOf(id).Kind() == reflect.String {
		result = id.(string)
	} else if reflect.TypeOf(id).Kind() == reflect.Int {
		result = fmt.Sprintf("%d", id)
	}
	result = strings.TrimSpace(result)
	return result, result != ""
}

func ParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func ParseAnyToString(value any) (string, error) {
	ref := reflect.ValueOf(value)
	if ref.Kind() == reflect.String {
		return value.(string), nil
	} else if InArray(ref.Kind(), []reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64}) {
		return fmt.Sprintf("%d", value), nil
	} else if InArray(ref.Kind(), []reflect.Kind{reflect.Float32, reflect.Float64}) {
		return fmt.Sprintf("%f", value), nil
	} else if ref.Kind() == reflect.Bool {
		return fmt.Sprintf("%t", value), nil
	} else if ref.Kind() == reflect.Slice {
		return fmt.Sprintf("%v", value), nil
	}
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func MustParseAnyToString(value any) string {
	str, err := ParseAnyToString(value)
	if err != nil {
		return ""
	}
	return str
}

func ParseStringToAny(value string, dest any) error {
	if err := json.Unmarshal([]byte(value), dest); err != nil {
		return err
	}
	return nil
}

func ParseAnyToAny(value any, dest any) (err error) {
	ref := reflect.ValueOf(value)
	var bytes []byte
	if ref.Kind() == reflect.String {
		bytes = []byte(value.(string))
	} else {
		bytes, err = json.Marshal(value)
		if err != nil {
			return err
		}
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

func InArray[T comparable](element T, slice []T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

var LETTER_RUNES = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

var NUMBER_RUNES = []rune("1234567890")

func GenerateRandomString(n int, letterRunes []rune) string {
	if letterRunes == nil || len(letterRunes) < 1 {
		letterRunes = LETTER_RUNES
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func InArrayContains(item string, array []string) bool {
	for _, v := range array {
		if strings.Contains(item, v) {
			return true
		}
	}
	return false
}

func ParseStringToTime(t string) *time.Time {
	if len(t) == 0 {
		return nil
	}
	c := carbon.Parse(t, "Asia/Ho_Chi_Minh")
	if c.Error != nil {
		return nil
	}
	tPtr := c.StdTime()
	return &tPtr
}

func GetCurrentPtrTime() *time.Time {
	t := time.Now()
	return &t
}

func CheckArrayContainValues(array []string, values []string) bool {
	for _, value := range values {
		if !InArrayContains(value, array) {
			return false
		}
	}
	return true
}

func ParseRecordToMap(headers, record []string) map[string]string {
	result := make(map[string]string)
	for i := 0; i < len(headers); i++ {
		result[headers[i]] = record[i]
	}
	return result
}

func ProcessSliceString(slice []string) []string {
	result := make([]string, 0)
	for _, v := range slice {
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}

func CheckFromAndToDateValid(from, to time.Time, isAllowZero bool) (isOk bool, err error) {
	if !isAllowZero {
		if from.IsZero() {
			return false, errors.New("is zero")
		}
		if to.IsZero() {
			return false, errors.New("is zero")
		}
	} else if from.After(to) {
		return false, errors.New("from date must be before to date")
	}
	isOk = true
	return
}

func ParseToAnyArray(value []string) []any {
	result := make([]any, 0)
	for _, v := range value {
		result = append(result, v)
	}
	return result
}

func RemoveElementFromSlice(slice []string, elementToRemove string) []string {
	var result []string
	for _, value := range slice {
		if value != elementToRemove {
			result = append(result, value)
		}
	}
	return result
}

func UnionSlices(first []string, second []string) []string {
	out := []string{}
	out = append(out, first...)
	out = append(out, second...)

	return out
}

func IntersectSlices(first []string, second []string) []string {
	out := []string{}
	bucket := map[string]bool{}

	for _, i := range first {
		for _, j := range second {
			if i == j && !bucket[i] {
				out = append(out, i)
				bucket[i] = true
			}
		}
	}

	return out
}

func RemoveDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func UnionWithoutDuplicate[T comparable](sliceList ...[]T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, slice := range sliceList {
		for _, item := range slice {
			if _, value := allKeys[item]; !value {
				allKeys[item] = true
				list = append(list, item)
			}
		}
	}
	return list
}

func ParseTimeFromString(timeString string) (time.Time, error) {
	c := carbon.Parse(timeString, "Asia/Ho_Chi_Minh")
	if c.IsInvalid() {
		return time.Time{}, errors.New("invalid time")
	} else if c.Error != nil {
		return time.Time{}, c.Error
	}
	return c.StdTime(), nil
}

func MustParseTimeFromString(timeString string) time.Time {
	result, err := ParseTimeFromString(timeString)
	if err != nil {
		return time.Time{}
	}
	return result
}

func ParseFloat64(value any) float64 {
	if value == nil {
		return 0
	}
	// convert to string
	str := MustParseAnyToString(value)
	// convert to float
	result, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return result
}

func IsUuid(value string) bool {
	_, err := uuid.Parse(value)
	return err == nil
}

func GetSubStringFromEnd(value string, length int) string {
	if len(value) < length {
		return value
	}
	return value[len(value)-length:]
}

func Base64Encode(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func Base64Decode(value string) (string, error) {
	result, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
