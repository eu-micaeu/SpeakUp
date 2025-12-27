package handlers

import (
	"net/http"

	"speakup/pkg/middlewares"
	"speakup/pkg/models"
	"speakup/pkg/repositories"
	"speakup/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repo repositories.UserRepository
}

func NewUserHandler(repo repositories.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

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
func (h *UserHandler) Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.Repo.FindByEmail(c, user.Email)
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid email")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid password")
		return
	}

	token, err := middlewares.GenerateJWT(result.ID, result.Name, result.Email, result.Language, result.Level)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	c.Set("authToken", token)
	utils.RespondWithJSON(c, http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

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
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if a user with the same email already exists
	_, err := h.Repo.FindByEmail(c, user.Email)
	if err == nil {
		utils.RespondWithError(c, http.StatusConflict, "User with this email already exists")
		return
	}
	if err != mongo.ErrNoDocuments {
		utils.RespondWithError(c, http.StatusInternalServerError, "Database error")
		return
	}

	// Generate a unique ID for the user
	user.ID = uuid.New().String()

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	user.Password = string(hashedPassword)

	// Insert the user into the database
	err = h.Repo.Create(c, user)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{"message": "User created successfully"})
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
func (h *UserHandler) GetUsers(c *gin.Context) {
	id := c.Param("id")
	user, err := h.Repo.FindByID(c, id)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get user")
		return
	}
	utils.RespondWithJSON(c, http.StatusOK, gin.H{"user": user})
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
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
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
			utils.RespondWithError(c, http.StatusInternalServerError, "Failed to hash password")
			return
		}
		updateDoc["password"] = string(hashedPassword)
	}

	// Se não houver campos para atualizar, retorna sucesso
	if len(updateDoc) == 0 {
		utils.RespondWithJSON(c, http.StatusOK, gin.H{"message": "No fields to update"})
		return
	}

	err := h.Repo.Update(c, id, updateDoc)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{"message": "User updated successfully"})
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
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := h.Repo.Delete(c, id)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}
	utils.RespondWithJSON(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}
