package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func makeEndPoint(path, ext string, reread bool) func(http.ResponseWriter, *http.Request) {
	var (
		json []byte
		err  error
	)
	if !reread {
		json, err = ioutil.ReadFile(path)
		if err != nil {
			json = []byte(`{"Error": "ReadFile(` + path + `)"}`)
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if reread {
			json, err = ioutil.ReadFile(path)
			if err != nil {
				json = []byte(`{"Error": "ReadFile(` + path + `)"}`)
			}
		}
		if ext == ".json" {
			w.Header().Set("Content-Type", "application/json")
		}
		w.Write(json)
	}
}

// make endpoint for each file in stub dir
func HandleStub(port int, reread bool) {
	fmt.Printf("server start on:  http://localhost:%d\n\nendpoints:\n", port)

	for _, path := range walkOnDir("./stub") {
		url, ext := path2url(path, "stub")

		endPoint := makeEndPoint(path, ext, reread)
		http.HandleFunc(url+"/", endPoint)
		http.HandleFunc(url, endPoint)

		fmt.Printf("  http://localhost:%d%s\n", port, url)
	}
	fmt.Print("\n")
}

//

func makeProxyPoint(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// https://bana.io/blog/go-newsinglehostreverseproxy/#main-go
		r.Host = r.URL.Host
		proxy.ServeHTTP(w, r)
	}
}

// make reverse proxy endpoint for each file in proxy dir
func HandleProxy(port int, proxy string) {
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			fmt.Println(" Error: wrong reverse proxy url:", err)
			return
		}
		reverseProxy := httputil.NewSingleHostReverseProxy(proxyUrl)

		fmt.Printf("\nreverse on:  %s\n\nproxy url:\n", proxy)

		for _, path := range walkOnDir("./proxy") {
			url, _ := path2url(path, "proxy")

			proxyPoint := makeProxyPoint(reverseProxy)
			http.HandleFunc(url+"/", proxyPoint)
			http.HandleFunc(url, proxyPoint)

			fmt.Printf("  http://localhost:%d%s =>\n", port, url)
			fmt.Printf("      %s%s\n", proxy, url)
		}
		fmt.Print("\n\n")
	}
}

//

func walkOnDir(path string) []string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("      required %s/\n", path)
		return nil
	}

	var out = []string{}
	var err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("   bad path%s\n", path)
			return err
		}
		if !info.IsDir() {
			out = append(out, path)
		}
		return nil
	})

	if err != nil {
		log.Printf("%v\n", err)
		return nil
	}
	return out
}

func path2url(path, prefix string) (url, ext string) {
	ext = filepath.Ext(path)
	url = strings.Replace(strings.TrimPrefix(strings.TrimSuffix(path, ext), prefix), "\\", "/", -1)

	if len(url) < 2 {
		log.Fatalf(" Error: wrong filename for endpoint: '%s'", path)
	}
	return
}
