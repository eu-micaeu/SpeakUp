package handlers

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
)

// TranscribeAudioRequest represents the request to transcribe audio
type TranscribeAudioRequest struct {
	Audio []byte `json:"audio"`
}

// TranscribeAudioResponse represents the response from audio transcription
type TranscribeAudioResponse struct {
	Text string `json:"text"`
}

// TranscribeAudio handles audio transcription using OpenAI Whisper API
func TranscribeAudio(c *gin.Context) {
	// Parse multipart form
	file, header, err := c.Request.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get audio file: " + err.Error()})
		return
	}
	defer file.Close()

	// Create temporary file
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "audio_"+header.Filename)
	
	out, err := os.Create(tempFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temp file: " + err.Error()})
		return
	}
	defer out.Close()
	defer os.Remove(tempFile)

	// Copy uploaded file to temp file
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save audio file: " + err.Error()})
		return
	}

	// Reopen file for reading
	audioFile, err := os.Open(tempFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open audio file: " + err.Error()})
		return
	}
	defer audioFile.Close()

	// Initialize OpenAI client
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OpenAI API key not configured"})
		return
	}

	client := openai.NewClient(apiKey)

	// Create transcription request
	ctx := context.Background()
	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: tempFile,
		Format:   openai.AudioResponseFormatJSON,
		Language: "en",
	}

	resp, err := client.CreateTranscription(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to transcribe audio: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, TranscribeAudioResponse{
		Text: resp.Text,
	})
}
