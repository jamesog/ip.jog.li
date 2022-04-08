package ipjogli

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/netip"
	"os"
	"strings"

	"github.com/jamesog/iptoasn"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/info", infoHandler)
}

// remoteAddr reads an HTTP request and returns an IP (in string form).
func remoteAddr(r *http.Request) string {
	if app := os.Getenv("FLY_APP_NAME"); app != "" {
		return r.Header.Get("Fly-Client-IP")
	}
	return parseXFF(r)
}

func parseXFF(r *http.Request) string {
	// Attempt to get the IP from the X-Forwarded-For header.
	// If that's empty then use the request's remote address.
	// N.B. if there are multiple X-Forwarded-For headers (meaning the user
	// went through multiple proxies) we use the last globally-routable
	// address.
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		addrs := strings.Split(xff, ",")
		// Walk backwards through the XFF addresses.
		// Typically a user of this library wants to know what they appear as
		// to a remote server so we don't want the first egress IP, rather
		// the last possible proxy that isn't internal to the network of this
		// server.
		// Note: If there is a reverse proxy in front of the this application
		// this will return the wrong IP. A list of "trusted" proxies is
		// needed.
		for i := len(addrs) - 1; i >= 0; i-- {
			ip, err := netip.ParseAddr(strings.TrimSpace(addrs[i]))
			if err != nil {
				// X-Forwarded-For is corrupted
				return ""
			}
			// If this is a globally-routable address, assume it's the last
			// known ingress and therefore how the request is seen by HTTP
			// servers.
			if ip.IsGlobalUnicast() && isRoutableAddr(ip) {
				return ip.String()
			}
		}
	}

	// Parse the HTTP connection address
	addr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// Assume that this is App Engine and thus there was no port
		addr = r.RemoteAddr
	}
	return addr
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Attempt to get the IP from the X-Forwarded-For header.
	// If that's empty then use the request's remote address.
	addr := remoteAddr(r)
	fmt.Fprintf(w, addr+"\n")
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	addr := remoteAddr(r)

	// Attempt to fetch from memcache
	info, err := iptoasn.LookupIP(addr)
	if err != nil {
		log.Printf("lookup error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	i, err := json.Marshal(info)
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", i)
}
