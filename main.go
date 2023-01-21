package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type TextCompletionsApiRequest struct {
	Model            string `json:"model"`
	Prompt           string `json:"prompt"`
	Temperature      int    `json:"temperature"`
	MaxTokens        int    `json:"max_tokens"`
	TopP             int    `json:"top_p"`
	FrequencePenalty int    `json:"frequence_penalty"`
	PresencePenalty  int    `json:"presence_penalty"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	openaiApiKey := os.Getenv("OPENAI_API_KEY")

	client := &http.Client{}

	url := "https://api.openai.com/v1/completions"

	prompt := createTextCompletionRequest("Write a terraform script that creates an AWS API Gateway that only allows requests from the IP 0.0.0.0")

	promptAsJson, err := json.Marshal(prompt)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("GET", url, bytes.NewReader(promptAsJson))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", openaiApiKey))

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}

func createTextCompletionRequest(prompt string) TextCompletionsApiRequest {
	return TextCompletionsApiRequest{
		Model:            "text-davinci-003",
		Prompt:           prompt,
		Temperature:      0,
		MaxTokens:        20,
		TopP:             1,
		FrequencePenalty: 0,
		PresencePenalty:  0,
	}
}
