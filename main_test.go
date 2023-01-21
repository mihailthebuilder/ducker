package main

import (
	"os"
	"testing"
)

func TestGetTextFromSite(t *testing.T) {
	url := "https://spec.openapis.org/oas/latest.html"
	out := getTextFromSite(url)
	writeStringToTestFile(out)
}

func writeStringToTestFile(input string) {
	file, err := os.Create("test.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	_, err = file.WriteString(input)
	if err != nil {
		panic(err)
	}
}
