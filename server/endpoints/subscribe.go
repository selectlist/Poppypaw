package endpoints

import (
	"encoding/json"
	"io"
	"net/http"

	"poppypaw/state"
	"strings"
)

func Subscribe(w http.ResponseWriter, r *http.Request) {
	type Credentials struct {
		DeviceTokens []string `json:"deviceTokens"`
	}

	type Body struct {
		ProviderID  string      `json:"providerId"`
		Credentials Credentials `json:"credentials"`
	}

	// Construct the URL with the correct value from the "id" query parameter
	url := "https://api.novu.co/v1/subscribers/" + r.URL.Query().Get("id") + "/credentials"

	// Create a JSON payload based on the request query parameter "token"
	jsonData, err := json.Marshal(Body{
		ProviderID: "fcm",
		Credentials: Credentials{
			DeviceTokens: []string{r.URL.Query().Get("token")},
		},
	})
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Create an HTTP PUT request with the constructed URL and JSON payload
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(jsonData)))
	if err != nil {
		http.Error(w, "Failed to create HTTP request", http.StatusInternalServerError)
		return
	}

	// Add necessary headers to the request
	req.Header.Add("Authorization", "ApiKey "+state.Config.NovuAPIKey)
	req.Header.Add("Content-Type", "application/json")

	// Send the HTTP request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send HTTP request", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	// Set the "Content-Type" header to "application/json"
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data from the response body to the response
	w.Write(responseBody)
}
