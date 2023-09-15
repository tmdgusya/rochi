package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// define flags
var (
	host, path, method string
	port               int
)

func main() {
	// initialize & parse flags
	flag.StringVar(&method, "method", "GET", "HTTP method to use")
	flag.StringVar(&host, "host", "localhost", "host to connect to")
	flag.StringVar(&path, "path", "/", "path to request")
	flag.IntVar(&port, "port", 8080, "port to connect to")
	flag.Parse()

	// ResolveTCP Addr is a slightly more convenient way of creating a TCPAddr
	ip, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}

	// dial(connect to) the remote host using the TCP addr we just created
	conn, err := net.DialTCP("tcp", nil, ip)
	if err != nil {
		panic(err)
	}

	log.Printf("Connected to %s (@ %s)", host, conn.RemoteAddr())

	defer conn.Close()

	var reqfields = []string{
		fmt.Sprintf("%s %s HTTP/1.1", method, path),
		"Host: " + host,
		"User-Agent: Roach",
		"",

		// body would go here, if we had one
	}

	request := strings.Join(reqfields, "\r\n") + "\r\n"

	if _, err = conn.Write([]byte(request)); err != nil {
		log.Printf("Error sending request to remote server(%s): %v", ip, err)
	}

	log.Printf("sent request: \n%s", request)

	for scanner := bufio.NewScanner(conn); scanner.Scan(); {
		line := scanner.Bytes()
		if _, err := fmt.Fprintf(os.Stdout, "%s\n", line); err != nil {
			log.Printf("Error writing to connection: %s", err)
		}
		if scanner.Err() != nil {
			log.Printf("Error reading from connection: %s", err)
			return
		}
	}
}
