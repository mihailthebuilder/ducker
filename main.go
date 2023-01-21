package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type TextCompletionApiRequest struct {
	Model            string `json:"model"`
	Prompt           string `json:"prompt"`
	Temperature      int    `json:"temperature"`
	MaxTokens        int    `json:"max_tokens"`
	TopP             int    `json:"top_p"`
	FrequencyPenalty int    `json:"frequency_penalty"`
	PresencePenalty  int    `json:"presence_penalty"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("error loading .env file: %s", err))
	}

	openaiApiKey := os.Getenv("OPENAI_API_KEY")

	client := &http.Client{}

	url := "https://api.openai.com/v1/completions"

	prompt := createTextCompletionRequest("Write a terraform script that creates an AWS API Gateway that only allows requests from the IP 0.0.0.0")

	promptAsJson, err := json.Marshal(prompt)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(promptAsJson))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", openaiApiKey))

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != 200 {
		panic(fmt.Errorf("response status %s", res.Status))
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	responseString := buf.String()

	fmt.Println(responseString)
}

func createTextCompletionRequest(prompt string) TextCompletionApiRequest {
	return TextCompletionApiRequest{
		Model:            "text-curie-001",
		Prompt:           prompt,
		Temperature:      0,
		MaxTokens:        20,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
}
