package connectors

import (
	"context"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIConnector struct {
	client *openai.Client
	model  string
}

func NewOpenAIConnector() *OpenAIConnector {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)
	return &OpenAIConnector{
		client: client,
		model:  "gpt-4o",
	}
}

func (o *OpenAIConnector) GenerateResponse(ctx context.Context, message string) (string, error) {
	resp, err := o.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: o.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: message,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
