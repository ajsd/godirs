package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
)

var (
	addrFlag          = flag.String("addr", "", "Address to use. [host]:port.")
	whitelistFileFlag = flag.String("cors-whitelist-file", "", "CORS whitelisted origins file (one origin per line).")
)
var whitelist []string

var m *martini.Martini

func init() {
	StringListVar(&whitelist, "cors-whitelist", []string{}, "CORS whitelisted origins (comma-separated)")

	m = martini.New()
	m.Use(martini.Recovery())
	m.Use(martini.Logger())

	r := martini.NewRouter()
	r.Get(dirsPath, ListFiles)

	m.Action(r.Handle)
}

func initWhitelist() {
	if *whitelistFileFlag == "" && len(whitelist) == 0 {
		log.Printf("No CORS whitelist specificied (-cors-whitelist, -cors-whitelist-file). Cross-domain requests will have default behaviour")
		return
	}
	var origins []string
	if *whitelistFileFlag != "" {
		file, err := os.Open(*whitelistFileFlag)
		if err != nil {
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") || strings.TrimSpace(line) == "" {
				continue
			}
			origins = append(origins, line)
		}
	}
	for _, origin := range whitelist {
		origins = append(origins, origin)
	}
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     origins,
		AllowCredentials: true,
	}))
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
