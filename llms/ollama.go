package llms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"` // one-go or as a stream? for UI or large requests
}

// Simple request took ~2.5 sec
// If "Stream" is on, will likely work better for client
func GetResponse(payload *OllamaRequest) (string, error) {
	// Important TODO: If Ollama does not have desired model, request silently fails and code keeps running.
	// Must have an error check otherwise will cause all sorts of bugs.
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	response, ok := result["response"].(string)
	if !ok {
		return "", nil
	}

	fmt.Println("Response is ", response)
	return response, nil
}
