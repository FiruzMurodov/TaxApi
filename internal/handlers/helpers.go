package handlers

import (
	"encoding/json"
	"taxApi/internal/models"
)

func ConvertToJson(message string) []byte {
	messageResponse := models.MessageToUser{Message: message}
	jsonResponse, _ := json.Marshal(messageResponse)
	return jsonResponse
}
