package main

import (
	"os"
	"testing"
)

func TestGetTextFromSite(t *testing.T) {
	url := "https://spec.openapis.org/oas/latest.html"
	out := getTextFromSite(url)

	file, err := os.Create("text.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	_, err = file.WriteString(out)
	if err != nil {
		panic(err)
	}
}
