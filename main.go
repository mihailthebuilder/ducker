package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type OpenAiClient struct {
	apiKey string
}

func main() {
	loadEnvironmentVariables()

	client := OpenAiClient{apiKey: os.Getenv("OPENAI_API_KEY")}

	url := "https://spec.openapis.org/oas/latest.html"

	doc := getTextFromSite(url)

	prompt := fmt.Sprintf("Question: What data types are supported in the OpenAPI specification? Base your answer on the text below:\n %s", doc)

	fmt.Println(prompt)

	res := client.callTextCompletion(prompt)

	buf := new(bytes.Buffer)
	buf.ReadFrom(res)
	responseString := buf.String()

	fmt.Println(responseString)
}

func getTextFromSite(url string) string {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != 200 {
		panic(fmt.Sprintf("get text from site error, status %s", res.Status))
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	responseString := buf.String()
	return responseString
}

func loadEnvironmentVariables() {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("error loading .env file: %s", err))
	}
}

func (o *OpenAiClient) callTextCompletion(prompt string) io.ReadCloser {

	client := &http.Client{}

	url := "https://api.openai.com/v1/completions"

	tcr := createTextCompletionRequest(prompt)

	promptAsJson, err := json.Marshal(tcr)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(promptAsJson))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.apiKey))

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != 200 {
		panic(fmt.Errorf("response status %s", res.Status))
	}

	return res.Body
}

type TextCompletionApiRequest struct {
	Model            string `json:"model"`
	Prompt           string `json:"prompt"`
	Temperature      int    `json:"temperature"`
	MaxTokens        int    `json:"max_tokens"`
	TopP             int    `json:"top_p"`
	FrequencyPenalty int    `json:"frequency_penalty"`
	PresencePenalty  int    `json:"presence_penalty"`
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
