package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	const method = "GET"
	v := make(url.Values)
	v.Add("q", `"of Emrakul"`) // note we use go's raw string syntax (`) to avoid having to escape the double quotes.
	v.Add("order", "released")
	v.Add("dir", "asc")
	const path = "https://scryfall.com/search"
	dst := path + "?" + v.Encode() // Encode() will escape the values for us. Remember the '?' separator!
	req, err := http.NewRequestWithContext(context.TODO(), method, dst, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Write(os.Stdout)
}
