package handlers

import (
	"fmt"
	"net/http"
	"nlip/llms"
	"nlip/models"
)

func InitiationHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// This does not produce err. Can it fail though?
	var msg models.Message
	err := handleIncomingMessage(w, r, &msg)
	if err != nil {
		return
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
	prepareJSONResponse(w, http.StatusOK, response)
}

// Interesting use case
func TestHandler(w http.ResponseWriter, r *http.Request) {
	var msg models.Message
	err := handleIncomingMessage(w, r, &msg)
	if err != nil {
		return
	}

	payload := llms.OllamaRequest{
		// Model must exist on the machine
		Model:  "llama3.2",
		Prompt: msg.Content,
		Stream: false,
	}

	resp, err := llms.GetResponse(&payload)
	if err != nil {
		http.Error(w, "Bad request! "+err.Error(), http.StatusBadRequest)
		return
	}

	jsonResp := models.Message{
		Format:    "text",
		Subformat: "english",
		Content:   resp,
	}

	prepareJSONResponse(w, http.StatusOK, jsonResp)
}
