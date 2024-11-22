package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/coreos/go-oidc"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

var (
	redirectURI      = "https://druid.eecs.umich.edu:80/auth/"
	oidcDiscoveryURL = "http://127.0.0.1:8080/o"
)
var (
	clientID     string
	clientSecret string
	oauth2Config *oauth2.Config
	oidcProvider *oidc.Provider
	verifier     *oidc.IDTokenVerifier
)

func init() {
	godotenv.Load()
	clientID = os.Getenv("OAUTH_CLIENT_ID")
	clientSecret = os.Getenv("OAUTH_CLIENT_SECRET")

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, oidcDiscoveryURL)
	if err != nil {
		log.Fatalf("Failed to get provider: %v", err)
	}
	oidcProvider = provider

	oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "read", "write"},
		Endpoint:     provider.Endpoint(),
	}

	verifier = provider.Verifier(&oidc.Config{ClientID: clientID})
}

func generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func generateCodeVerifier() (string, error) {
	return generateRandomString(43)
}

func deriveCodeChallenge(verifier string) string {
	h := sha256.New()
	h.Write([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func HandleLogin(c echo.Context) error {
	state, err := generateRandomString(16)
	if err != nil {
		return c.String(500, fmt.Sprintf("Failed to generate state: %v", err))
	}
	nonce, err := generateRandomString(16)
	if err != nil {
		return c.String(500, fmt.Sprintf("Failed to generate nonce: %v", err))
	}

	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		return c.String(500, fmt.Sprintf("Failed to generate code verifier: %v", err))
	}
	codeChallenge := deriveCodeChallenge(codeVerifier)

	// These are to make sure nothing is intercepted
	c.SetCookie(&http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
		Path:     "/",
	})
	c.SetCookie(&http.Cookie{
		Name:     "oauth_nonce",
		Value:    nonce,
		HttpOnly: true,
		Path:     "/",
	})
	c.SetCookie(&http.Cookie{
		Name:     "oauth_code_verifier",
		Value:    codeVerifier,
		HttpOnly: true,
		Path:     "/",
	})

	authCodeURL := oauth2Config.AuthCodeURL(state,
		oauth2.SetAuthURLParam("nonce", nonce),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
	fmt.Println("Generated Authorization URL:", authCodeURL)
	return c.Redirect(302, authCodeURL)
}

func HandleCallback(c echo.Context) error {
	ctx := context.Background()

	stateCookie, err := c.Cookie("oauth_state")
	if err != nil {
		return c.String(400, "State cookie not found")
	}
	nonceCookie, err := c.Cookie("oauth_nonce")
	if err != nil {
		return c.String(400, "Nonce cookie not found")
	}

	// needed for CSRF attacks
	state := c.QueryParam("state")
	if state != stateCookie.Value {
		return c.String(400, "State parameter mismatch")
	}

	// to ensure intercepted requests still can't get a token
	codeVerifierCookie, err := c.Cookie("oauth_code_verifier")
	if err != nil {
		return c.String(400, "Code verifier cookie not found")
	}
	codeVerifier := codeVerifierCookie.Value

	code := c.QueryParam("code")
	if code == "" {
		return c.String(400, "Code not found in callback")
	}

	token, err := oauth2Config.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		return c.String(500, fmt.Sprintf("Failed to exchange token: %v", err))
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return c.String(500, "No id_token field in oauth2 token.")
	}

	idToken, err := verifier.Verify(ctx, rawIDToken) // Verifier comes back here
	if err != nil {
		return c.String(500, fmt.Sprintf("Failed to verify ID Token: %v", err))
	}

	var claims struct {
		Nonce string `json:"nonce"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return c.String(500, fmt.Sprintf("Failed to parse claims: %v", err))
	}
	if claims.Nonce != nonceCookie.Value {
		return c.String(400, "Nonce mismatch")
	}

	userInfo, err := oidcProvider.UserInfo(ctx, oauth2Config.TokenSource(ctx, token))
	if err != nil {
		return c.String(500, fmt.Sprintf("Failed to get userinfo: %v", err))
	}

	var userClaims map[string]interface{}
	if err := userInfo.Claims(&userClaims); err != nil {
		return c.String(500, fmt.Sprintf("Failed to parse userinfo claims: %v", err))
	}

	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    token.AccessToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	return c.JSON(200, echo.Map{
		"message":   "User authenticated",
		"user_info": userClaims,
	})
}

func ProtectedHandler(c echo.Context) error {
	accessTokenCookie, err := c.Cookie("access_token")
	if err != nil {
		return c.String(401, "Access token cookie missing")
	}
	accessToken := accessTokenCookie.Value

	ctx := context.Background()
	_, err = oidcProvider.UserInfo(ctx, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accessToken,
	}))
	if err != nil {
		return c.String(401, fmt.Sprintf("Failed to complete verification: %v", err))
	}

	// Don't really need below because don't need user info
	// var claims map[string]interface{}
	// if err := userInfo.Claims(&claims); err != nil {
	// 	return c.String(500, fmt.Sprintf("Failed to parse userinfo claims: %v", err))
	// }

	return c.JSON(200, echo.Map{
		"message": "User token successfully mapped to a user.",
	})
}
