package proxy

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Timeout int               `json:"timeout"`
}

func Handler(token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if token != "" {
			auth := r.Header.Get("Authorization")
			if auth != "Bearer "+token {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		if req.URL == "" {
			http.Error(w, `{"error":"url is required"}`, http.StatusBadRequest)
			return
		}

		if req.Method == "" {
			req.Method = "GET"
		}

		timeout := time.Duration(req.Timeout) * time.Second
		if timeout <= 0 {
			timeout = 30 * time.Second
		}

		var body io.Reader
		if req.Body != "" {
			body = strings.NewReader(req.Body)
		}

		upstream, err := http.NewRequest(req.Method, req.URL, body)
		if err != nil {
			http.Error(w, `{"error":"failed to create upstream request"}`, http.StatusBadRequest)
			return
		}

		for k, v := range req.Headers {
			upstream.Header.Set(k, v)
		}

		start := time.Now()
		client := &http.Client{Timeout: timeout}
		resp, err := client.Do(upstream)
		duration := time.Since(start)

		if err != nil {
			log.Printf("[PROXY] %s %s -> ERROR %s (%s)", req.Method, req.URL, err.Error(), duration)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		log.Printf("[PROXY] %s %s -> %d (%s)", req.Method, req.URL, resp.StatusCode, duration)

		// Copy upstream response headers
		for k, vals := range resp.Header {
			for _, v := range vals {
				w.Header().Add(k, v)
			}
		}

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}
