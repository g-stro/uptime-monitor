package main

import (
	"encoding/json"
	"fmt"
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
	// http server config
	Url      string
	HttpPort int
	// smtp config
	SmptHost   string
	SmtpPort   int
	Sender     string
	SenderPass string
	Recipients []string
}

func getStatus(r *Record) {
	resp, err := http.Get(r.Address)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != 200 {
		// send mail to recipients
		msg := fmt.Sprintf("Subject: Monitoring - %v\r\nUh-Oh!\r\n\r\n"+
			"%v responded with status code %v when performing monitoring.",
			config.Url, config.Url, resp.StatusCode)
		SendMail(msg)
		log.Printf("Uh-Oh! %v appears to be down. Status code: %v\n", r.Address, resp.StatusCode)
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
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.HttpPort), nil))
}
