package main

import (
	"net/http"
	"log"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	hakaruHandler := func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("mysql", "root:hakaru-pass@/hakarudb")
		if err != nil {
			panic(err.Error())
		}
		defer db.Close() // 関数がリターンする直前に呼び出される

		// TODO: insertする
		_, e := db.Query("SELECT * FROM users") //
		if e != nil {
			panic(e.Error())
		}

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

	// start server
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
