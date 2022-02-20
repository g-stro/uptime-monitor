package main

import (
	"time"
)

type Record struct {
	Address    string
	Last30Days []DayRecord
	pingCount  int
	respCount  int
	Uptime     float32
	NextCheck  time.Time
}

type DayRecord struct {
	Responses []HttpResponse
	Uptime    float32
	pingCount int
}

type HttpResponse struct {
	DateTime time.Time
	status   int
}

func (r *Record) update(status int, interval time.Duration) {
	var todaysRecord DayRecord
	// check whether it's a new day or not
	if !time.Now().Truncate(24 * time.Hour).Equal(time.Now().Add(time.Duration(-interval) * time.Minute).Truncate(24 * time.Hour)) {
		if len(r.Last30Days) >= 30 {
			r.respCount = r.respCount - len(r.Last30Days[0].Responses)   // update response counter
			r.pingCount = r.pingCount - r.Last30Days[0].pingCount        // update ping counter
			r.Last30Days = append(r.Last30Days[:0], r.Last30Days[1:]...) // remove day 1
		}
		r.Last30Days = append(r.Last30Days, todaysRecord) // append empty record for now
	} else {
		todaysRecord = r.Last30Days[len(r.Last30Days)-1] // get existing record
	}

	todaysRecord.pingCount++
	r.pingCount++

	if status != 200 {
		resp := HttpResponse{
			DateTime: time.Now(),
			status:   status,
		}
		todaysRecord.Responses = append(todaysRecord.Responses, resp)
		r.respCount++
	}

	r.Uptime = float32(r.pingCount-r.respCount) / float32(r.pingCount) * 100
	todaysRecord.Uptime = float32(todaysRecord.pingCount-len(todaysRecord.Responses)) / float32(todaysRecord.pingCount) * 100
	r.Last30Days[len(r.Last30Days)-1] = todaysRecord
	r.NextCheck = time.Now().Add(time.Minute * interval)
}
