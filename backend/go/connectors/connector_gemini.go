package connectors

import (
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiConnector struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiConnector() *GeminiConnector {
	apiKey := os.Getenv("GEMINI_API_KEY")
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	model := client.GenerativeModel("gemini-2.0-flash")

	return &GeminiConnector{
		client: client,
		model:  model,
	}
}

func (g *GeminiConnector) GenerateResponse(ctx context.Context, message string) (string, error) {
	resp, err := g.model.GenerateContent(ctx, genai.Text(message))
	if err != nil {
		return "", err
	}

	if resp != nil && len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			return string(text), nil
		}
	}

	return "", nil
}

func (g *GeminiConnector) Close() {
	g.client.Close()
}
