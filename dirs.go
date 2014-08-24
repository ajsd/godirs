package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/go-martini/martini"
)

const (
	dirsPath = "/p/**"
)

var (
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

func ListFiles(params martini.Params) (int, string) {
	p := params["_1"]
	dir := path.Join(*baseFlag, p)
	if !strings.HasPrefix(dir, *baseFlag) {
		return http.StatusForbidden, "Bad path"
	}
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	var ret []*FileInfo
	for _, info := range infos {
		if skip(info.Name()) {
			continue
		}
		finfo := &FileInfo{
			IsDir: info.IsDir(),
			Name:  info.Name(),
			Path:  path.Join("/", p, info.Name()),
			Size:  info.Size(),
		}
		ret = append(ret, finfo)
	}
	data, err := json.Marshal(ret)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, string(data)
}

func init() {
	flag.Parse()
	if *baseFlag == "" {
		log.Fatalln("-base is required")
	}
	if !strings.HasPrefix(*baseFlag, "/") {
		log.Fatalln("-base must be absolute [%s]\n", *baseFlag)
	}
}
