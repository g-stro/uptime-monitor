package main

import (
	"log"
	"net/http"
	"time"
)

const interval = 5 // in minutes

var record Record

func getStatus(r *Record) {
	resp, err := http.Get(r.Address)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != 200 {
		log.Printf("Uh-Oh! %v appears to be down. Response code: %v\n", r.Address, resp.StatusCode)
	}

	r.update(resp.StatusCode, interval)
}

func loop(r *Record) {
	ticker := time.NewTicker(interval * time.Second)
	for range ticker.C {
		getStatus(r)
	}
}

func init() {
	record = Record{
		Address:    "https://www.justice-defenders.org/",
		Last30Days: make([]DayRecord, 1, 30),
	}
}

func main() {
	getStatus(&record)
	go loop(&record)
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
