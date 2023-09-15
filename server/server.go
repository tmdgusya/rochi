package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// echoUpper reads lines from r, uppercases them, and writes them to w.
func echoUpper(w io.Writer, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		// note that scanner.Text() strips the newline character from the end of the line,
		// so we need to add it back in when we write to w.
		fmt.Fprintf(w, "%s\n", strings.ToUpper(line))
	}
	if err := scanner.Err(); err != nil {
		log.Printf("error: %s", err)
	}
}

func main() {
	const name = "rochi"
	log.SetPrefix("name" + "\t")

	port := flag.Int("p", 8080, "port to listen on")
	flag.Parse()

	// ListenTCP creates a TCP listener accepting connections on the given address
	// TCPAddr represents the address of a TCP end point; it has an IP, Port, and Zone, all of which are optional.
	// Zone only matters for IPv6; we'll ignore it for now
	// If we omit the IP, it means we are listening on all available IP address
	// If we omit the Port, it means we are listening on a random port.
	// We want to listen on a port specified by the user on the command-line.
	// see https://golang.org/pkg/net/#ListenTCP and https://golang.org/pkg/net/#Dial for details.
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: *port})

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	log.Printf("listening at localhost: %s", listener.Addr())

	for {
		// loop forever, acceptiong connections one at a time
		// we need to use go routine for acceptiong many connections at a time
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go echoUpper(conn, conn)
	}
}
