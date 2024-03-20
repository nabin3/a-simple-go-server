package main

import (
	"html/template"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getfileserverHitsNumber() int {
	return cfg.fileserverHits
}

func (cfg *apiConfig) resetfileserverHitsNumber() {
	cfg.fileserverHits = 0
}

// Defining handler for "/admin/metrics"
func (cfg *apiConfig) handlerAdminMetrics(w http.ResponseWriter, r *http.Request) {
	adminHtml, err := template.ParseFiles("./html_files/admin_metrics.html")

	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error occured in parsing /admin/metrics related html template file\n")
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err1 := adminHtml.Execute(w, cfg.getfileserverHitsNumber())
	if err1 != nil {
		w.WriteHeader(500)
		log.Printf("error occured in injecting the metrics-counter in /admi/metrics related html template file\n")
	}
}

// Defining handler for "/reset"
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.resetfileserverHitsNumber()
	w.Write([]byte("fileserverHit counter reseted successfully to 0 \n"))
}
