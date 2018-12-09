package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	kinesis "github.com/aws/aws-sdk-go/service/kinesis"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	sess := session.Must(session.NewSession())
	//cred := credentials.NewSharedCredentials("", "sunrise2018")

	//auth := kinesis.New(sess, &aws.Config{Credentials: cred, Region: aws.String("ap-northeast-1")})
	auth := kinesis.New(sess, &aws.Config{Region: aws.String("ap-northeast-1")})
	streamName := "hakaru-stream"

	hakaruHandler := func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		value := r.URL.Query().Get("value")
		str := name + ":" + value
		id := strconv.Itoa(rand.Int())

		record := &kinesis.PutRecordInput{
			Data:         []byte(str),
			PartitionKey: &id,
			StreamName:   &streamName,
		}

		putsOutput, err := auth.PutRecord(record)

		if err != nil {
			panic(err)
		}

		fmt.Printf("%v\n", putsOutput)

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
