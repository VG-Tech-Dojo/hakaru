package main

import (
	"container/list"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const INSERT_TIME = 10

type EventLog struct {
	Name  string
	Value string
	Now   string
}

func NewEventLog(name, value string) EventLog {
	// TODO mysqlのNOW()が何を作っているの確認
	// https://stackoverflow.com/questions/23415612/insert-datetime-using-now-with-go
	return EventLog{
		Name:  name,
		Value: value,
		Now:   time.Now().Format(time.RFC3339),
	}
}

func RunDB(db *sql.DB, eventlogStack *list.List, queCount int) {
	query := "INSERT INTO eventlog(at, name, value) values"
	for i := 0; i < queCount; i++ {
		query += "(NOW(), ?, ?)"
		if i != queCount {
			query += ","
		}
	}
	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	v := []string{}
	for i := 0; i < queCount; i++ {
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

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// insertの通知をするためのgoroutine
	requestCh := make(chan *http.Request)
	// バルクインサート周りを管理するgoroutine
	go func(requestCh chan *http.Request) {
		timeNotification := time.NewTicker(INSERT_TIME * time.Second)
		queCount := 0
		eventlogStack := list.New()
		for {
			select {
			case r := <-requestCh:
				name := r.URL.Query().Get("name")
				value := r.URL.Query().Get("value")
				eventlogStack.PushBack(NewEventLog(name, value))
				queCount += 1
			case <-timeNotification.C:
				if queCount != 0 {
					RunDB(db, eventlogStack, queCount)
					eventlogStack = list.New()
					queCount = 0
				}
			}
		}
		timeNotification.Stop()
	}(requestCh)

	hakaru := HakaruHandler{
		DB:        db,
		requestCH: requestCh,
	}

	http.Handle("/hakaru", hakaru)
	http.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })

	// start server
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}

type HakaruHandler struct {
	DB        *sql.DB
	requestCH chan *http.Request
}

func (h HakaruHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.requestCH <- r

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
