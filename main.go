package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	httpPort string
	httpAddr string
)

func init() {
	flag.StringVar(&httpPort, "http.port", "8080", "-http.port=8080")
	flag.StringVar(&httpAddr, "http.addr", "", "-http.addr=localhost")
}

func main() {
	log.Fatal(http.ListenAndServe(httpAddr+":"+httpPort, nil))
}

func budgetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		return
	}
	if r.Method == "POST" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}
