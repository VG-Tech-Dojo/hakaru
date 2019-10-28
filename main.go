package main

import (
	"net/http"
	"log"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"os"
)

func main() {
	dataSourceName := os.Getenv("HAKARU_DATASOURCENAME")
	if dataSourceName == "" {
		dataSourceName = "root:password@tcp(127.0.0.1:13306)/hakaru"
	}

	hakaruHandler := func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		stmt, e := db.Prepare("INSERT INTO eventlog(at, name, value) values(NOW(), ?, ?)")
		if e != nil {
			panic(e.Error())
		}

		defer stmt.Close()

		name := r.URL.Query().Get("name")
		value := r.URL.Query().Get("value")

		_, _ = stmt.Exec(name, value)

		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
	}

	http.HandleFunc("/hakaru", hakaruHandler)
	http.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })

	// start server
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
