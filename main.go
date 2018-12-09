package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Value struct {
	Name  string    `db:"name"`
	Value string    `db:"value"`
	Now   time.Time `db:"now"`
}

var (
	db    *sqlx.DB
	queCh = make(chan Value, 1000)
)

func inserter() {
	ticker := time.NewTicker(1 * time.Second)
	valueQue := make([]Value, 0, 1000)
	for {
		select {
		case <-ticker.C:
			// TODO: Insert 処理
			//stmt, e := db.Prepare("INSERT INTO eventlog(at, name, value) values(NOW(), ?, ?)")
			//if e != nil {
			//	panic(e.Error())
			//}
			if len(valueQue) == 0 {
				continue
			}
			query := "INSERT INTO eventlog(at, name, value) values(?, ?, ?)" + strings.Repeat(", (?, ?, ?)", len(valueQue)-1)
			stmt, e := db.Prepare(query)
			if e != nil {
				panic(e.Error())
			}

			defer stmt.Close()

			args := make([]interface{}, 3*len(valueQue))
			for i, que := range valueQue {
				args[i] = que.Now
				args[i+1] = que.Name
				args[i+2] = que.Value
			}
			_, _ = stmt.Exec(args...)
			valueQue = make([]Value, 0, 1000)

		case que := <-queCh:
			valueQue = append(valueQue, que)

		}
	}
}

func hakaruHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	name := r.URL.Query().Get("name")
	value := r.URL.Query().Get("value")

	queCh <- Value{
		Now:   now,
		Name:  name,
		Value: value,
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

func main() {
	dataSourceName := os.Getenv("HAKARU_DATASOURCENAME")
	if dataSourceName == "" {
		dataSourceName = "root:hakaru-pass@tcp(127.0.0.1:13306)/hakaru-db"

	}

	_db, err := sqlx.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	defer _db.Close()

	db = _db

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)

	http.HandleFunc("/hakaru", hakaruHandler)
	http.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })

	// start server
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
