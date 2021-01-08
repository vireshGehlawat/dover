package main

import (
	"bufio"
	"dover/api"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"
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

	cliArguments := os.Args[1:]
	if len(cliArguments) > 0 {
		if cliArguments[0] == "ingestprofiles" {
			file, err := os.Open(cliArguments[1])
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	fmt.Println(srv.ListenAndServe())
}
