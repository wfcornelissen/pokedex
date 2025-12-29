package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	type testInput struct {
		test string
		want []string
	}

	inputText := []testInput{
		{
			test: "pikachu is the most popular pokemon",
			want: []string{"pikachu", "is", "the", "most", "popular", "pokemon"},
		},
		{
			test: "   this        one   has   multiple spaces   everywhere    ",
			want: []string{"this", "one", "has", "multiple", "spaces", "everywhere"},
		},
	}

	for _, test := range inputText {
		result := CleanInput(test.test)
		if !reflect.DeepEqual(result, test.want) {
			t.Errorf("\n%v\n Does not match to:\n%v", result, test.want)
			continue
		}
		fmt.Printf("\n%v\n equals\n%v\n", result, test.want)
	}

}
