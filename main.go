package main

import (
	"bufio"
	"dover/api"
	"dover/services"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"
	"time"
)

func main() {
	db := sqlx.MustConnect("mysql", "root@(localhost:3306)/dover")
	ingestion := services.New(db)
	if len(os.Args) == 1 {
		// init API and block to listen
		initializeForAPIRole(db)
	}
	cliArguments := os.Args[1:]
	if cliArguments[0] == "ingestprofiles" {
		if len(cliArguments) < 2 {
			fmt.Println("invalid usage of the role, sample start command\n" +
				" ./main ingestprofiles ./path/to/file.json\n")
			return
		}
		file, err := os.Open(cliArguments[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		err = ingestion.IngestBulk(scanner)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func initializeForAPIRole(db *sqlx.DB) {
	// add proper initialization flow for router
	r := mux.NewRouter()
	profilesService := services.NewProfilesService(db)
	view := api.New(profilesService)
	api.InitializeRoutes(r, view)
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println(srv.ListenAndServe())
}
