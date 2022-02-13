package main

import (
	"log"
	"net/http"
	"time"
)

type record struct {
	address   string
	responses []httpResponse
}

type httpResponse struct {
	dateTime time.Time
	status   int
}

func (r *record) update(status int) {
	resp := httpResponse{
		dateTime: time.Now(),
		status:   status,
	}
	r.responses = append(r.responses, resp)
}

func getStatus(r *record) {
	resp, err := http.Get(r.address)

	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != 200 {
		log.Printf("Uh-Oh!\nwww.justice-defenders.org appears to be down.\n%v", resp.StatusCode)
	} else {
		log.Printf("Success!\nwww.justice-defenders.org is live!\n%v", resp.StatusCode)
	}

	r.update(resp.StatusCode)
}

func main() {
	record := record{
		address: "https://www.justice-defenders.org/",
	}
	getStatus(&record)

	ticker := time.NewTicker(5 * time.Second)

	for _ = range ticker.C {
		getStatus(&record)
	}
}
