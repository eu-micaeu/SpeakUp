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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"golang.org/x/crypto/bcrypt"
)

// TestLogin tests the Login handler
func TestLogin(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	gin.SetMode(gin.TestMode)

	t.Run("should login successfully", func(t *testing.T) {
		mt.Run("find success", func(mt *mtest.T) {
			// Mock MongoDB find response
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
			user := models.User{
				ID:       "12345",
				Email:    "john.doe@example.com",
				Password: string(hashedPassword),
				Language: "en",
			}
			var mockResponse bson.D
			bsonBytes, _ := bson.Marshal(user)
			_ = bson.Unmarshal(bsonBytes, &mockResponse)
			mt.AddMockResponses(mtest.CreateCursorResponse(1, "speakup.users", mtest.FirstBatch, mockResponse))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create login credentials
			credentials := map[string]string{
				"email":    "john.doe@example.com",
				"password": "password123",
			}
			credentialsJSON, _ := json.Marshal(credentials)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer(credentialsJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the Login handler
			Login(c)

			// Assert the response
			assert.Equal(t, http.StatusOK, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, "Login successful", response["message"])
			assert.NotEmpty(t, response["token"])
		})

	})

	t.Run("should return error for invalid email", func(t *testing.T) {
		mt.Run("find failure", func(mt *mtest.T) {
			// Mock MongoDB find failure response
			mt.AddMockResponses(mtest.CreateCursorResponse(0, "speakup.users", mtest.FirstBatch))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create login credentials
			credentials := map[string]string{
				"email":    "invalid.email@example.com",
				"password": "password123",
			}
			credentialsJSON, _ := json.Marshal(credentials)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer(credentialsJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the Login handler
			Login(c)

			// Assert the response
			assert.Equal(t, http.StatusUnauthorized, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Invalid email", response["error"])
		})
	})

	t.Run("should return error for invalid password", func(t *testing.T) {
		mt.Run("find success but invalid password", func(mt *mtest.T) {
			// Mock MongoDB find response
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
			user := models.User{
				ID:       "12345",
				Email:    "john.doe@example.com",
				Password: string(hashedPassword),
				Language: "en",
			}
			var mockResponse bson.D
			bsonBytes, _ := bson.Marshal(user)
			_ = bson.Unmarshal(bsonBytes, &mockResponse)
			mt.AddMockResponses(mtest.CreateCursorResponse(1, "speakup.users", mtest.FirstBatch, mockResponse))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create login credentials
			credentials := map[string]string{
				"email":    "john.doe@example.com",
				"password": "wrongpassword",
			}
			credentialsJSON, _ := json.Marshal(credentials)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer(credentialsJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the Login handler
			Login(c)

			// Assert the response
			assert.Equal(t, http.StatusUnauthorized, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Invalid password", response["error"])
		})
	})

	t.Run("should return error for invalid request body", func(t *testing.T) {
		mt.Run("invalid body", func(mt *mtest.T) {
			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create an invalid JSON request body
			invalidJSON := `{"email": "john.doe@example.com", "password":}`

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodPost, "/user/login", bytes.NewBufferString(invalidJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the Login handler
			Login(c)

			// Assert the response
			assert.Equal(t, http.StatusBadRequest, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Contains(t, response["error"], "invalid character")
		})
	})

}

// TestCreateUser tests the CreateUser handler
func TestCreateUser(t *testing.T) {
    mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

    gin.SetMode(gin.TestMode)

    t.Run("should create user successfully", func(t *testing.T) {
        mt.Run("success", func(mt *mtest.T) {
            mt.AddMockResponses(
                mtest.CreateCursorResponse(0, "speakup.users", mtest.FirstBatch),
            )
            
            insertedID := primitive.NewObjectID()
            mt.AddMockResponses(
                mtest.CreateSuccessResponse(
                    bson.E{Key: "ok", Value: 1},
                    bson.E{Key: "insertedId", Value: insertedID},
                ),
            )

            config.SetMongoClient(mt.Client)

            user := models.User{
                Name:     "John Doe",
                Email:    "john.doe@example.com",
                Password: "password123",
                Language: "en",
            }

            userJSON, _ := json.Marshal(user)
            req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userJSON))
            req.Header.Set("Content-Type", "application/json")

            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)
            c.Request = req

            CreateUser(c)

            assert.Equal(t, http.StatusOK, w.Code)
            
            var response map[string]string
            err := json.Unmarshal(w.Body.Bytes(), &response)
            assert.NoError(t, err)
            assert.Equal(t, "User created successfully", response["message"])
        })
    })
}

// TestGetUsers tests the GetUsers handler
func TestGetUsers(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	gin.SetMode(gin.TestMode)

	t.Run("should retrieve user successfully", func(t *testing.T) {
		mt.Run("find success", func(mt *mtest.T) {
			// Mock MongoDB find response
			user := models.User{
				ID:       "12345",
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: "hashedpassword",
				Language: "en",
			}
			var mockResponse bson.D
			bsonBytes, _ := bson.Marshal(user)
			_ = bson.Unmarshal(bsonBytes, &mockResponse)
			mt.AddMockResponses(mtest.CreateCursorResponse(1, "speakup.users", mtest.FirstBatch, mockResponse))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodGet, "/user/12345", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: "12345"}}

			// Call the GetUsers handler
			GetUsers(c)

			// Assert the response
			assert.Equal(t, http.StatusOK, w.Code)
			var response map[string]models.User
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, user, response["user"])
		})
	})

	t.Run("should return error for user not found", func(t *testing.T) {
		mt.Run("find failure", func(mt *mtest.T) {
			// Mock MongoDB find failure response
			mt.AddMockResponses(mtest.CreateCursorResponse(0, "speakup.users", mtest.FirstBatch))

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodGet, "/user/12345", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: "12345"}}

			// Call the GetUsers handler
			GetUsers(c)

			// Assert the response
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Failed to get user", response["error"])
		})
	})
}

// TestUpdateUser tests the UpdateUser handler
func TestUpdateUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	gin.SetMode(gin.TestMode)

	t.Run("should update user successfully", func(t *testing.T) {
		mt.Run("update success", func(mt *mtest.T) {
			// Mock MongoDB update response
			mt.AddMockResponses(mtest.CreateSuccessResponse())

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test user
			user := models.User{
				Name:     "Updated Name",
				Email:    "updated.email@example.com",
				Password: "updatedpassword123",
				Language: "fr",
			}

			// Marshal user to JSON
			userJSON, _ := json.Marshal(user)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodPut, "/user/12345", bytes.NewBuffer(userJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: "12345"}}

			// Call the UpdateUser handler
			UpdateUser(c)

			// Assert the response
			assert.Equal(t, http.StatusOK, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "User updated successfully", response["message"])
		})
	})

	t.Run("should return error for invalid request body", func(t *testing.T) {
		mt.Run("invalid body", func(mt *mtest.T) {
			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create an invalid JSON request body - CORRIGIDO: usando 'name' ao invés de 'first_name' e 'last_name'
			invalidJSON := `{"name": "Updated Name", "email":}`

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodPut, "/user/12345", bytes.NewBufferString(invalidJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: "12345"}}

			// Call the UpdateUser handler
			UpdateUser(c)

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

			// Create a test user
			user := models.User{
				Name:     "Updated Name",
				Email:    "updated.email@example.com",
				Password: "updatedpassword123",
				Language: "fr",
			}

			// Marshal user to JSON
			userJSON, _ := json.Marshal(user)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodPut, "/user/12345", bytes.NewBuffer(userJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: "12345"}}

			// Call the UpdateUser handler
			UpdateUser(c)

			// Assert the response
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Failed to update user", response["error"])
		})
	})
}

// TestDeleteUser tests the DeleteUser handler
func TestDeleteUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	gin.SetMode(gin.TestMode)

	t.Run("should delete user successfully", func(t *testing.T) {
		mt.Run("delete success", func(mt *mtest.T) {
			// Mock MongoDB delete response
			mt.AddMockResponses(mtest.CreateSuccessResponse())

			// Mock MongoDB client
			config.SetMongoClient(mt.Client)

			// Create a test HTTP request
			req, _ := http.NewRequest(http.MethodDelete, "/user/12345", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: "12345"}}

			// Call the DeleteUser handler
			DeleteUser(c)

			// Assert the response
			assert.Equal(t, http.StatusOK, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "User deleted successfully", response["message"])
		})
	})

	t.Run("should return error for database failure", func(t *testing.T) {
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
			req, _ := http.NewRequest(http.MethodDelete, "/user/12345", nil)

			// Create a test HTTP response recorder
			w := httptest.NewRecorder()

			// Create a Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: "12345"}}

			// Call the DeleteUser handler
			DeleteUser(c)

			// Assert the response
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Failed to delete user", response["error"])
		})
	})
}
