package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "    ",
			expected: []string{},
		},
		{
			input:    "Hello, World!",
			expected: []string{"hello,", "world!"},
		},
		{
			input:    "This is a test Sentence",
			expected: []string{"this", "is", "a", "test", "sentence"},
		},
		{
			input:    "    hasleadingspace",
			expected: []string{"hasleadingspace"},
		},
		{
			input:    "hastrailingspace     ",
			expected: []string{"hastrailingspace"},
		},
		{
			input:    "      hasbothspace      ",
			expected: []string{"hasbothspace"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("lengths don't match: '%v' vs '%v'", actual, c.expected)
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("%v does not match with %v", word, expectedWord)
				t.Fail()
			}
		}

	}
}
