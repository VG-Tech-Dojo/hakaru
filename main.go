package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
	"log"
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
	eventChan = make(chan EventLog, 2000)
)

func timerInsert() {
	var db *sql.DB
	ticker := time.NewTicker(1 * time.Second)
	querys := make([]EventLog, 0, 2000)
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
			querys = make([]EventLog, 0, 2000) // 初期化
		case event := <-eventChan:
			querys = append(querys, event)
		}
	}
}

func hakaruHandler(ctx *fasthttp.RequestCtx) {
	var event EventLog
	event.Name = string(ctx.QueryArgs().Peek("name"))
	event.Value = string(ctx.URI().QueryArgs().Peek("value"))
	event.Now = time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	eventChan <- event

	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET")
}

func main() {
	go timerInsert()
	router := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/ok":
			ctx.SetStatusCode(200)
		case "/hakaru":
			hakaruHandler(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

	// start server
	if err := fasthttp.ListenAndServe(":8081", router); err != nil {
		log.Fatal(err)
	}
}
