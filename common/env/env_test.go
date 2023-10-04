package env

import (
	"testing"

	"gotest.tools/assert"
)

func TestGetStringENV(t *testing.T) {
	testCases := []struct {
		setVar       string
		envVar       string
		defaultValue string
		output       string
	}{
		{
			setVar:       "TEST",
			envVar:       "TEST",
			defaultValue: "",
			output:       "TEST",
		},
		{
			setVar:       "",
			envVar:       "TEST",
			defaultValue: "abc",
			output:       "abc",
		},
	}
	for _, tc := range testCases {
		t.Setenv("TEST", tc.setVar)
		assert.Equal(t, tc.output, GetStringENV(tc.envVar, tc.defaultValue))
	}
}

func TestGetIntENV(t *testing.T) {
	testCases := []struct {
		setVar       string
		envVar       string
		defaultValue int
		output       int
	}{
		{
			setVar:       "80",
			envVar:       "TEST",
			defaultValue: 0,
			output:       80,
		},
		{
			setVar:       "",
			envVar:       "TEST",
			defaultValue: 100,
			output:       100,
		},
	}
	for _, tc := range testCases {
		t.Setenv("TEST", tc.setVar)
		assert.Equal(t, tc.output, GetIntENV(tc.envVar, tc.defaultValue))
	}
}

func TestGetBoolENV(t *testing.T) {
	testCases := []struct {
		setVar       string
		envVar       string
		defaultValue bool
		output       bool
	}{
		{
			setVar:       "true",
			envVar:       "TEST",
			defaultValue: false,
			output:       true,
		},
		{
			setVar:       "",
			envVar:       "TEST",
			defaultValue: false,
			output:       false,
		},
	}
	for _, tc := range testCases {
		t.Setenv("TEST", tc.setVar)
		assert.Equal(t, tc.output, GetBoolENV(tc.envVar, tc.defaultValue))
	}
}
