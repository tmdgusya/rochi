package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	const method = "GET"
	const url = "https://eblog.fly.dev/index.html"
	var body io.Reader = nil
	req, err := http.NewRequestWithContext(context.TODO(), method, url, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Accept-Encoding", "deflate")
	req.Header.Set("User-Agent", "eblog/1.0")
	req.Header.Set("some-key", "a value")
	req.Header.Set("SOMe-KEY", "somevalue") // will overwrite the above and be canonicalized since we used Set rather than Add
	req.Write(os.Stdout)
}
