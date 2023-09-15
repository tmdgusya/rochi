package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	const name = "writetcp"
	log.SetPrefix(name + "\t")

	// register the command-line flages: -p specifies the port to connect to
	port := flag.Int("p", 8080, "port to connect to")
	flag.Parse()

	// connect to a server at an IP address and port
	// bidirectional TCP connection
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: *port})

	if err != nil {
		log.Fatalf("Error connecting to localhost:%d: %v", *port, err)
	}

	log.Printf("Connected to %s: will forward stdin", conn.RemoteAddr())

	defer conn.Close()

	// spawn a goroutine to read incoming lines from the server and print them to stdout
	go func() {

		for connScanner := bufio.NewScanner(conn); connScanner.Scan(); {
			fmt.Printf("%s\n", connScanner.Text())

			if err := connScanner.Err(); err != nil {
				log.Fatalf("Error reading from %s: %v", conn.RemoteAddr(), err)
			}

			if connScanner.Err() != nil {
				log.Fatalf("Error reading from %s: %v", conn.RemoteAddr(), err)
			}
		}
	}()

	// read incoming lines from stdin and forware the to the server
	for stdinScanner := bufio.NewScanner(os.Stdin); stdinScanner.Scan(); {
		log.Printf("sent: %s\n", stdinScanner.Text())

		// scanner.Bytes() returns a slice of bytes up to but not including the next newline
		if _, err := conn.Write(stdinScanner.Bytes()); err != nil {
			log.Fatalf("Error writing to %s: %v", conn.RemoteAddr(), err)
		}

		// we need to add the newline back in
		if _, err := conn.Write([]byte("\n")); err != nil {
			log.Fatalf("Error writing to %s: %v", conn.RemoteAddr(), err)
		}

		if stdinScanner.Err() != nil {
			log.Fatalf("Error reading from %s: %v", conn.RemoteAddr(), err)
		}
	}
}
