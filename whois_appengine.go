// +build appengine

package ipjogli

import (
	"net"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/socket"
)

func newWhoisConn(r *http.Request) (net.Conn, error) {
	ctx := appengine.NewContext(r)
	conn, err := socket.DialTimeout(ctx, "tcp", "whois.cymru.com:43", 3*time.Second)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
