package main

import (
	"net/http"
	"log"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

func main() {
	dataSourceName := "root:hakaru-pass@tcp(127.0.0.1:13306)/hakaru-db"
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

	probe := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("on probe:")

		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		rows, e := db.Query("SELECT name, value FROM eventlog")
		if e != nil {
			panic(e.Error())
		}

		for rows.Next() {
			var name string
			var value int

			if err := rows.Scan(&name, &value); err != nil {
				log.Fatal(err)
			}
			fmt.Println(name, value)
		}
	}


	http.HandleFunc("/hakaru", hakaruHandler)
	http.HandleFunc("/probe", probe)

	// start server
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
