package ipjogli

import (
	"fmt"
	"net"
	"net/http"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var addr string
	if _, ok := r.Header["X-Appengine-Remote-Addr"]; ok {
		addr = r.RemoteAddr
	} else {
		var err error
		addr, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	fmt.Fprintf(w, addr+"\n")
}
