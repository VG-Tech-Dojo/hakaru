package main

import (
	"database/sql"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initDB() *sql.DB {
	dataSourceName := "root:hakaru-pass@tcp(127.0.0.1:13307)/hakaru-db-test"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}

	db.Exec(`DROP TABLE eventlog`)
	_, err = db.Exec(`
			CREATE TABLE eventlog (
			  id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
			  at datetime DEFAULT NULL,
			  name varchar(255) NOT NULL,
			  value int(10) unsigned,
			  PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
		`)
	if err != nil {
		panic(err)
	}
	return db
}

func Test_Main(t *testing.T) {

	db := initDB()
	defer db.Close()

	h := HakaruHandler{DB: db}
	ts := httptest.NewServer(h)
	defer ts.Close()

	nameWant := "test"
	valueWant := "1"
	url := fmt.Sprintf("%s?name=%s&value=%s", ts.URL, nameWant, valueWant)
	fmt.Println(url)

	r, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	if r.StatusCode != http.StatusOK {
		t.Fatalf("got %d", r.StatusCode)
	}

	acaoGot := r.Header["Access-Control-Allow-Origin"]
	acaoWant := []string{"*"}
	if !cmp.Equal(acaoGot, acaoWant) {
		t.Fatalf("got %s. want %s", acaoGot, acaoWant)
	}

	acahGot := r.Header["Access-Control-Allow-Headers"]
	acahWant := []string{"Content-Type"}
	if !cmp.Equal(acahGot, acahWant) {
		t.Fatalf("got %s. want %s", acahGot, acahWant)
	}

	acamGot := r.Header["Access-Control-Allow-Methods"]
	acamWant := []string{"GET"}
	if !cmp.Equal(acamGot, acamWant) {
		t.Fatalf("got %s. want %s", acamGot, acamWant)
	}

	rows, err := db.Query(`SELECT name, value FROM eventlog`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Fatal("no rows")
	}
	var nameGot, valueGot string
	if err := rows.Scan(&nameGot, &valueGot); err != nil {
		t.Fatal(err)
	}

	if nameGot != nameWant {
		t.Fatalf("got: %s, want: %s", nameGot, nameWant)
	}

	if valueGot != valueWant {
		t.Fatalf("got: %s, want: %s", nameGot, nameWant)
	}

	if rows.Next() {
		t.Fatal("got: more than 1 rows.")
	}

}
