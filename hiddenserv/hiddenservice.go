// hiddenserv contains the HiddenServiceConn function which establishes a
// connection to a Tor Hidden Service through Tor's SOCKS5-proxy.
package hiddenserv

import (
	"net"

	"golang.org/x/net/proxy"
)

// HiddenServiceConn tries to create a net.Conn through Tor's SOCKS5 proxy
// to a Hidden Service. The torSocks-parameter should be the Tor's SOCKS5 proxy
// (like "localhost:9050") and the hiddenService-parameter should be your
// destination (like "foobar2323.onion:1234").
func HiddenServiceConn(torSocks, hiddenService string) (net.Conn, error) {
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
