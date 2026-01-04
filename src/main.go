package main

import (
	"errors"
	"log"
	"net/http"
	"time"
)

func main() {
	initLogLevel()
	port := envString("PORT", "8080")
	timeout := envDuration("HBOX_LABEL_MAKER_LABEL_SERVICE_TIMEOUT", 30*time.Second)
	maxUpload := envInt("HBOX_WEB_MAX_UPLOAD_SIZE", defaultMaxUpload)

	logInfo("HomeBox Label Service starting")
	logDebug("  port: %s", port)
	logDebug("  timeout: %v", timeout)
	logDebug("  max upload size: %d bytes", maxUpload)

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

	logInfo("HomeBox Label Service listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
