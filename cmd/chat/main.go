package main

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	key := os.Getenv("ANTHROPIC_API_KEY")
	model := os.Getenv("MODEL")
	if model == "" {
		model = "claude-opus-5"
	}

	client := anthropic.NewClient(
		option.WithAPIKey(key),
	)

	messages := []anthropic.MessageParam{}

	system := []anthropic.TextBlockParam{}

	prompt := "Generate a very short event bridge rule as json."
	messages = addUserMessage(messages, prompt)

	response, err := chat(client, messages, model, system)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", response)
}

// chat sends the prompt to the API and returns the response.
func chat(
	client anthropic.Client,
	messages []anthropic.MessageParam,
	model string,
	system []anthropic.TextBlockParam,
) (string, error) {

	var temp param.Opt[float64]
	temp.Value = 0.0

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens:   1024,
		Messages:    messages,
		Model:       model,
		System:      system,
		Temperature: temp,
		OutputConfig: anthropic.OutputConfigParam{
			Format: anthropic.JSONOutputFormatParam{
				Schema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"source":      map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
						"detail-type": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
						"detail": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"state": map[string]any{
									"type":  "array",
									"items": map[string]any{"type": "string"},
								},
							},
							"required":             []string{"state"},
							"additionalProperties": false,
						},
					},
					"required":             []string{"source", "detail-type", "detail"},
					"additionalProperties": false,
				},
			},
		},
	},
	)
	if err != nil {
		return "", err
	}

	response := getResponse(message)
	return response, nil
}

// getResponse returns the response from the API to the console.
// Not needed for streaming, but useful for testing.
func getResponse(response *anthropic.Message) string {
	for _, block := range response.Content {
		if textBlock, ok := block.AsAny().(anthropic.TextBlock); ok {
			return textBlock.Text
		}
	}
	return ""
}

func addUserMessage(messages []anthropic.MessageParam, text string) []anthropic.MessageParam {
	messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(text)))
	return messages
}

func addAssistantMessage(messages []anthropic.MessageParam, text string) []anthropic.MessageParam {
	messages = append(messages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(text)))
	return messages
}
