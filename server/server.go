package server

import (
	"encoding/json"
	"net/http"
)

// AnnounceResponse is a JSON response returned when a GET is issued to the /announce endpoint
type AnnounceResponse struct {
	Title string `json:"title"`
}

func createHandler(title string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data := AnnounceResponse{
			Title: title,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
	}
}

func Start(addr, title string) error {
	http.HandleFunc("/announce", createHandler(title))
	return http.ListenAndServe(addr, nil)
}
