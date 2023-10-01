package server

import (
	"net/http"
	"poppypaw/server/endpoints"
)

func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

func Initialize() http.Handler {
	// Create a new ServeMux to register the handlers
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", CORS(endpoints.Ping))
	mux.HandleFunc("/subscribe", CORS(endpoints.Subscribe))
	mux.HandleFunc("/notification", CORS(endpoints.Notification))

	return mux
}
