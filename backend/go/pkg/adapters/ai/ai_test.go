package ai

import (
	"context"
	"os"
	"testing"

	"speakup/pkg/adapters/connectors"
)

func TestGetConnector_Selection(t *testing.T) {
	// Test env variable AI_PROVIDER=ollama
	os.Setenv("AI_PROVIDER", "ollama")
	conn := getConnector(context.Background())
	if _, ok := conn.(*connectors.OllamaConnector); !ok {
		t.Errorf("Expected OllamaConnector when AI_PROVIDER=ollama")
	}

	// Test env variable AI_PROVIDER=gemini
	os.Setenv("AI_PROVIDER", "gemini")
	conn = getConnector(context.Background())
	if _, ok := conn.(*connectors.GeminiConnector); !ok {
		t.Errorf("Expected GeminiConnector when AI_PROVIDER=gemini")
	}

	// Test context override
	ctx := context.WithValue(context.Background(), "aiProvider", "ollama")
	conn = getConnector(ctx)
	if _, ok := conn.(*connectors.OllamaConnector); !ok {
		t.Errorf("Expected OllamaConnector when context aiProvider=ollama")
	}
}
