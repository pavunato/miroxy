package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/viptony/miroxy/proxy"
)

var version = "dev"

func main() {
	port := flag.Int("port", 0, "Port to listen on (default 8080)")
	token := flag.String("token", "", "Bearer token for auth (overrides MIROXY_TOKEN)")
	showVersion := flag.Bool("version", false, "Print version and exit")
	update := flag.Bool("update", false, "Self-update to latest release")
	flag.Parse()

	if *showVersion {
		fmt.Println("miroxy " + version)
		return
	}

	if *update {
		fmt.Println("Updating miroxy...")
		cmd := exec.Command("sh", "-c", "curl -fsSL https://pavunato.github.io/miroxy/install.sh | MIROXY_SKIP_SURVEY=1 sh")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal("Update failed: ", err)
		}
		return
	}

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
