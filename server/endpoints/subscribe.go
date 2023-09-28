package endpoints

import (
	"encoding/json"
	"net/http"
)

// WIP
func Subscribe(w http.ResponseWriter, r *http.Request) {
	// Convert the bot struct to JSON
	jsonData, err := json.Marshal(map[string]interface{}{
		"message": "PONG!",
	})
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Set the "Content-Type" header to "application/json"
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response
	w.Write(jsonData)
}
