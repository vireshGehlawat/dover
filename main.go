package main

import (
	"dover/api"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func main() {

	r := mux.NewRouter()

	// add proper initialization flow
	api.InitializeRoutes(r, api.New())

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// initialize command line utils

	fmt.Println(srv.ListenAndServe())
}
