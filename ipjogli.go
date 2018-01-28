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
	addr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, addr+"\n")
}
