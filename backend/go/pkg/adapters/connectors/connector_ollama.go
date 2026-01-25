package connectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type OllamaConnector struct {
	baseURL string
	model   string
	client  *http.Client
}

type OllamaRequest struct {
	Model   string         `json:"model"`
	Prompt  string         `json:"prompt"`
	Stream  bool           `json:"stream"`
	System  string         `json:"system,omitempty"`
	Options map[string]any `json:"options,omitempty"`
}

type OllamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

func NewOllamaConnector() *OllamaConnector {
	baseURL := os.Getenv("OLLAMA_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "mistral"
	}

	return &OllamaConnector{
		baseURL: baseURL,
		model:   model,
		client:  &http.Client{},
	}
}

func (o *OllamaConnector) GenerateResponse(ctx context.Context, message string) (string, error) {
	return o.generate(ctx, message, "", nil)
}

func (o *OllamaConnector) GenerateResponseWithOptions(ctx context.Context, message string, system string, options map[string]any) (string, error) {
	return o.generate(ctx, message, system, options)
}

func (o *OllamaConnector) generate(ctx context.Context, message string, system string, options map[string]any) (string, error) {
	reqBody := OllamaRequest{
		Model:   o.model,
		Prompt:  message,
		Stream:  false,
		System:  system,
		Options: options,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.baseURL+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	return ollamaResp.Response, nil
}
