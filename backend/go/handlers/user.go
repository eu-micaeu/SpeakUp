package handlers

import (
	"net/http"
	"speakup/config"
	"speakup/middlewares"
	"speakup/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	token, err := middlewares.GenerateJWT(result.ID, result.Email, result.Language, result.Level)
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
// @Tags         User
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

	collection := client.Database("speakup").Collection("users")

	// Check if a user with the same email already exists
	var existingUser models.User
	err := collection.FindOne(c, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}
	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
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

	// Insert the user into the database
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
// @Tags         User
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
// @Tags         User
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

	// Preparar o documento de atualização apenas com campos fornecidos
	updateDoc := bson.M{}

	if user.Name != "" {
		updateDoc["name"] = user.Name
	}
	if user.Email != "" {
		updateDoc["email"] = user.Email
	}
	if user.Language != "" {
		updateDoc["language"] = user.Language
	}
	if user.Level != "" {
		updateDoc["level"] = user.Level
	}
	if user.Password != "" {
		// Hash do novo password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		updateDoc["password"] = string(hashedPassword)
	}

	// Se não houver campos para atualizar, retorna sucesso
	if len(updateDoc) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No fields to update"})
		return
	}

	collection := client.Database("speakup").Collection("users")
	_, err := collection.UpdateOne(c, bson.M{"id": id}, bson.M{"$set": updateDoc})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser handles the deletion of a user
// DeleteUser godoc
// @Summary      Delete a user
// @Description  Delete a user by their ID
// @Tags         User
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
