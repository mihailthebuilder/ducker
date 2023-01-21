package main

import (
	"testing"
)

func TestGetTextFromSite(t *testing.T) {
	url := "https://spec.openapis.org/oas/latest.html"
	out := getTextFromSite(url)
	writeStringToTestFile(out, "test.txt")
}
