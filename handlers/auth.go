package handlers

import (
	"net/http"
	"nlip/initializers"
	"nlip/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func Register(c echo.Context) error {
	// Bind the incoming JSON request body to AuthBody struct
	var body AuthBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Hash the password using bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to hash password: " + err.Error()})
	}

	// Create a new user record
	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user) // pass pointer of data to Create

	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to create user: " + result.Error.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func Login(c echo.Context) error {
	// Bind the incoming JSON request body to AuthBody struct
	var body AuthBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Retrieve the user from the database and validate the password
	user := GetUser(&body)
	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User does not exist."})
	}

	// Create a JWT token with a 30-day expiration
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign the token with a secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Respond with the token in JSON format
	resp := map[string]string{
		"token": tokenString,
	}
	return c.JSON(http.StatusOK, resp)
}

// Returns a "copy", not direct User reference from the DB!
func GetUser(body *AuthBody) *models.User {
	var user models.User
	result := initializers.DB.First(&user, "email = ?", body.Email)
	if result.Error != nil {
		return nil
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return nil
	}
	return &user
}
