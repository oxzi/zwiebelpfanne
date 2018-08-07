package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"

	"golang.org/x/net/proxy"
)

// State's variables, will be set via command line flags
var (
	torSocks      string
	localAddr     string
	hiddenService string
)

// hiddenServiceConn tries to create a net.Conn through Tor's SOCKS5 proxy
// to a Hidden Service. The torSocks-parameter should be the Tor's SOCKS5 proxy
// (like "localhost:9050") and the hiddenService-parameter should be your
// destination (like "foobar2323.onion:1234").
func hiddenServiceConn(torSocks, hiddenService string) (net.Conn, error) {
	dial, err := proxy.SOCKS5("tcp", torSocks, nil, nil)
	if err != nil {
		return nil, err
	}

	conn, err := dial.Dial("tcp", hiddenService)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// handleConn takes an incomming net.Conn from a listener and tries to bind it
// to a connection to a Hidden Service created by hiddenServiceConn.
func handleConn(conn net.Conn) {
	defer conn.Close()

	torConn, err := hiddenServiceConn(torSocks, hiddenService)
	if err != nil {
		fmt.Printf("Failed to establish Hidden Service connection for %s: %v\n",
			conn.RemoteAddr().String(), err.Error())
		return
	}
	defer torConn.Close()

	go io.Copy(torConn, conn)
	io.Copy(conn, torConn)
}

func main() {
	flag.StringVar(&torSocks, "tor-socks5", "localhost:9050", "Tor's SOCKS5 proxy")
	flag.StringVar(&localAddr, "listen", "localhost:1337", "Bind to this address")
	flag.StringVar(&hiddenService, "onion", "", "Hidden Service to be bound to")
	flag.Parse()

	if hiddenService == "" {
		fmt.Println("zwiebelpfanne: --onion is missing")
		os.Exit(1)
	}

	fmt.Printf("zwiebelpfanne: %s -> %s\n", hiddenService, localAddr)

	ln, err := net.Listen("tcp", localAddr)
	if err != nil {
		fmt.Printf("Could not listen on %s: %v\n", localAddr, err.Error())
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err.Error())
			continue
		}

		fmt.Printf("Established connection to %s\n", conn.RemoteAddr().String())
		go handleConn(conn)
	}
}
