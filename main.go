package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	openaiApiKey := os.Getenv("OPENAI_API_KEY")

	client := &http.Client{}

	url := "https://api.openai.com/v1/completions"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", openaiApiKey))

	res, _ := client.Do(req)
	fmt.Println(res)
}
