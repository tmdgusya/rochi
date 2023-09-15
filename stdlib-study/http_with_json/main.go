package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// Request https://go.dev/play/p/A8QVJwFEeq3
type Request struct {
	Format string `json:"format"` // Format, as in time.Format. If empty, use time.RFC3339.
	TZ     string `json:"tz"`     // TZ, as in time.LoadLocation. If empty, use time.Local.
}

// Response The time, formatted according to the request's Format and TZ.
type Response struct {
	Time time.Time `json:"time"`
}

// Error no need for omitempty here; we'll never send a zero time.
type Error struct {
	Error string `json:"error"`
}

func getTime(w http.ResponseWriter, r *http.Request) {
	var req Request
	w.Header().Set("Content-Type", "encoding/json")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(Error{err.Error()})
		return
	}
	r.Body.Close()
	var tz *time.Location = time.Local
	if req.TZ != "" {
		var err error
		tz, err = time.LoadLocation(req.TZ)
		if err != nil || tz == nil {
			w.WriteHeader(400) // bad request
			json.NewEncoder(w).Encode(Error{err.Error()})
			return
		}
	}

	resp := Response{Time: time.Now().In(tz)}
	json.NewEncoder(w).Encode(resp)
}

var client = &http.Client{Timeout: 2 * time.Second}

func sendRequest(tz, format string) {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(Request{TZ: tz, Format: format})
	log.Printf("request body: %v", body)
	req, err := http.NewRequestWithContext(context.TODO(), "GET", "http://localhost:8080", body)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Write(os.Stdout)
	resp.Body.Close() // always close response bodies when you're done with them.
}

func main() {
	server := http.Server{Addr: ":8080", Handler: http.HandlerFunc(getTime)}
	go server.ListenAndServe()

	sendRequest("", "") // rely on defaults
	sendRequest("America/Los_Angeles", time.RFC3339)
	sendRequest("America/New_York", time.RFC822Z) // "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	sendRequest("faketz", "")                     // should get 400 Bad Request

}
