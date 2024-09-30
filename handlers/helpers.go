package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nlip/models"

	"github.com/go-playground/validator/v10"
)

func init() {
	validate = validator.New()
}

func prepareJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func prepareTextResponse(w http.ResponseWriter, status int, data string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(data))
}

// Documentation recommends single instance as it does caching
var validate *validator.Validate

func decodeJSONBody(r *http.Request, object interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // @Sugih mentioned maybe removing this
	return decoder.Decode(&object)
}

func validateStruct(object interface{}) error {
	return validate.Struct(object)
}

var getFormat = map[models.Subformat]models.Format{
	models.English:   models.Text,
	models.Spanish:   models.Text,
	models.German:    models.Text,
	models.JSON:      models.Structured,
	models.URI:       models.Structured,
	models.XML:       models.Structured,
	models.HTML:      models.Structured,
	models.ImageBMP:  models.Binary,
	models.ImageGIF:  models.Binary,
	models.ImageJPEG: models.Binary,
	models.ImageJPG:  models.Binary,
	models.ImagePNG:  models.Binary,
	models.ImageTIFF: models.Binary,
	models.AudioMP3:  models.Binary,
	models.TextSF:    models.Location,
	models.GPS:       models.Location,
}

func validateMessage(msg *models.Message) error {
	// Using a map here to map from subformat to format.
	if msg.Format != getFormat[msg.Subformat] {
		return fmt.Errorf("subformat incompatible with Format type")
	}

	// Need to make sure Format and Subformat are valid!

	// TODO: Add if control != nil && control == true check here
	// so that logistical requests (server policy, parameter negotiating)
	// can be handled here.
	return nil
}

func handleIncomingMessage(w http.ResponseWriter, r *http.Request, msg *models.Message) error {
	err := decodeJSONBody(r, msg)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return err
	}
	err = validateStruct(msg)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return err
	}
	err = validateMessage(msg)
	if err != nil {
		http.Error(w, "Erroneous message format", http.StatusBadRequest)
		return err
	}
	return nil
}

func handleAuthBody(w http.ResponseWriter, r *http.Request, body *AuthBody) error {
	err := decodeJSONBody(r, body)
	if err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return err
	}
	err = validateStruct(body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return err
	}
	return nil
}
