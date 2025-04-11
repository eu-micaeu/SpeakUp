package handlers

import (
	"net/http"
	"speakup/config"
	"speakup/middlewares"
	"speakup/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Login handles the login of a user
// Login godoc
// @Summary      Login de usuário
// @Description  Autentica um usuário e retorna um token JWT
// @Tags         Autenticação
// @Accept       json
// @Produce      json
// @Param         credentials body object{email=string,password=string} true "Credenciais de login"
// @Example      {object} credentials {"email":"user@example.com","password":"123456"}
// @Success      200         {object}  object{token=string}
// @Router       /user/login [post]
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
	token, err := middlewares.GenerateJWT(result.ID, result.Email, result.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.Set("authToken", token)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

// CRUD operations for user

// CreateUser handles the creation of a new user
// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user with the provided details
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body models.User true "User details"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /user [post]
func CreateUser(c *gin.Context) {
	client := config.GetMongoClient()
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a unique ID for the user
	user.ID = uuid.New().String()

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	// Insert the user object directly into the database
	collection := client.Database("speakup").Collection("users")
	_, err = collection.InsertOne(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// GetUsers handles the retrieval of a user
// GetUsers godoc
// @Summary      Get a user by ID
// @Description  Retrieve a user's details by their ID
// @Tags         users
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200  {object}  models.User
// @Failure      500  {object}  object{error=string}
// @Router       /user/{id} [get]
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
// UpdateUser godoc
// @Summary      Update a user
// @Description  Update a user's details by their ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        user body models.User true "User details"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /user/{id} [put]
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
		"password":   user.Password,
		"language":   user.Language,
	}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
		return
	}
}

// DeleteUser handles the deletion of a user
// DeleteUser godoc
// @Summary      Delete a user
// @Description  Delete a user by their ID
// @Tags         users
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200  {object}  object{message=string}
// @Failure      500  {object}  object{error=string}
// @Router       /user/{id} [delete]
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
