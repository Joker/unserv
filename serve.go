package main

import (
	"io/fs"
	"net/http"
	"os"
	"syscall"
)

// http.HandleFunc("/", index)
func index(wr http.ResponseWriter, req *http.Request) {
	var index = app + "/index.html"
	if file, err := os.Stat(index); os.IsNotExist(err) || file.IsDir() || req.URL.Path != "/" {
		err404(wr, req.URL.Path)
		return
	}
	http.ServeFile(wr, req, index)
}

// http.HandleFunc("/static/", staticServer)
func staticServer(wr http.ResponseWriter, req *http.Request) {
	path := app + req.URL.Path
	if file, err := os.Stat(path); err == nil && !file.IsDir() {
		http.ServeFile(wr, req, path)
		return
	}
	err404(wr, path)
}

// http.Handle("/", http.FileServer(filterFileSystem{http.Dir(app)}))
type filterFileSystem struct{ http.FileSystem }

func (fsys filterFileSystem) Open(name string) (http.File, error) {

	file, err := fsys.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}

	if s, _ := file.Stat(); s.IsDir() && name != "/" {
		return nil, &fs.PathError{Op: "open", Path: name, Err: syscall.ENOENT}
	}

	return file, err
}
