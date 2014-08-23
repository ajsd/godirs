package whitelist

import (
	"bufio"
	"container/list"
	"log"
	"net/http"
	"os"
	"strings"
)

type Whitelist interface {
	http.Handler
	IsWhitelisted(origin string) bool
}

type whitelist struct {
	whitelist *list.List
}

func (w *whitelist) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if w.IsWhitelisted(origin) {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
	}
}

func (w *whitelist) IsWhitelisted(origin string) bool {
	if w.whitelist == nil {
		return false
	}
	for e := w.whitelist.Front(); e != nil; e = e.Next() {
		if e.Value == origin {
			return true
		}
	}
	return false
}

func (w *whitelist) set(origins []string) {
	if w.whitelist == nil {
		w.whitelist = list.New()
	}
	w.whitelist.Init()
	for _, origin := range origins {
		w.whitelist.PushBack(origin)
		log.Printf("[CORS] Whitelisted '%s'\n", origin)
	}
}

func (w *whitelist) loadFromFile(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	var origins []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") || strings.TrimSpace(line) == "" {
			continue
		}
		origins = append(origins, line)
	}
	w.set(origins)
	return nil
}

func New(origins []string) Whitelist {
	w := &whitelist{whitelist: list.New()}
	w.set(origins)
	return w
}

func NewFromFile(name string) (Whitelist, error) {
	w := &whitelist{whitelist: list.New()}
	if err := w.loadFromFile(name); err != nil {
		return nil, err
	}
	return w, nil
}

func init() {
	log.SetPrefix("[Whitelist]")
}
