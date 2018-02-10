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

func handler(w http.ResponseWriter, r *http.Request) {
	addr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// Assume that this is App Engine and thus there was no port
		addr = r.RemoteAddr
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
	addr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// Assume that this is App Engine and thus there was no port
		addr = r.RemoteAddr
	}

	// Attempt to fetch from memcache
	var info ipInfo
	info, err = mcGet(r, addr)
	if err != nil {
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
