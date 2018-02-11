package ipjogli

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

// This fakes the Team Cymru WHOIS service
func testWhoisServer(t *testing.T) {
	l, err := net.Listen("tcp", whoisAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	info := bytes.NewBufferString("NA      | 192.0.2.1        | NA                  |    | other    |            | NA\n")

	for {
		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}
		go func(c net.Conn) {
			c.Write(info.Bytes())
			c.Close()
		}(conn)
	}
}

func TestIP(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	expected := "192.0.2.1\n"

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, resp.StatusCode)
	}

	if string(body) != expected {
		t.Errorf("expected body %q, got %q", expected, string(body))
	}
}

func TestInfo(t *testing.T) {
	whoisAddr = "[::1]:10043"

	go testWhoisServer(t)

	req := httptest.NewRequest("GET", "/info", nil)
	w := httptest.NewRecorder()
	infoHandler(w, req)

	info := ipInfo{
		IP:        "192.0.2.1",
		BGPPrefix: "NA",
		AS:        0,
		ASName:    "NA",
		Country:   "",
		Registry:  "other",
	}
	want, err := json.Marshal(info)
	if err != nil {
		t.Fatal(err)
	}

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("wanted status %v, got %v", http.StatusOK, resp.StatusCode)
	}

	if string(body) != string(want) {
		t.Errorf("wanted body %q, got %q", string(want), string(body))
	}
}
