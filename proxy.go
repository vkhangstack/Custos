package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func StartProxy() {
	target, _ := url.Parse("http://127.0.0.1:8080") // backend
	proxy := httputil.NewSingleHostReverseProxy(target)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Proxy:", r.Method, r.URL.String())
		proxy.ServeHTTP(w, r)
	})

	log.Println("HTTP proxy on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
