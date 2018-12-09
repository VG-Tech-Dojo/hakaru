package main

import (
	"container/list"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
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

func RunDB(db *sql.DB, eventlogStack list.List) {
	fmt.Println("INSERT NUM", eventlogStack.Len())
	query := "INSERT INTO eventlog(at, name, value) values"
	for i := 0; i < eventlogStack.Len(); i++ {
		query += "(?, ?, ?)"
		if i != eventlogStack.Len()-1 {
			query += ","
		}
	}
	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	s := make([]interface{}, 0)
	for i := 0; i < eventlogStack.Len(); i++ {
		event := eventlogStack.Remove(eventlogStack.Front())
		eventLog := event.(EventLog)
		s = append(s, eventLog.Now)
		s = append(s, eventLog.Name)
		s = append(s, eventLog.Value)
	}
	fmt.Println("args Num", len(s), "event num", eventlogStack.Len())
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
		eventlogStack := list.New()
		cpyEveneLogStack := list.New()
		_ = cpyEveneLogStack
		mux := new(sync.Mutex)
		for {
			mux.Lock()
			select {
			case r := <-requestCh:
				name := r.URL.Query().Get("name")
				value := r.URL.Query().Get("value")

				eventlogStack.PushBack(NewEventLog(name, value))
				mux.Unlock()
			case <-timeNotification.C:
				if eventlogStack.Len() != 0 {
					cpyEveneLogStack = eventlogStack
					go RunDB(db, *eventlogStack)
					eventlogStack = list.New()
					mux.Unlock()
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
