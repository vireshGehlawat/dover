package api

import (
	"github.com/gorilla/mux"
)

func InitializeRoutes(r *mux.Router, view View) {
	r.HandleFunc("/api/profiles/", view.GetProfileList).Methods("GET")
	r.HandleFunc("/api/profiles/{id}/", view.GetProfile).Methods("GET")
}
