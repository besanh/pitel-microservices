package util

import "regexp"

func CheckPatternNumeric(Key string) bool {
	// pattern := "^[1]\\d{9}|[9]\\d{8}|[2-8]\\d{8,9}$"
	r, _ := regexp.Compile("^\\d+$")
	return r.MatchString(Key)
}
