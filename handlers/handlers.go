package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"nlip/llms"
	"nlip/models"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
)

func StartConversationHandler(c echo.Context) error {
	var msg models.Message
	if err := c.Bind(&msg); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := validate.Struct(&msg); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed", "details": err.Error()})
	}

	// For testing right now
	fmt.Println(msg.Control)
	fmt.Println(msg.Format)
	fmt.Println(msg.Subformat)
	fmt.Println(msg.Content)
	fmt.Println(msg.Submessages)
	fmt.Println(msg.Token)
	fmt.Println(msg.Subtokens)

	// Dummy, hardcoded response:
	response := &models.Message{
		Format:    models.Text,
		Subformat: models.English,
		Content: "Use Authentication token 0x0567564.\n" +
			"Authentication-token must be specified.\n" +
			"Only last 5 exchanges will be remembered by the server.\n" +
			"You need to remember and provide all exchanges older than the last 5.",
	}
	return c.JSON(http.StatusOK, response)
}

func TextHandler(c echo.Context) error {
	var msg models.Message
	if err := c.Bind(&msg); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if msg.Format != "text" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Text endpoint received non-text data"})
	}

	if err := validate.Struct(&msg); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed", "details": err.Error()})
	}

	payload := llms.OllamaRequest{
		Model:  "llama3.2", // Model must exist on the machine
		Prompt: msg.Content,
		Stream: false,
	}

	resp, err := llms.GetTextResponse(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request: " + err.Error()})
	}

	jsonResp := models.Message{
		Format:    "text",
		Subformat: "english",
		Content:   resp,
	}

	return c.JSON(http.StatusOK, jsonResp)
}

func ImageHandler(c echo.Context) error {
	var msg models.Message
	if err := c.Bind(&msg); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := validate.Struct(&msg); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed", "details": err.Error()})
	}

	if msg.Format != "binary" || !isValidImageSubformat(msg.Subformat) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid format or subformat"})
	}

	imageData, err := base64.StdEncoding.DecodeString(msg.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unable to decode base64 content"})
	}

	uniqueID := uuid.New().String()
	extension := strings.ToLower(string(msg.Subformat))
	filename := fmt.Sprintf("%s.%s", uniqueID, extension)
	basePath := "/Users/hbzengin/src/go-server-example/uploads"
	filepath := filepath.Join(basePath, filename)

	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		if err := os.Mkdir(basePath, 0755); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Unable to create uploads directory",
			})
		}
	}

	if err := os.WriteFile(filepath, imageData, 0644); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to save file"})
	}

	ollamaPrompt := "What do you see in this image?"
	payload := llms.OllamaRequest{
		Model:  "llava",
		Prompt: ollamaPrompt,
		Image:  filepath,
		Stream: false,
	}

	ollamaResponse, err := llms.GetImageResponse(&payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get response from Ollama", "details": err.Error()})
	}

	jsonResp := models.Message{
		Format:    "text",
		Subformat: "english",
		Content:   ollamaResponse,
	}

	return c.JSON(http.StatusOK, jsonResp)
}

func isValidImageSubformat(subformat models.Subformat) bool {
	switch strings.ToLower(string(subformat)) {
	case "jpeg", "jpg", "png", "gif", "bmp":
		return true
	default:
		return false
	}
}
