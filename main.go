package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const (
	port         = "8080"
	filePathRoot = "."
)

func main() {

	dbg := flag.Bool("debug", false, "activate debug mode")

	flag.Parse()

	if *dbg {
		os.Remove("./database.json")
		fmt.Printf("deleting databse.json...\n")
	}

	// by default, godotenv will look for a file named .env in the current directory
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")

	// storing secret key inside apiConfig istance cfg
	cfg := apiConfig{
		jwtSecret: jwtSecret,
	}

	mux := http.NewServeMux()

	// Setting handler for "/app/*"
	fsRoot := http.FileServer(http.Dir(filePathRoot))
	mux.Handle("/app/*", cfg.middlewareMetricsInc(http.StripPrefix("/app", fsRoot)))

	// Setting handler for "POST /api/users"
	mux.HandleFunc("POST /api/users", handlerUsersPost)

	// Setting handler for "PUT /api/users"
	mux.HandleFunc("PUT /api/users", cfg.handlerUsersPut)

	// Setting handler for "POST /api/login"
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)

	// Setting handler for "POST /api/refresh"
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)

	// Setting handler for "POST /api/revoke"
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevokeRefreshToken)

	// Setting handler for "POST /api/chirps"
	mux.HandleFunc("POST /api/chirps", handlerChirpsPost)

	// Setting handler for "GET /api/chirps"
	mux.HandleFunc("GET /api/chirps", handlerChirpsGet)

	// Setting handler for "Get /api/chirps/{chirp_id}"
	mux.HandleFunc("GET /api/chirps/{chirp_id}", handlerChirpGetByID)

	// Setting handle for /healthz
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	// Setting handle for "/admin/metrics"
	mux.HandleFunc("GET /admin/metrics", cfg.handlerAdminMetrics)

	// Setting handle for "/reset"
	mux.HandleFunc("/api/reset", cfg.handlerReset)

	corsMux := middlewareCors(mux)

	ourServer := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	fmt.Printf("Listening and Serving %s/ on port: %s\n", filePathRoot, port)
	log.Fatal(ourServer.ListenAndServe())

}
