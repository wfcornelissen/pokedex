package main

import (
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "Hello, World!",
			expected: []string{"hello,", "world!"},
		},
		{
			input:    "This is a test",
			expected: []string{"this", "is", "a", "test"},
		},
		{
			input:    "12345",
			expected: []string{"12345"},
		},
		{
			input:    "Hello, World! This is a test",
			expected: []string{"hello,", "world!", "this", "is", "a", "test"},
		},
	}

	for _, c := range cases {
		result := cleanInput(c.input)
		if !reflect.DeepEqual(result, c.expected) {
			t.Errorf("Expected %v, got %v", c.expected, result)
		}
	}
}
