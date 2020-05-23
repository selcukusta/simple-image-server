package helper

import "testing"

type TestDataItem struct {
	input    string
	minValue int
	maxValue int
	expected bool
}

func TestValidateRange(t *testing.T) {
	dataItems := []TestDataItem{
		{"100", 0, 100, true},
		{"0", 0, 100, true},
		{"101", 0, 100, false},
		{"-1", 0, 100, false},
		{"", 0, 100, false},
	}

	for _, item := range dataItems {
		result := ValidateRange(item.input, item.minValue, item.maxValue)
		if result != item.expected {
			t.Errorf(`ValidateRange("%s", %d, %d) is FAILED. Expected "%t" but got value "%t"`, item.input, item.minValue, item.maxValue, item.expected, result)
		}
	}
}
