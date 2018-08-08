package hiddenserv

import (
	"io"
	"strings"
	"testing"
)

const torSocks = "localhost:9050"

// TestHiddenServiceConn tests the HiddenServiceConn function by establishing
// a connection to Facebook's Hidden Service. This test requires the Tor daemon
// to be running on localhost:9050 and Facebook to not stop serving their
// "social network" through Tor. So, it relies on some side effects.
func TestHiddenServiceConn(t *testing.T) {
	conn, err := HiddenServiceConn(torSocks, "facebookcorewwwi.onion:80")
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	_, err = conn.Write(
		[]byte("GET / HTTP/1.0\r\nHost: facebookcorewwwi.onion\r\n\r\n"))
	if err != nil {
		t.Error(err)
	}

	buf := make([]byte, 256)
	for {
		len, err := conn.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			t.Error(err)
		}

		if strings.Contains(string(buf[:len]), "https://facebookcorewwwi.onion/") {
			return
		}
	}

	t.Error("Facebook's redirect wasn't in the response")
}
