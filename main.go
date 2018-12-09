package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type EventLog struct {
	Name  string    `db:"name"`
	Value string    `db:"value"`
	Now   time.Time `db:"at"`
}

var (
	db        *sql.DB
	eventChan = make(chan EventLog, 500)
)

func timerInsert() {
	ticker := time.NewTicker(1 * time.Second)
	querys := make([]EventLog, 0, 500)
	for {
		select {
		case <-ticker.C:
			if len(querys) == 0 {
				continue
			}
			sqlStatement := "INSERT INTO eventlog(at, name, value) values(?, ?, ?)" + strings.Repeat(", (?,?,?)", len(querys)-1)
			values := make([]interface{}, 3*len(querys))
			for i, query := range querys {
				values[3*i] = query.Now
				values[3*i+1] = query.Name
				values[3*i+2] = query.Value
			}
			_, err := db.Exec(sqlStatement, values...)
			if err != nil {
				fmt.Println("SQL could not be executed")
			}
		case event := <-eventChan:
			querys = append(querys, event)
		}
	}
}

func main() {
	dataSourceName := os.Getenv("HAKARU_DATASOURCENAME")
	if dataSourceName == "" {
		dataSourceName = "root:hakaru-pass@tcp(127.0.0.1:13306)/hakaru-db"
	}

	_db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("DB could not be opened")
	}
	defer _db.Close()

	db = _db

	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)

	go timerInsert()

	hakaruHandler := func(w http.ResponseWriter, r *http.Request) {
		var event EventLog
		event.Name = r.URL.Query().Get("name")
		event.Value = r.URL.Query().Get("value")
		event.Now = time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
		eventChan <- event

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
