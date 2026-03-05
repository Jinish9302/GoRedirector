package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"fmt"

	"github.com/joho/godotenv"
)

func matchHost(host string, patterns []string) bool {
	host = strings.ToLower(host)

	for _, p := range patterns {
		p = strings.TrimSpace(strings.ToLower(p))

		if strings.HasPrefix(p, "*.") {
			suffix := p[1:]
			if strings.HasSuffix(host, suffix) {
				return true
			}
		} else if host == p {
			return true
		}
	}

	return false
}

func main() {

	// load .env if exists
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, continuing with system env")
	}

	target := os.Getenv("REDIRECT_URL")
	hostEnv := os.Getenv("ALLOWED_HOSTS")

	if target == "" {
		log.Fatal("REDIRECT_URL must be set")
	}

	patterns := strings.Split(hostEnv, ",")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		host := r.Host
		fmt.Println(host)
		if strings.Contains(host, ":") {
			host = strings.Split(host, ":")[0]
		}

		if matchHost(host, patterns) {
			http.Redirect(w, r, target+r.RequestURI, http.StatusTemporaryRedirect)
			return
		}

		http.NotFound(w, r)
	})

	log.Println("Listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}