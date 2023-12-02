package main

import (
	"testing"
)

type stopWordTest struct {
	word  string
	there bool
}

func TestStop(t *testing.T) {

	tests := []stopWordTest{
		{"the", true},
		{"benvolio", false},
		{"diary", false},
		{"face", true},
	}

	stopWordsMap := stopWords()
	for _, val := range tests {
		if got, ok := stopWordsMap[val.word]; ok {

			if got != val.there {
				t.Errorf("Test failed Expected: %v, %v Got: %v", val.word, val.there, ok)
			}

		}
	}

}
