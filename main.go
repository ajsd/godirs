package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
)

var (
	addrFlag = flag.String("addr", "", "Address to use. [host]:port.")
	baseFlag = flag.String("base", "", "Absolute path to the base directory.")
)

type FileInfo struct {
	IsDir bool   `json:"is_dir"`
	Name  string `json:"name"`
	Path  string `json:"path"`
	Size  int64  `json:"size"`
}

func skip(name string) bool {
	return strings.HasPrefix(name, ".")
}

func handler(w http.ResponseWriter, r *http.Request) {
	p := path.Join(*baseFlag, r.URL.Path)
	if !strings.HasPrefix(p, *baseFlag) {
		http.Error(w, "Bad path", http.StatusForbidden)
		return
	}
	infos, err := ioutil.ReadDir(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var ret []FileInfo
	for _, info := range infos {
		if skip(info.Name()) {
			continue
		}
		finfo := FileInfo{
			IsDir: info.IsDir(),
			Name:  info.Name(),
			Path:  path.Join(r.URL.Path, info.Name()),
			Size:  info.Size(),
		}
		ret = append(ret, finfo)
	}
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
}

func main() {
	flag.Parse()
	if *addrFlag == "" || *baseFlag == "" {
		log.Fatalln("Both -addr and -base are required flags")
	}
	if !strings.HasPrefix(*baseFlag, "/") {
		log.Fatalln("-base must be absolute")
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(*addrFlag, nil)
}
