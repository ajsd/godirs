package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	addrFlag = flag.String("addr", "", "Address to use. [host]:port.")
)

func reloadHandler(w http.ResponseWriter, r *http.Request) {
	ReloadWhitelist()
}

func handler(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if IsWhitelisted(origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	ListFilesHandler(w, r)
}

func init() {
	http.HandleFunc("/rerere", reloadHandler)
	http.HandleFunc(dirsPath, handler)
}

func main() {
	flag.Parse()
	if *addrFlag == "" {
		log.Fatalln("-addr is required")
	}
	http.ListenAndServe(*addrFlag, nil)
}
