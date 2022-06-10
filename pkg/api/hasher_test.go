package api

import (
	"context"
	"os"
	"os/signal"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {
	ctx := context.Background()
	var results chan HashData
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	testTable := []struct {
		name        string
		inputValues HashData
		expected    []HashData
	}{
		{
			name: "exist data",
			inputValues: HashData{
				Hash:         "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				FileName:     "new",
				FullFilePath: "/local_path/gi/g/new",
				Algorithm:    "SHA256",
			},
			expected: []HashData{{
				Hash:         "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				FileName:     "new",
				FullFilePath: "/local_path/gi/g/new",
				Algorithm:    "SHA256",
			}}},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			go func(v HashData) {
				results <- v
			}(testCase.inputValues)

			got := Result(ctx, results, sig)

			assert.Equal(t, testCase.expected, got)
		})
	}
}
