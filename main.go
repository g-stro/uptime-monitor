package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const interval = 5 // in minutes

var record Record
var config configuration

type configuration struct {
	Url  string
	Port int
}

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
	contents, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(contents, &config)
	if err != nil {
		log.Fatal(err)
	}

	record = Record{
		Address:    config.Url,
		Last30Days: make([]DayRecord, 1, 30),
	}
}

func main() {
	getStatus(&record)
	go loop(&record)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}
