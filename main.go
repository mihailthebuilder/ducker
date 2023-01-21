package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/net/html"
)

type OpenAiClient struct {
	apiKey string
}

func main() {
	loadEnvironmentVariables()

	client := OpenAiClient{apiKey: os.Getenv("OPENAI_API_KEY")}

	// url := "https://spec.openapis.org/oas/latest.html"

	doc := getTextFromTestFile("input.txt")

	prompt := fmt.Sprintf(`Question: What data types are supported in the OpenAPI specification? Base your answer on the text below:\n"\n %s\n"`, doc)

	res := client.callTextCompletion(prompt)

	defer res.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(res)
	responseString := buf.String()

	fmt.Println(responseString)
}

func getTextFromTestFile(fileName string) string {
	data, err := os.ReadFile("test/" + fileName)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func getTextFromSite(url string) string {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic(fmt.Sprintf("get text from site error, status %s", res.Status))
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		panic(err)
	}

	var bodyText string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			var g func(*html.Node)
			g = func(n *html.Node) {
				if n.Type == html.TextNode {
					bodyText += n.Data
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					g(c)
				}
			}
			g(n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return bodyText
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

	writeStringToTestFile(string(promptAsJson), "test.json")

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
		Model:            "text-davinci-003",
		Prompt:           prompt,
		Temperature:      0,
		MaxTokens:        1000,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
}

func writeStringToTestFile(input string, fileName string) {
	file, err := os.Create("test/" + fileName)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	_, err = file.WriteString(input)
	if err != nil {
		panic(err)
	}
}
