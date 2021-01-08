package main

import (
	"dover/api"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
)

func main() {

	r := mux.NewRouter()

	db := sqlx.MustConnect("mysql", "root@(localhost:3306)/dover")
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
