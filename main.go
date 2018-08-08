package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	"github.com/geistesk/zwiebelpfanne/hiddenserv"
)

// handleConn takes an incomming net.Conn from a listener and tries to bind it
// to a connection to a Hidden Service created by hiddenserv.HiddenServiceConn.
func handleConn(conn net.Conn, torSocks, hiddenService string) {
	defer conn.Close()

	torConn, err := hiddenserv.HiddenServiceConn(torSocks, hiddenService)
	if err != nil {
		fmt.Printf("Failed to establish Hidden Service connection for %s: %v\n",
			conn.RemoteAddr(), err.Error())
		return
	}
	defer torConn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		io.Copy(torConn, conn)
		wg.Done()
	}()

	go func() {
		io.Copy(conn, torConn)
		wg.Done()
	}()

	wg.Wait()
	fmt.Printf("Close connection to %s\n", conn.RemoteAddr())
}

func main() {
	var localAddr, torSocks, hiddenService string

	flag.StringVar(&localAddr, "listen", "localhost:1337", "Bind to this address")
	flag.StringVar(&torSocks, "tor-socks5", "localhost:9050", "Tor's SOCKS5 proxy")
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
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err.Error())
			continue
		}

		fmt.Printf("Established connection to %s\n", conn.RemoteAddr())
		go handleConn(conn, torSocks, hiddenService)
	}
}
