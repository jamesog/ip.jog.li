// +build !appengine

package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/jamesog/ip.jog.li"
)

func main() {
	s := &http.Server{
		Addr:         ":8000",
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
