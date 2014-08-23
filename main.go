package main

import (
	"flag"
	"log"

	"github.com/ajsd/godirs/whitelist"
	"github.com/go-martini/martini"
)

var (
	addrFlag          = flag.String("addr", "", "Address to use. [host]:port.")
	whitelistFileFlag = flag.String("cors-whitelist-file", "", "CORS whitelisted origins file (one origin per line).")
)

var m *martini.Martini

func init() {
	m = martini.New()
	m.Use(martini.Recovery())
	m.Use(martini.Logger())

	r := martini.NewRouter()
	r.Get(dirsPath, ListFiles)

	m.Action(r.Handle)
}

func initWhitelist() {
	if *whitelistFileFlag != "" {
		w, err := whitelist.NewFromFile(*whitelistFileFlag)
		if err != nil {
			log.Fatalf("Error loading whitelist file: %v\n", err)
		}
		m.Use(w.ServeHTTP)
	} else {
		log.Printf("No CORS whitelist specificied (-cors-whitelist-file). Cross-domain requests will have default behaviour")
	}
}

func main() {
	log.SetPrefix("[godirs] ")
	flag.Parse()
	if *addrFlag == "" {
		log.Fatalln("-addr is required")
	}
	initWhitelist()
	m.RunOnAddr(*addrFlag)
}
