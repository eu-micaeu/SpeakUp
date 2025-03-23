package handlers

import (
	"net/http"
	"speakup/models"
	"github.com/gin-gonic/gin"
	"speakup/config"
	"speakup/middlewares"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Login handles the login of a user
func Login(c *gin.Context) {
	client := config.GetMongoClient()
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	collection := client.Database("speakup").Collection("users")
	var result models.User
	err := collection.FindOne(c, bson.M{"email": user.Email}).Decode(&result)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}
	token, err := middlewares.GenerateJWT(result.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

// CreateUser handles the creation of a new user
func CreateUser(c *gin.Context) {
	client := config.GetMongoClient()
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	collection := client.Database("speakup").Collection("users")
	user.ID = uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)
	_, err = collection.InsertOne(c, bson.M{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"password":   user.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// GetUsers handles the retrieval of a user
func GetUsers(c *gin.Context) {
	client := config.GetMongoClient()
	id := c.Param("id")
	collection := client.Database("speakup").Collection("users")
	var user models.User
	err := collection.FindOne(c, bson.M{"id": id}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUser handles the update of a user
func UpdateUser(c *gin.Context) {
	client := config.GetMongoClient()
	id := c.Param("id")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	collection := client.Database("speakup").Collection("users")
	_, err := collection.UpdateOne(c, bson.M{"id": id}, bson.M{"$set": bson.M{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
	}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
}

// DeleteUser handles the deletion of a user
func DeleteUser(c *gin.Context) {
	client := config.GetMongoClient()
	id := c.Param("id")
	collection := client.Database("speakup").Collection("users")
	_, err := collection.DeleteOne(c, bson.M{"id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}