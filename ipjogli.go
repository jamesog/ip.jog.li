package ipjogli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/info", infoHandler)
}

// remoteAddr reads an HTTP request and returns an IP (in string form).
func remoteAddr(r *http.Request) string {
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
		for i := len(addrs) - 1; i >= 0; i-- {
			ip := net.ParseIP(strings.TrimSpace(addrs[i]))
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

type ipInfo struct {
	IP        string `json:"ip"`
	BGPPrefix string `json:"bgp_prefix"`
	AS        int    `json:"as"`
	ASName    string `json:"as_name"`
	Country   string `json:"country"`
	Registry  string `json:"registry"`
}

var whoisAddr = "whois.cymru.com:43"

func whoisIP(r *http.Request, ip string) (ipInfo, error) {
	conn, err := newWhoisConn(r)
	if err != nil {
		return ipInfo{}, err
	}

	// -v enable all fields
	// -f disable header
	// Order matters!
	fmt.Fprintf(conn, "-v -f %s\n", ip)
	res, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return ipInfo{}, nil
	}
	conn.Close()

	fields := strings.Split(res, "|")
	for i, _ := range fields {
		fields[i] = strings.TrimSpace(fields[i])
	}

	// Convert the AS to an integer
	// The AS may be returned as "NA", e.g. for non-global IPs
	as, err := strconv.Atoi(fields[0])
	if err != nil && fields[0] != "NA" {
		return ipInfo{}, err
	}

	info := ipInfo{
		IP:        fields[1],
		BGPPrefix: fields[2],
		AS:        as,
		ASName:    fields[6],
		Country:   fields[3],
		Registry:  fields[4],
	}

	return info, nil
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	addr := remoteAddr(r)

	// Attempt to fetch from memcache
	var info ipInfo
	var err error
	if info, err = mcGet(r, addr); err != nil {
		// Memcache miss, so fetch from Team Cymru and save to memcache
		info, err = whoisIP(r, addr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		mcSet(r, info)
	}

	i, err := json.Marshal(info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", i)
}
