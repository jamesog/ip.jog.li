// +build !appengine

package ipjogli

import (
	"net"
	"net/http"
	"time"
)

func newWhoisConn(r *http.Request) (net.Conn, error) {
	d := net.Dialer{
		// Don't wait too long to establish the connection
		Timeout: 3 * time.Second,
		// Enable Happy Eyeballs
		DualStack: true,
	}

	conn, err := d.Dial("tcp", "whois.cymru.com:43")
	if err != nil {
		return nil, err
	}

	return conn, nil
}
