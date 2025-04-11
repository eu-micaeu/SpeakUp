package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"speakup/config"
	"speakup/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// TestCreateUser tests the CreateUser handler
func TestCreateUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	gin.SetMode(gin.TestMode)

	t.Run("should create user successfully", func(t *testing.T) {
		mt.Run("insert success", func(mt *mtest.T) {
			// Mock MongoDB insert response
			mt.AddMockResponses(mtest.CreateSuccessResponse())

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test user
			user := models.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "password123",
				Language:  "en",
			}

			// Marshal user to JSON
			userJSON, _ := json.Marshal(user)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the CreateUser handler
			CreateUser(c)

			// Assert the response
			assert.Equal(t, http.StatusOK, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "User created successfully", response["message"])
		})
	})

	t.Run("should return error for database failure", func(t *testing.T) {
		mt.Run("insert failure", func(mt *mtest.T) {
			// Mock MongoDB insert failure response
			mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
				Index:   0,
				Code:    11000,
				Message: "duplicate key error",
			}))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test user
			user := models.User{
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane.doe@example.com",
				Password:  "password123",
				Language:  "en",
			}

			// Marshal user to JSON
			userJSON, _ := json.Marshal(user)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the CreateUser handler
			CreateUser(c)

			// Assert the response
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Failed to create user", response["error"])
		})
	})
}
