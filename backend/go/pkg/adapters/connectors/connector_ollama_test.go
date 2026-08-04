package connectors

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestOllamaConnector_GenerateResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/generate" {
			t.Errorf("Expected path /api/generate, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		var req OllamaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if req.Prompt != "Hello Ollama" {
			t.Errorf("Expected prompt 'Hello Ollama', got '%s'", req.Prompt)
		}

		resp := OllamaResponse{
			Model:    req.Model,
			Response: "Hello back!",
			Done:     true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	os.Setenv("OLLAMA_HOST", ts.URL)
	os.Setenv("OLLAMA_MODEL", "llama3.2:latest")

	connector := NewOllamaConnector()
	res, err := connector.GenerateResponse(context.Background(), "Hello Ollama")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if res != "Hello back!" {
		t.Errorf("Expected 'Hello back!', got '%s'", res)
	}
}

func TestOllamaConnector_GenerateResponseWithOptions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req OllamaRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.System != "You are a tutor." {
			t.Errorf("Expected system prompt 'You are a tutor.', got '%s'", req.System)
		}

		if req.Options == nil {
			t.Errorf("Expected options to be non-nil")
		}

		resp := OllamaResponse{
			Model:    req.Model,
			Response: "Corrected output",
			Done:     true,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	os.Setenv("OLLAMA_HOST", ts.URL)
	os.Setenv("OLLAMA_MODEL", "llama3.2:latest")

	connector := NewOllamaConnector()
	options := map[string]any{
		"temperature": 0.1,
		"num_predict": 100,
	}
	res, err := connector.GenerateResponseWithOptions(context.Background(), "Test prompt", "You are a tutor.", options)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if res != "Corrected output" {
		t.Errorf("Expected 'Corrected output', got '%s'", res)
	}
}
