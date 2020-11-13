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

func endPoint(path, ext string) func(http.ResponseWriter, *http.Request) {
	var (
		json []byte
		err  error
	)
	if !cache {
		json, err = ioutil.ReadFile(path)
		if err != nil {
			json = []byte(`{"Error": "ReadFile(` + path + `)"}`)
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if cache {
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

func HandleStub(port int) {
	fmt.Printf("server start on:  http://localhost:%d\n\nend points:\n", port)

	for _, path := range walkStab("./stub") {
		url, ext := path2url(path, "stub")
		http.HandleFunc(url+"/", endPoint(path, ext))
		http.HandleFunc(url, endPoint(path, ext))
		fmt.Printf("  http://localhost:%d%s\n", port, url)
	}
	fmt.Print("\n")
}

//

func HandleProxy(proxy string) {
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			log.Printf("  bad proxy url: %v\n", err)
			return
		}
		reverseProxy := httputil.NewSingleHostReverseProxy(proxyUrl)

		fmt.Printf("\nreverse on:  %s\n\nproxy url:\n", proxy)

		for _, path := range walkStab("./proxy") {
			url, _ := path2url(path, "proxy")

			http.Handle(url+"/", reverseProxy)
			http.Handle(url, reverseProxy)

			fmt.Printf("  %s%s\n", proxy, url)
		}
		fmt.Print("\n\n")
	}
}

//

func walkStab(path string) []string {
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
