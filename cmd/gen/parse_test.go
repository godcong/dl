package gen

import (
	"testing"
)

func TestExtractKeyValueWellFormed(t *testing.T) {
	typo := "[string]int"
	key, value := mapKeyValueTypes(typo)
	if key != "string" || value != "int" {
		t.Errorf("Expected key 'string' and value 'int', got key '%s' and value '%s'", key, value)
	}
}

func TestExtractKeyValueIncorrectOrder(t *testing.T) {
	typo := "]string[int["
	key, value := mapKeyValueTypes(typo)
	if key != "" || value != "" {
		t.Errorf("Expected empty key and value, got key '%s' and value '%s'", key, value)
	}
}

func TestExtractKeyValueNestedBrackets(t *testing.T) {
	typo := "[map[string]int]float64"
	key, value := mapKeyValueTypes(typo)
	if key != "map[string]int" || value != "float64" {
		t.Errorf("Expected key 'map[string]int' and value 'float64', got key '%s' and value '%s'", key, value)
	}
}

func TestMapKeyIndex(t *testing.T) {
	tests := []struct {
		value         string
		expectedStart int
		expectedEnd   int
	}{
		{value: "[string]int", expectedStart: 0, expectedEnd: 7},
		{value: "[int]string", expectedStart: 0, expectedEnd: 4},
		{value: "[*[]byte]*[]byte", expectedStart: 0, expectedEnd: 8},
		// Add more test cases as needed
	}

	for _, test := range tests {
		start, end := mapKeyIndex(test.value)
		if start != test.expectedStart || end != test.expectedEnd {
			t.Errorf("For value '%s', expected start: %d, end: %d, but got start: %d, end: %d", test.value, test.expectedStart, test.expectedEnd, start, end)
		}
	}
}
