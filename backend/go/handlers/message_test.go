package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"speakup/config"
	"speakup/models"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// TestCreateMessage tests the CreateMessage handler
func TestCreateMessage(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	gin.SetMode(gin.TestMode)

	t.Run("should create message successfully", func(t *testing.T) {
		mt.Run("insert success", func(mt *mtest.T) {
			// Mock MongoDB insert response
			mt.AddMockResponses(mtest.CreateSuccessResponse())

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test message
			message := models.Message{
				ChatID:    "chat123",
				Content:   "Hello, world!",
				Sender:    "user123",
				Type:      "text",
				CreatedAt: time.Now().Format(time.RFC3339),
			}

			// Marshal message to JSON
			messageJSON, _ := json.Marshal(message)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodPost, "/message", bytes.NewBuffer(messageJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the CreateMessage handler
			CreateMessage(c)

			// Assert the response
			assert.Equal(t, http.StatusCreated, w.Code)
			var response models.Message
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, message.ChatID, response.ChatID)
			assert.Equal(t, message.Content, response.Content)
			assert.Equal(t, message.Sender, response.Sender)
			assert.Equal(t, message.Type, response.Type)
			assert.NotEmpty(t, response.ID)
		})
	})

	t.Run("should return error for invalid request body", func(t *testing.T) {
		mt.Run("invalid body", func(mt *mtest.T) {
			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create an invalid JSON request body
			invalidJSON := `{"chat_id": "chat123", "content":}`

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodPost, "/message", bytes.NewBufferString(invalidJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the CreateMessage handler
			CreateMessage(c)

			// Assert the response
			assert.Equal(t, http.StatusBadRequest, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Contains(t, response["error"], "invalid character")
		})
	})

	t.Run("should return error for database failure", func(t *testing.T) {
		mt.Run("insert failure", func(mt *mtest.T) {
			// Mock MongoDB insert failure response
			mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
				Index:   0,
				Code:    11000,
				Message: "insert error",
			}))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test message
			message := models.Message{
				ChatID:    "chat123",
				Content:   "Hello, world!",
				Sender:    "user123",
				Type:      "text",
				CreatedAt: time.Now().Format(time.RFC3339),
			}

			// Marshal message to JSON
			messageJSON, _ := json.Marshal(message)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodPost, "/message", bytes.NewBuffer(messageJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the CreateMessage handler
			CreateMessage(c)

			// Assert the response
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Failed to create message", response["error"])
		})
	})
}

// TestGetMessageById tests the GetMessageById handler
func TestGetMessageById(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	gin.SetMode(gin.TestMode)

	t.Run("should retrieve message successfully", func(t *testing.T) {
		mt.Run("find success", func(mt *mtest.T) {
			// Mock MongoDB find response
			message := models.Message{
				ID:        "message123",
				ChatID:    "chat123",
				Content:   "Hello, world!",
				Sender:    "user123",
				Type:      "text",
				CreatedAt: time.Now().Format(time.RFC3339),
			}
			mt.AddMockResponses(mtest.CreateCursorResponse(1, "speakup.messages", mtest.FirstBatch, bson.D{
				{Key: "id", Value: message.ID},
				{Key: "chat_id", Value: message.ChatID},
				{Key: "content", Value: message.Content},
				{Key: "sender", Value: message.Sender},
				{Key: "type", Value: message.Type},
				{Key: "created_at", Value: message.CreatedAt},
			}))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodGet, "/message/message123", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "message123"}}

			// Call the GetMessageById handler
			GetMessageById(c)

			// Assert the response
			assert.Equal(t, http.StatusOK, w.Code)
			var response models.Message
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, message.ID, response.ID)
			assert.Equal(t, message.ChatID, response.ChatID)
			assert.Equal(t, message.Content, response.Content)
			assert.Equal(t, message.Sender, response.Sender)
			assert.Equal(t, message.Type, response.Type)
		})
	})

	t.Run("should return error for non-existent message", func(t *testing.T) {
		mt.Run("find failure", func(mt *mtest.T) {
			// Mock MongoDB find failure response
			mt.AddMockResponses(mtest.CreateCursorResponse(0, "speakup.messages", mtest.FirstBatch))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodGet, "/message/nonexistent", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}

			// Call the GetMessageById handler
			GetMessageById(c)

			// Assert the response
			assert.Equal(t, http.StatusNotFound, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Message not found", response["error"])
		})
	})
}

// TestGetMessages tests the GetMessages handler
func TestGetMessages(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	gin.SetMode(gin.TestMode)

	t.Run("should retrieve all messages successfully", func(t *testing.T) {
		mt.Run("find success", func(mt *mtest.T) {
			// Mock MongoDB find response
			messages := []bson.D{
				{
					{Key: "id", Value: "message1"},
					{Key: "chat_id", Value: "chat123"},
					{Key: "content", Value: "Hello, world!"},
					{Key: "sender", Value: "user123"},
					{Key: "type", Value: "text"},
					{Key: "created_at", Value: time.Now().Format(time.RFC3339)},
				},
				{
					{Key: "id", Value: "message2"},
					{Key: "chat_id", Value: "chat123"},
					{Key: "content", Value: "Hi again!"},
					{Key: "sender", Value: "user456"},
					{Key: "type", Value: "text"},
					{Key: "created_at", Value: time.Now().Format(time.RFC3339)},
				},
			}
			mt.AddMockResponses(mtest.CreateCursorResponse(0, "speakup.messages", mtest.FirstBatch, messages...))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodGet, "/message", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the GetMessages handler
			GetMessages(c)

			// Assert the response
			assert.Equal(t, http.StatusOK, w.Code)
			var response []models.Message
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, 2)
			assert.Equal(t, "message1", response[0].ID)
			assert.Equal(t, "message2", response[1].ID)
		})
	})

	t.Run("should return error for database failure", func(t *testing.T) {
		mt.Run("find failure", func(mt *mtest.T) {
			// Mock MongoDB find failure response
			mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
				Code:    11000,
				Message: "find error",
			}))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodGet, "/message", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the GetMessages handler
			GetMessages(c)

			// Assert the response
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Failed to get messages", response["error"])
		})
	})
}

// TestUpdateMessage tests the UpdateMessage handler
func TestUpdateMessage(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	gin.SetMode(gin.TestMode)

	t.Run("should update message successfully", func(t *testing.T) {
		mt.Run("update success", func(mt *mtest.T) {
			// Mock MongoDB update response
			mt.AddMockResponses(mtest.CreateSuccessResponse())

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test message update
			update := models.Message{
				Content: "Updated content",
				Type:    "text",
			}

			// Marshal update to JSON
			updateJSON, _ := json.Marshal(update)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodPut, "/message/message123", bytes.NewBuffer(updateJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "message123"}}

			// Call the UpdateMessage handler
			UpdateMessage(c)

			// Assert the response
			assert.Equal(t, http.StatusOK, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Message updated successfully", response["message"])
		})
	})

	t.Run("should return error for invalid request body", func(t *testing.T) {
		mt.Run("invalid body", func(mt *mtest.T) {
			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create an invalid JSON request body
			invalidJSON := `{"content":}`

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodPut, "/message/message123", bytes.NewBufferString(invalidJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "message123"}}

			// Call the UpdateMessage handler
			UpdateMessage(c)

			// Assert the response
			assert.Equal(t, http.StatusBadRequest, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Contains(t, response["error"], "invalid character")
		})
	})

	t.Run("should return error for database failure", func(t *testing.T) {
		mt.Run("update failure", func(mt *mtest.T) {
			// Mock MongoDB update failure response
			mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
				Index:   0,
				Code:    11000,
				Message: "update error",
			}))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test message update
			update := models.Message{
				Content: "Updated content",
				Type:    "text",
			}

			// Marshal update to JSON
			updateJSON, _ := json.Marshal(update)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodPut, "/message/message123", bytes.NewBuffer(updateJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "message123"}}

			// Call the UpdateMessage handler
			UpdateMessage(c)

			// Assert the response
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Failed to update message", response["error"])
		})
	})
}

// TestDeleteMessage tests the DeleteMessage handler
func TestDeleteMessage(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	gin.SetMode(gin.TestMode)

	t.Run("should delete message successfully", func(t *testing.T) {
		mt.Run("delete success", func(mt *mtest.T) {
			// Mock MongoDB delete response
			mt.AddMockResponses(mtest.CreateSuccessResponse())

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodDelete, "/message/message123", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "message123"}}

			// Call the DeleteMessage handler
			DeleteMessage(c)

			// Assert the response
			assert.Equal(t, http.StatusOK, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Message deleted successfully", response["message"])
		})
	})

	t.Run("should return error for non-existent message", func(t *testing.T) {
		mt.Run("delete failure", func(mt *mtest.T) {
			// Mock MongoDB delete failure response
			mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
				Index:   0,
				Code:    11000,
				Message: "delete error",
			}))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodDelete, "/message/nonexistent", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}

			// Call the DeleteMessage handler
			DeleteMessage(c)

			// Assert the response
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Failed to delete message", response["error"])
		})
	})
}

// TestGetMessagesByChatId tests the GetMessagesByChatId handler
func TestGetMessagesByChatId(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	gin.SetMode(gin.TestMode)

	t.Run("should retrieve messages by chat ID successfully", func(t *testing.T) {
		mt.Run("find success", func(mt *mtest.T) {
			// Mock MongoDB find response
			messages := []bson.D{
				{
					{Key: "id", Value: "message1"},
					{Key: "chat_id", Value: "chat123"},
					{Key: "content", Value: "Hello, world!"},
					{Key: "sender", Value: "user123"},
					{Key: "type", Value: "text"},
					{Key: "created_at", Value: time.Now().Format(time.RFC3339)},
				},
				{
					{Key: "id", Value: "message2"},
					{Key: "chat_id", Value: "chat123"},
					{Key: "content", Value: "Hi again!"},
					{Key: "sender", Value: "user456"},
					{Key: "type", Value: "text"},
					{Key: "created_at", Value: time.Now().Format(time.RFC3339)},
				},
			}
			mt.AddMockResponses(mtest.CreateCursorResponse(0, "speakup.messages", mtest.FirstBatch, messages...))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodGet, "/message/chat/chat123", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "chat123"}}

			// Call the GetMessagesByChatId handler
			GetMessagesByChatId(c)

			// Assert the response
			assert.Equal(t, http.StatusOK, w.Code)
			var response []models.Message
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, 2)
			assert.Equal(t, "message1", response[0].ID)
			assert.Equal(t, "message2", response[1].ID)
		})
	})

	t.Run("should return error for database failure", func(t *testing.T) {
		mt.Run("find failure", func(mt *mtest.T) {
			// Mock MongoDB find failure response
			mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
				Code:    11000,
				Message: "find error",
			}))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodGet, "/message/chat/chat123", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "chat123"}}

			// Call the GetMessagesByChatId handler
			GetMessagesByChatId(c)

			// Assert the response
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Failed to get messages", response["error"])
		})
	})
}