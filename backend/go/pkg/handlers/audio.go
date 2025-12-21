package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TranscribeAudioRequest represents the request to transcribe audio
type TranscribeAudioRequest struct {
	Audio []byte `json:"audio"`
}

// TranscribeAudioResponse represents the response from audio transcription
type TranscribeAudioResponse struct {
	Text string `json:"text"`
}

// TranscribeAudio handles audio transcription (removed OpenAI dependency)
func TranscribeAudio(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Audio transcription not implemented without OpenAI"})
}
