package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	app    string
	rqlog  bool
	reread bool
)

func main() {
	var (
		js    = flag.Bool("react", false, "write setupProxy.js file (for React CRA with http-proxy-middleware)")
		port  = flag.Int("p", 8080, "server port")
		proxy string
	)
	flag.BoolVar(&reread, "reread", false, "disable endpoint file cache (reread file on every request)")
	flag.BoolVar(&rqlog, "log", false, "show HTTP request log")
	flag.StringVar(&app, "root", "./build", "root path for index.html and ./static foder")
	flag.StringVar(&proxy, "proxy", "", `reverse proxy url (example:  http://localhost:9000 )`)
	flag.Parse()

	if *js {
		setupProxy_js()
		return
	}
	if _, err := os.Stat(app); os.IsNotExist(err) {
		app = "../build"
	}

	//

	HandleStub(*port)
	HandleProxy(*port, proxy)
	http.HandleFunc("/static/", staticServer)
	http.HandleFunc("/", index)

	var hdl http.Handler
	if rqlog {
		hdl = logMiddleware(http.DefaultServeMux)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), hdl))
}

//

func index(wr http.ResponseWriter, req *http.Request) {
	var index = app + "/index.html"
	if _, err := os.Stat(index); os.IsNotExist(err) || req.URL.Path != "/" {
		err404(wr, req.URL.Path)
		return
	}
	http.ServeFile(wr, req, index)
}

func staticServer(wr http.ResponseWriter, req *http.Request) {
	path := app + req.URL.Path
	if file, err := os.Stat(path); err == nil && !file.IsDir() {
		http.ServeFile(wr, req, path)
		return
	}
	err404(wr, path)
}

func err404(w http.ResponseWriter, error string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(404)
	fmt.Fprintf(w, "<!DOCTYPE html><html style=\"font-family: sans-serif;\">UnServ: %s<h2>404 page not found</h2></html>", error)
}

func logMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("    %s: %s\n", r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
