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

type GeminiConnector struct {
	apiKey string
	model  string
	client *http.Client
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiSystemInstruction struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiThinkingConfig struct {
	ThinkingBudget int `json:"thinkingBudget"`
}

type GeminiGenerationConfig struct {
	Temperature     float64               `json:"temperature,omitempty"`
	TopP            float64               `json:"topP,omitempty"`
	MaxOutputTokens int                   `json:"maxOutputTokens,omitempty"`
	StopSequences   []string              `json:"stopSequences,omitempty"`
	ThinkingConfig  *GeminiThinkingConfig `json:"thinkingConfig,omitempty"`
}

type GeminiRequest struct {
	Contents          []GeminiContent          `json:"contents"`
	SystemInstruction *GeminiSystemInstruction `json:"systemInstruction,omitempty"`
	GenerationConfig  *GeminiGenerationConfig  `json:"generationConfig,omitempty"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []GeminiPart `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func NewGeminiConnector() *GeminiConnector {
	apiKey := os.Getenv("GEMINI_API_KEY")
	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = "gemini-1.5-flash"
	}

	return &GeminiConnector{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{},
	}
}

func (g *GeminiConnector) GenerateResponse(ctx context.Context, message string) (string, error) {
	return g.generate(ctx, message, "", nil)
}

func (g *GeminiConnector) GenerateResponseWithOptions(ctx context.Context, message string, system string, options map[string]any) (string, error) {
	return g.generate(ctx, message, system, options)
}

func (g *GeminiConnector) generate(ctx context.Context, message string, system string, options map[string]any) (string, error) {
	if g.apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY is not set in environment variables")
	}

	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: message},
				},
			},
		},
	}

	if system != "" {
		reqBody.SystemInstruction = &GeminiSystemInstruction{
			Parts: []GeminiPart{
				{Text: system},
			},
		}
	}

	if options != nil {
		config := &GeminiGenerationConfig{}
		hasConfig := false

		if temp, ok := options["temperature"]; ok {
			if t, ok := temp.(float64); ok {
				config.Temperature = t
				hasConfig = true
			} else if t, ok := temp.(float32); ok {
				config.Temperature = float64(t)
				hasConfig = true
			} else if t, ok := temp.(int); ok {
				config.Temperature = float64(t)
				hasConfig = true
			}
		}

		if topP, ok := options["top_p"]; ok {
			if p, ok := topP.(float64); ok {
				config.TopP = p
				hasConfig = true
			} else if p, ok := topP.(float32); ok {
				config.TopP = float64(p)
				hasConfig = true
			}
		}

		if maxTokens, ok := options["num_predict"]; ok {
			var limit int
			if m, ok := maxTokens.(int); ok {
				limit = m
			} else if m, ok := maxTokens.(float64); ok {
				limit = int(m)
			}
			if limit > 0 {
				config.MaxOutputTokens = limit
				hasConfig = true
				if limit <= 256 {
					config.ThinkingConfig = &GeminiThinkingConfig{ThinkingBudget: 0}
				}
			}
		}

		if stop, ok := options["stop"]; ok {
			if sList, ok := stop.([]string); ok {
				config.StopSequences = sList
				hasConfig = true
			} else if sList, ok := stop.([]any); ok {
				for _, s := range sList {
					if str, ok := s.(string); ok {
						config.StopSequences = append(config.StopSequences, str)
					}
				}
				if len(config.StopSequences) > 0 {
					hasConfig = true
				}
			}
		}

		if hasConfig {
			reqBody.GenerationConfig = config
		}
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", g.model, g.apiKey)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request to Gemini: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gemini API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini returned an empty response")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}
