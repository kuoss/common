package tester

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaseName(t *testing.T) {
	testCases := []struct {
		name     string
		input    []any
		expected string
	}{
		{
			name:     "Single integer",
			input:    []any{1},
			expected: "01",
		},
		{
			name:     "Integer and string",
			input:    []any{1, "test"},
			expected: "01 test",
		},
		{
			name:     "Slice of strings",
			input:    []any{1, []string{"test", "value"}},
			expected: "01 test value",
		},
		{
			name:     "Map of strings",
			input:    []any{1, map[string]string{"key1": "value1", "key2": "value2"}},
			expected: "01 value1 value2",
		},
		{
			name: "Struct with fields",
			input: []any{1, struct {
				Field1 string
				Field2 int
				Field3 string
			}{"value1", 2, ""}},
			expected: "01 value1 2",
		},
		{
			name: "Different types",
			input: []any{1, "test", struct {
				Field1 string
				Field2 int
				Field3 string
			}{"value1", 2, ""}},
			expected: "01 test value1 2",
		},
		{
			name:     "Long string",
			input:    []any{1, strings.Repeat("a", 70)},
			expected: "01 aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa...",
		},
		{
			name:     "Special characters",
			input:    []any{1, "test", "/path"},
			expected: "01 test %path",
		},
		{
			name:     "Empty struct",
			input:    []any{1, struct{}{}},
			expected: "01",
		},
		{
			name:     "Nil pointer",
			input:    []any{1, (*string)(nil)},
			expected: "01",
		},
	}

	for i, tc := range testCases {
		t.Run(CaseName(i, tc.name), func(t *testing.T) {
			result := CaseName(tc.input...)
			assert.Equal(t, tc.expected, result, "Unexpected result for case: %s", tc.name)
		})
	}
}

func TestToString(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "Nil case",
			input:    nil,
			expected: "",
		},
		{
			name:     "Integer",
			input:    42,
			expected: "42",
		},
		{
			name:     "Float",
			input:    3.14,
			expected: "3.14",
		},
		{
			name:     "String",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "Slice of ints",
			input:    []int{1, 2, 3},
			expected: "1 2 3",
		},
		{
			name:     "Array of strings",
			input:    [3]string{"a", "b", "c"},
			expected: "a b c",
		},
		{
			name:     "Map",
			input:    map[string]int{"a": 1, "b": 2},
			expected: "1 2",
		},
		{
			name: "Struct with public fields",
			input: struct {
				Name  string
				Value int
			}{"Test", 10},
			expected: "Test 10",
		},
		{
			name: "Pointer to struct",
			input: &struct {
				Name  string
				Value int
			}{"Pointer", 20},
			expected: "Pointer 20",
		},
		{
			name: "Struct with unexported field",
			input: struct {
				Name   string
				Value  int
				hidden string
			}{"Visible", 30, "hidden"},
			expected: "Visible 30",
		},
		{
			name: "Struct with nil pointer field",
			input: struct {
				A int
				B *int
			}{A: 5, B: nil},
			expected: "5",
		},
	}

	for i, tc := range testCases {
		t.Run(CaseName(i, tc.name), func(t *testing.T) {
			result := toString(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
