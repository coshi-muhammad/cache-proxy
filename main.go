package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	cache_path     = "./cache/"
	max_cache_size = 500 // posibly use this to free up the cache if it gets full
)

var (
	//NOTE: comandline arguments and flags
	port        = flag.Int("port", 0, "The port to listen to")
	origin      = flag.String("origin", "", "the domain to proxy and cache")
	clear_cache = flag.Bool("clear-cache", false, "clears the cached responses. when this flag is present no server will be setup")
)

var (
	// NOTE: global state modify with caution
	cache = make(map[string]bool) //posibly add a mechanizme to free the cache
)

// NOTE: this is ai generated just so the broswers shuts up
func setupCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle browser preflight checks
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return true // Indicates it was a preflight request
	}

	return false
}

// TODO: handle the requests
func handleProxy(w http.ResponseWriter, r *http.Request) {
	//FIX: the css doesnt get loaded in the first go around for some reason
	if setupCORS(w, r) {
		return
	}

	url := *origin + r.URL.Path
	log.Println(url)
	_, ok := cache[url]
	path := strings.Replace(url, "/", "__", -1)
	path = strings.Replace(path, ":", "--", 1)
	ext := filepath.Ext(path)
	switch ext {
	case "css":
		w.Header().Set("Content-type", "text/css")
	case "js":
		w.Header().Set("Content-type", "text/javascript")
	case "ttf":
		w.Header().Set("Content-type", "font/ttf")
		// w.header.set("Content-type")
	}
	if ok {
		w.Header().Add("X-Cache", "HIT")
		http.ServeFile(w, r, cache_path+path)
	} else {
		res, err := http.Get(url)
		if err != nil {
			w.WriteHeader(res.StatusCode)
			w.Write([]byte("Couldnt resolve the requested endpoint 1"))
			return
		}
		defer res.Body.Close()
		buf, err := io.ReadAll(res.Body)
		if err != nil && err != io.EOF {
			log.Println(err)
			w.WriteHeader(res.StatusCode)
			w.Write([]byte("Couldnt resolve the requested endpoint 2"))
			return
		}
		err = os.WriteFile(cache_path+path, buf, 0755)
		if err != nil {
			log.Println("Failed to write the buffer to disk: ", err)
		} else {
			cache[url] = true
		}
		w.Header().Add("X-Cache", "MISS")
		w.Write(buf)
	}

}

func main() {
	//NOTE: this doesnt handle anything other than a simple url and it assumes all the cached data is html content
	flag.Parse()
	if *clear_cache {
		if err := os.RemoveAll(cache_path); err != nil {
			log.Fatal("failed to delete file: ", err)
		}
		//HACK: Im almost certain there is a better way to delete the content but keeping the directory intact
		if err := os.MkdirAll(cache_path, 0755); err != nil {
			log.Fatal("Failed to crate directory: ", err)
		}
	} else {
		if err := os.MkdirAll(cache_path, 0755); err != nil {
			log.Fatal("Failed to crate directory: ", err)
		}
		cached_urls, err := fs.Glob(os.DirFS(cache_path), "*")
		if err != nil {
			log.Fatal("Failed to get file list: ", err)
		}
		for _, url := range cached_urls {
			url = strings.Replace(url, "__", "/", -1)
			url = strings.Replace(url, "--", ":", 1)
			cache[url] = true
		}
		log.Println(cache)
		log.Println("listening on port :", *port)
		if err := http.ListenAndServe(fmt.Sprint(":", *port), http.HandlerFunc(handleProxy)); err != nil {
			log.Fatal("Failed to listen: ", err)
		}
	}
}
