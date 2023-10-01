package endpoints

import (
	"context"
	"encoding/json"
	"net/http"

	"poppypaw/state"

	"github.com/getsentry/sentry-go"
	novu "github.com/novuhq/go-novu/lib"
)

func Notification(w http.ResponseWriter, r *http.Request) {
	// Check if the request is not coming from localhost
	if r.Host != "localhost:" + state.Config.Port && r.Host != "127.0.0.1:" + state.Config.Port {
		http.Error(w, "Page not Found.", http.StatusNotFound)
		return
	}

	// Sentry
	sentry.Init(sentry.ClientOptions{
		Dsn:              "https://0af9ed73773d95692388f88fb8b68777@trace.select-list.xyz/5",
		TracesSampleRate: 1.0,
	})

	// Extract data from the request URL
	subscriberID := r.URL.Query().Get("id")
	title := r.URL.Query().Get("title")
	description := r.URL.Query().Get("description")
	eventName := r.URL.Query().Get("event_name")

	if subscriberID == "" {
		http.Error(w, "Missing Fields: id", http.StatusBadRequest)
		return
	} 
	
	if title == "" {
		http.Error(w, "Missing Fields: title", http.StatusBadRequest)
		return
	}
	
	if description == "" {
		http.Error(w, "Missing Fields: description", http.StatusBadRequest)
		return
	}
	
	if eventName == "" {
		http.Error(w, "Missing Fields: event_name", http.StatusBadRequest)
		return
	}

	// Prepare data for Novu
	ctx := context.Background()
	to := map[string]interface{}{
		"subscriberId": subscriberID,
	}

	payload := map[string]interface{}{
		"title":       title,
		"description": description,
	}

	data := novu.ITriggerPayloadOptions{To: to, Payload: payload}

	// Trigger the Novu event
	var err error
	_, err = state.Novu.EventApi.Trigger(ctx, eventName, data)
	if err != nil {
		// Capture and report the error to Sentry
		sentry.CaptureException(err)
		http.Error(w, "Failed to trigger the event", http.StatusInternalServerError)
		return
	}

	// Response
	responseData := map[string]interface{}{
		"success": true,
	}

	// Set the "Content-Type" header to "application/json"
	w.Header().Set("Content-Type", "application/json")

	// Serialize the JSON response and send it to the client
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}
