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

// TestCreateChat tests the CreateChat handler
func TestCreateChat(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	gin.SetMode(gin.TestMode)

	t.Run("should create chat successfully", func(t *testing.T) {
		mt.Run("insert success", func(mt *mtest.T) {
			// Mock MongoDB insert response
			mt.AddMockResponses(mtest.CreateSuccessResponse())

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test chat
			chat := models.Chat{
				UserID:    "user123",
				Topic:     "Test Topic",
				StartTime: time.Now().Format(time.RFC3339),
			}

			// Marshal chat to JSON
			chatJSON, _ := json.Marshal(chat)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodPost, "/chat", bytes.NewBuffer(chatJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the CreateChat handler
			CreateChat(c)

			// Assert the response
			assert.Equal(t, http.StatusCreated, w.Code)
			var response models.Chat
			err := json.Unmarshal(w.Body.Bytes(), &response)

			assert.NoError(t, err)
			// assert.Equal(t, chat.UserID, response.UserID) TODO: Fix this
			assert.Equal(t, chat.Topic, response.Topic)
			assert.NotEmpty(t, response.ID)
		})
	})

	t.Run("should return error for invalid request body", func(t *testing.T) {
		mt.Run("invalid body", func(mt *mtest.T) {
			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create an invalid JSON request body
			invalidJSON := `{"user_id": "user123", "topic":}`

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodPost, "/chat", bytes.NewBufferString(invalidJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the CreateChat handler
			CreateChat(c)

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

			// Create a test chat
			chat := models.Chat{
				UserID:    "user123",
				Topic:     "Test Topic",
				StartTime: time.Now().Format(time.RFC3339),
			}

			// Marshal chat to JSON
			chatJSON, _ := json.Marshal(chat)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodPost, "/chat", bytes.NewBuffer(chatJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the CreateChat handler
			CreateChat(c)

			// Assert the response
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Failed to create chat", response["error"])
		})
	})
}

// TestGetChatById tests the GetChatById handler
func TestGetChatById(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	gin.SetMode(gin.TestMode)

	t.Run("should retrieve chat successfully", func(t *testing.T) {
		mt.Run("find success", func(mt *mtest.T) {
			// Mock MongoDB find response
			chat := models.Chat{
				ID:        "chat123",
				UserID:    "user123",
				Topic:     "Test Topic",
				StartTime: time.Now().Format(time.RFC3339),
			}
			var chatBson bson.D
			chatJSON, _ := json.Marshal(chat)
			_ = bson.UnmarshalExtJSON(chatJSON, true, &chatBson)
			mt.AddMockResponses(mtest.CreateCursorResponse(1, "speakup.chats", mtest.FirstBatch, chatBson))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodGet, "/chat/chat123", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: "chat123"}}

			// Call the GetChatById handler
			GetChatById(c)

			// Assert the response
			assert.Equal(t, http.StatusOK, w.Code)
			var response models.Chat
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, chat.ID, response.ID)
			assert.Equal(t, chat.Topic, response.Topic)
		})
	})

	t.Run("should return error if chat not found", func(t *testing.T) {
		mt.Run("find failure", func(mt *mtest.T) {
			// Mock MongoDB find response with no documents
			mt.AddMockResponses(mtest.CreateCursorResponse(0, "speakup.chats", mtest.FirstBatch))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodGet, "/chat/nonexistent", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: "nonexistent"}}

			// Call the GetChatById handler
			GetChatById(c)

			// Assert the response
			assert.Equal(t, http.StatusNotFound, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Chat not found", response["error"])
		})
	})
}

// TestGetChats tests the GetChats handler
func TestGetChats(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	gin.SetMode(gin.TestMode)

	t.Run("should retrieve all chats successfully", func(t *testing.T) {
		mt.Run("find success", func(mt *mtest.T) {
			// Mock MongoDB find response
			chats := []models.Chat{
				{
					ID:        "chat123",
					UserID:    "user123",
					Topic:     "Test Topic 1",
					StartTime: time.Now().Format(time.RFC3339),
				},
				{
					ID:        "chat456",
					UserID:    "user456",
					Topic:     "Test Topic 2",
					StartTime: time.Now().Format(time.RFC3339),
				},
			}
			var chatsBson []bson.D
			for _, chat := range chats {
				chatJSON, _ := json.Marshal(chat)
				var chatBson bson.D
				_ = bson.UnmarshalExtJSON(chatJSON, true, &chatBson)
				chatsBson = append(chatsBson, chatBson)
			}
			mt.AddMockResponses(mtest.CreateCursorResponse(0, "speakup.chats", mtest.FirstBatch, chatsBson...))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodGet, "/chat", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the GetChats handler
			GetChats(c)

			// Assert the response
			assert.Equal(t, http.StatusOK, w.Code)
			var response []models.Chat
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response, len(chats))
			assert.Equal(t, chats[0].ID, response[0].ID)
			assert.Equal(t, chats[1].ID, response[1].ID)
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
			req, _ := http.NewRequest(http.MethodGet, "/chat", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the GetChats handler
			GetChats(c)

			// Assert the response
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Failed to get chats", response["error"])
		})
	})
}

// TestDeleteChat tests the DeleteChat handler
func TestDeleteChat(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	gin.SetMode(gin.TestMode)

	t.Run("should delete chat successfully", func(t *testing.T) {
		mt.Run("delete success", func(mt *mtest.T) {
			// Mock MongoDB delete response
			mt.AddMockResponses(mtest.CreateSuccessResponse())

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodDelete, "/chat/chat123", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: "chat123"}}

			// Call the DeleteChat handler
			DeleteChat(c)

			// Assert the response
			assert.Equal(t, http.StatusOK, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Chat deleted successfully", response["message"])
		})
	})

	t.Run("should return error for database failure", func(t *testing.T) {
		mt.Run("delete failure", func(mt *mtest.T) {
			// Mock MongoDB delete failure response
			mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
				Code:    11000,
				Message: "delete error",
			}))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodDelete, "/chat/chat123", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: "chat123"}}

			// Call the DeleteChat handler
			DeleteChat(c)

			// Assert the response
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Failed to delete chat", response["error"])
		})
	})
}

// TestGetChatsByUserId tests the GetChatsByUserId handler
func TestGetChatsByUserId(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	gin.SetMode(gin.TestMode)

	t.Run("should retrieve user's chats successfully", func(t *testing.T) {
		mt.Run("find success", func(mt *mtest.T) {
			// Mock MongoDB find response
			chats := []models.Chat{
				{
					ID:        "chat123",
					UserID:    "user123",
					Topic:     "Test Topic 1",
					StartTime: time.Now().Format(time.RFC3339),
				},
				{
					ID:        "chat456",
					UserID:    "user123",
					Topic:     "Test Topic 2",
					StartTime: time.Now().Format(time.RFC3339),
				},
			}
			var chatsBson []bson.D
			for _, chat := range chats {
				chatJSON, _ := json.Marshal(chat)
				var chatBson bson.D
				_ = bson.UnmarshalExtJSON(chatJSON, true, &chatBson)
				chatsBson = append(chatsBson, chatBson)
			}
			mt.AddMockResponses(mtest.CreateCursorResponse(0, "speakup.chats", mtest.FirstBatch, chatsBson...))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodGet, "/chat/user", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Mock user ID in context
			c.Set("user_id", "user123")

			// Call the GetChatsByUserId handler
			GetChatsByUserId(c)

			// Assert the response
			assert.Equal(t, http.StatusOK, w.Code)
			var response struct {
				Chats []models.Chat `json:"chats"`
			}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response.Chats, len(chats))
			if len(response.Chats) > 0 {
				assert.Equal(t, chats[0].ID, response.Chats[0].ID)
				assert.Equal(t, chats[1].ID, response.Chats[1].ID)
			}
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
			req, _ := http.NewRequest(http.MethodGet, "/chat/user", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Mock user ID in context
			c.Set("user_id", "user123")

			// Call the GetChatsByUserId handler
			GetChatsByUserId(c)

			// Assert the response
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Failed to get chats", response["error"])
		})
	})
}