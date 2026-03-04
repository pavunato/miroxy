package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/viptony/miroxy/proxy"
)

func main() {
	port := flag.Int("port", 0, "Port to listen on (default 8080)")
	token := flag.String("token", "", "Bearer token for auth (overrides MIROXY_TOKEN)")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)
	if *port == 0 {
		addr = os.Getenv("MIROXY_ADDR")
		if addr == "" {
			addr = ":8080"
		}
	}

	authToken := *token
	if authToken == "" {
		authToken = os.Getenv("MIROXY_TOKEN")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("POST /proxy", proxy.Handler(authToken))

	log.Printf("miroxy listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
