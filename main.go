package main

import (
	"container/list"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

const INSERT_COUNT = 100

type EventLog struct {
	Name  string
	Value string
}

func NewEventLog(name, value string) EventLog {
	return EventLog{
		Name:  name,
		Value: value,
	}
}

func RunDB(db *sql.DB, eventlogStack *list.List) {
	query := "INSERT INTO eventlog(at, name, value) values"
	for i := 0; i < INSERT_COUNT; i++ {
		query += "(NOW(), ?, ?)"
		if i != INSERT_COUNT {
			query += ","
		}
	}
	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	v := []string{}
	for i := 0; i < INSERT_COUNT; i++ {
		event := eventlogStack.Remove(eventlogStack.Front())
		eventLog := event.(EventLog)
		name := eventLog.Name
		value := eventLog.Value
		v = append(v, name)
		v = append(v, value)
	}
	s := make([]interface{}, len(v))
	for i, v := range v {
		s[i] = v
	}
	_, err = stmt.Exec(s...)
	if err != nil {
		panic(err)
	}
}

func main() {
	dataSourceName := os.Getenv("HAKARU_DATASOURCENAME")
	if dataSourceName == "" {
		dataSourceName = "root:hakaru-pass@tcp(127.0.0.1:13306)/hakaru-db"
	}

	eventlogsStack := list.New()
	_ = eventlogsStack

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var mux sync.Mutex
	requestCount := 0
	hakaruHandler := func(w http.ResponseWriter, r *http.Request) {
		go func() {
			name := r.URL.Query().Get("name")
			value := r.URL.Query().Get("value")

			eventlogsStack.PushFront(NewEventLog(name, value))
			mux.Lock()
			defer mux.Unlock()
			requestCount += 1
			if requestCount%INSERT_COUNT == 0 {
				RunDB(db, eventlogsStack)
			}
		}()

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
