package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	PORT      = ":8000"
	BASE_HOST = ""
)

func main() {
	http.HandleFunc("/", reverseProxyHandler)
	log.Printf("Reverse Proxy Running on %s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

func reverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	host := r.Host // e.g., "p1.localhost:8000"
	subdomain := strings.Split(host, ".")[0]

	// Construct target base like: https://name.s3.us-east-1.amazonaws.com
	target, err := url.Parse(BASE_HOST)
	if err != nil {
		log.Printf("[ERROR] Failed to parse base host: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Incoming Request - Host: %s | Subdomain: %s | Path: %s", host, subdomain, r.URL.Path)

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Capture default director
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// Prefix the path with __outputs/{subdomain}
		if req.URL.Path == "/" {
			req.URL.Path = "/__outputs/" + subdomain + "/index.html"
		} else {
			req.URL.Path = "/__outputs/" + subdomain + req.URL.Path
		}

		// Update host header
		req.Host = target.Host

		log.Printf("[INFO] Final Proxy Request → URL: %s%s", BASE_HOST, req.URL.Path)
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		log.Printf("[INFO] Response Code: %d for %s", resp.StatusCode, resp.Request.URL)
		return nil
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		log.Printf("[ERROR] Proxy error: %v", err)
		http.Error(w, "Proxy error", http.StatusBadGateway)
	}

	proxy.ServeHTTP(w, r)
}