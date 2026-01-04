package main

import (
	"errors"
	"log"
	"net/http"
	"time"
)

func main() {
	port := envString("PORT", "8080")
	timeout := envDuration("HBOX_LABEL_MAKER_LABEL_SERVICE_TIMEOUT", 30*time.Second)

	mux := http.NewServeMux()
	mux.HandleFunc("/", labelHandler)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       timeout,
		WriteTimeout:      timeout,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("homebox label service listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
