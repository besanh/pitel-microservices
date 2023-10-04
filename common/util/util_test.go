package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tel4vn/fins-microservices/common/log"
)

func TestParseInt(t *testing.T) {
	testCases := []struct {
		input  string
		output int
	}{
		{
			input:  "123",
			output: 123,
		},
		{
			input:  "abc",
			output: 0,
		},
	}
	for _, tc := range testCases {
		assert.Equal(t, tc.output, ParseInt(tc.input))
	}
}

func TestParseAnyToString(t *testing.T) {
	testCases := []struct {
		input   any
		output  string
		isError bool
	}{
		{
			input: map[string]any{
				"key": "123",
			},
			output:  "{\"key\":\"123\"}",
			isError: false,
		},
		{
			input:   make(chan int),
			output:  "",
			isError: true,
		},
	}
	for _, tc := range testCases {
		output, err := ParseAnyToString(tc.input)
		log.Infof(output)
		log.Error(err)
		if tc.isError {
			assert.Equal(t, true, err != nil)
		} else {
			assert.Equal(t, tc.output, output)
		}
	}
}

func TestParseLimit(t *testing.T) {
	testCases := []struct {
		input  any
		output int
	}{
		{
			input:  100,
			output: 100,
		},
		{
			input:  "100",
			output: 100,
		},
		{
			input:  "abc",
			output: 10,
		},
		{
			input:  -1,
			output: 10,
		},
		{
			input:  50001,
			output: 50000,
		},
	}
	for _, tc := range testCases {
		assert.Equal(t, tc.output, ParseLimit(tc.input))
	}
}

func TestParseOffset(t *testing.T) {
	testCases := []struct {
		input  any
		output int
	}{
		{
			input:  100,
			output: 100,
		},
		{
			input:  "100",
			output: 100,
		},
		{
			input:  "abc",
			output: 0,
		},
		{
			input:  -1,
			output: 0,
		},
	}
	for _, tc := range testCases {
		assert.Equal(t, tc.output, ParseOffset(tc.input))
	}
}

func TestParseStringToAny(t *testing.T) {
	testCases := []struct {
		input   string
		isError bool
	}{
		{
			input:   "abc",
			isError: true,
		}, {
			input:   "{\"key\":\"123\"}",
			isError: false,
		},
	}
	for _, tc := range testCases {
		var val any
		err := ParseStringToAny(tc.input, &val)
		assert.Equal(t, tc.isError, err != nil)
	}
}

func TestParseAnyToAny(t *testing.T) {
	testCases := []struct {
		input   any
		isError bool
	}{
		{
			input:   make(chan int),
			isError: true,
		},
		{
			input:   "{\"key\":\"123\"}",
			isError: false,
		},
		{
			input:   "abc",
			isError: false,
		},
	}
	for _, tc := range testCases {
		if tc.isError {
			var val int
			err := ParseAnyToAny(tc.input, &val)
			assert.Equal(t, true, err != nil)
		} else {
			var val any
			err := ParseAnyToAny(tc.input, &val)
			assert.Equal(t, false, err != nil)
		}
	}
}

func TestParseString(t *testing.T) {
	testCases := []struct {
		input  any
		output string
	}{
		{
			input:  "abc",
			output: "abc",
		},
		{
			input:  123,
			output: "",
		},
	}
	for _, tc := range testCases {
		assert.Equal(t, tc.output, ParseString(tc.input))
	}
}

func TestInArray(t *testing.T) {
	testCases := []struct {
		arr    any
		item   any
		output bool
	}{
		{
			arr:    "abc",
			output: false,
		},
		{
			arr:    []string{"abc"},
			item:   "abc",
			output: true,
		},
		{
			arr:    []string{"abc"},
			item:   "abcd",
			output: false,
		},
	}
	for _, tc := range testCases {
		assert.Equal(t, tc.output, InArray(tc.item, tc.arr))
	}
}

func TestGenerateRandomString(t *testing.T) {
	assert.Equal(t, true, len(GenerateRandomString(3)) == 3)
}
