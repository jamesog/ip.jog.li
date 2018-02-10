package ipjogli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/info", infoHandler)
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

type ipInfo struct {
	IP        string `json:"ip"`
	BGPPrefix string `json:"bgp_prefix"`
	AS        int    `json:"as"`
	ASName    string `json:"as_name"`
	Country   string `json:"country"`
	Registry  string `json:"registry"`
}

func whoisIP(ip string) (ipInfo, error) {
	d := net.Dialer{
		// Don't wait too long to establish the connection
		Timeout: 3 * time.Second,
		// Enable Happy Eyeballs
		DualStack: true,
	}

	conn, err := d.Dial("tcp", "whois.cymru.com:43")
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

	// Attempt to fetch from memcache
	var info ipInfo
	var err error
	info, err = mcGet(r, addr)
	if err != nil {
		// Memcache miss, so fetch from Team Cymru and save to memcache
		info, err = whoisIP(addr)
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
