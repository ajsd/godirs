package main

import (
	"bufio"
	"container/list"
	"errors"
	"flag"
	"log"
	"os"
	"strings"
)

const (
	noWhitelistError = "No whitelist specified for CORS. Cross-origin requests will have default behaviour."
)

// Global (see #init())
var whitelist *list.List

// Flags
var (
	whitelistFileFlag = flag.String("cors-whitelist", "", "File containing allowed CORS origins")
)

func IsWhitelisted(origin string) bool {
	if whitelist == nil {
		return false
	}
	for e := whitelist.Front(); e != nil; e = e.Next() {
		if e.Value == origin {
			return true
		}
	}
	return false
}

func readWhitelistFile(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	newlist := list.New()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") || strings.TrimSpace(line) == "" {
			continue
		}
		newlist.PushBack(line)
		log.Printf("[CORS] Whitelisted '%s'\n", line)
	}
	whitelist.Init()
	whitelist.PushBackList(newlist)
	return nil
}

func GetWhitelist() []string {
	if whitelist == nil {
		return nil
	}
	w := make([]string, whitelist.Len())
	for e := whitelist.Front(); e != nil; e = e.Next() {
		if origin, ok := e.Value.(string); ok && origin != "" {
			w = append(w, origin)
		}
	}
	return w
}

func ReloadWhitelist() error {
	if *whitelistFileFlag == "" {
		return errors.New(noWhitelistError)
	}
	return readWhitelistFile(*whitelistFileFlag)
}

func init() {
	flag.Parse()
	if *whitelistFileFlag == "" {
		log.Println(noWhitelistError)
		whitelist = nil
		return
	}
	whitelist = list.New()
	readWhitelistFile(*whitelistFileFlag)
}
