package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	target := query.Get("target")
	if len(query["target"]) != 1 || target == "" {
		http.Error(w, "'target' parameter must be specified once", 400)
		return
	}
	registry := prometheus.NewRegistry()
	collector := StatsCollector{ctx: r.Context(), donor: target}
	registry.MustRegister(collector)
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/metrics", handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head>
              <title>Folding@Home User Stats Exporter</title>
              <style>
                form label, form input { margin: 10px; }
              </style>
            </head>
            <body>
              <h1>Folding@Home User Stats Exporter</h1>
              <form action="/metrics">
                <label>User name or ID:</label> <input type="text" name="target"><br>
                <input type="submit" value="Submit">
              </form>
            </body>
            </html>`))
	})

	if err := http.ListenAndServe(":9702", nil); err != nil {
		fmt.Printf("Error starting HTTP server: %v", err)
		os.Exit(1)
	}
}
