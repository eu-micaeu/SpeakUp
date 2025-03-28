package handlers

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	openai "github.com/sashabaranov/go-openai"
)

// GenerateResponse generates a response for the AI
func GenerateResponseDialog(c *gin.Context) {

    prompt := `Your name is SpeakUp. You are a helpful assistant. You will be given a message and you need to respond to it in a friendly manner. Your response should be clear and concise, and you should avoid using overly technical language. If the message is a question, provide a direct answer. If it is a statement, acknowledge it and provide any additional information that may be helpful.`

    var request struct {
        Message string `json:"message"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    apiKey := os.Getenv("OPENAI_API_KEY")
    client := openai.NewClient(apiKey)

    // Generate a response for the dialogue
    dialogueResp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
        Model: "gpt-3.5-turbo",
        Messages: []openai.ChatCompletionMessage{
            {
                Role:    "user",
                Content: prompt + "Please pretend to be a friend and answer:" + request.Message,
            },
        },
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"response": dialogueResp.Choices[0].Message.Content})
}

// GenerateResponseCorrection generates a correction for the AI
func GenerateResponseCorrection(c *gin.Context) {
    var request struct {
        Message string `json:"message"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    apiKey := os.Getenv("OPENAI_API_KEY")
    client := openai.NewClient(apiKey)

    // Generate a correction for the dialogue
    correctionResp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
        Model: "gpt-3.5-turbo",
        Messages: []openai.ChatCompletionMessage{
            {
                Role:    "user",
                Content: "Please correct the spelling and grammar of the following text: " + request.Message,
            },
        },
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"response": correctionResp.Choices[0].Message.Content})
}