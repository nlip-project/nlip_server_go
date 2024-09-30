package handlers

import (
	"net/http"
	"nlip/initializers"
	"nlip/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body AuthBody
	var hash []byte
	var err error

	err = handleAuthBody(w, r, &body)
	if err != nil {
		return
	}

	// bcrypt already handles salting...
	hash, err = bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		http.Error(w, "Failed to hash pasword. "+err.Error(), http.StatusBadRequest)
		return
	}

	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user) // pass pointer of data to Create

	if result.Error != nil {
		http.Error(w, "Failed to create user. "+result.Error.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body AuthBody
	err := handleAuthBody(w, r, &body)
	if err != nil {
		return
	}

	user := GetUser(&body)
	if user == nil {
		prepareTextResponse(w, http.StatusNotFound, "User does not exist.")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// If mobile, send in JSON
	resp := map[string]string{
		"token": tokenString,
	}
	prepareJSONResponse(w, http.StatusOK, resp)

	// TODO: Consider the below code later.

	// 	userAgent := r.Header.Get("User-Agent")
	// 	if strings.Contains(userAgent, "Mobile") {
	// 		// If mobile, send in JSON
	// 		resp := map[string]string{
	// 			"token": tokenString,
	// 		}
	// 		prepareJSONResponse(w, http.StatusOK, resp)
	// 	} else {
	// 		// If web,
	// 		cookie := &http.Cookie{
	// 			Name:     "access_token",
	// 			Value:    tokenString,
	// 			Path:     "/",
	// 			HttpOnly: true, // Prevent JavaScript access to the cookie
	// 			Secure:   true, // Send cookie over HTTPS only
	// 			SameSite: http.SameSiteStrictMode,
	// 			Expires:  time.Now().Add(time.Hour * 24 * 30), // how long to store on the browser
	// 		}
	// 		http.SetCookie(w, cookie)
	// 		w.WriteHeader(http.StatusOK)
	// 	}
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
