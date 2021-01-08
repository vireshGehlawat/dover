package main

import (
	"dover/api"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"time"
)

func main() {

	r := mux.NewRouter()

	db := sqlx.MustConnect("sqlite3", ":memory:")
	fmt.Println(db.Ping())

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
