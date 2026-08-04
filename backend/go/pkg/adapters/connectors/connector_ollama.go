package connectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type OllamaConnector struct {
	host   string
	model  string
	client *http.Client
}

type OllamaRequest struct {
	Model   string         `json:"model"`
	Prompt  string         `json:"prompt"`
	System  string         `json:"system,omitempty"`
	Stream  bool           `json:"stream"`
	Options map[string]any `json:"options,omitempty"`
}

type OllamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

func NewOllamaConnector() *OllamaConnector {
	host := os.Getenv("OLLAMA_HOST")
	if host == "" {
		host = "http://localhost:11434"
	}
	host = strings.TrimSuffix(host, "/")

	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "llama3.2:latest"
	}

	return &OllamaConnector{
		host:   host,
		model:  model,
		client: &http.Client{},
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
		System:  system,
		Stream:  false,
		Options: options,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request for Ollama: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", o.host)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request for Ollama: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request to Ollama at %s: %w", o.host, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading Ollama response: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("error unmarshaling Ollama response: %w", err)
	}

	return ollamaResp.Response, nil
}
