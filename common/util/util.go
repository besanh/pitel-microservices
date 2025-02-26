package util

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
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
	if val < 0 {
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

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func GenerateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
