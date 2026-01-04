package main

import (
	"net/http"
	"strconv"
	"time"
)

func labelHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	remoteAddr := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		remoteAddr = forwarded
	}

	logInfo("%s %s from %s", r.Method, r.URL.Path, remoteAddr)

	if r.Method != http.MethodGet {
		logError("method not allowed: %s (expected GET)", r.Method)
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	params, err := parseLabelParams(r.URL.Query())
	if err != nil {
		logError("parameter parsing failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logDebug("params: size=%dx%d dpi=%.1f margin=%d padding=%d qrSize=%d title=%q secondary=%q id=%q url=%q",
		params.width, params.height, params.dpi, params.margin, params.padding,
		params.qrSize, params.titleText, params.secondaryText, params.idText, shortURLFrom(params.url))

	img, err := renderLabel(params)
	if err != nil {
		logError("rendering failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	maxUpload := envInt("HBOX_WEB_MAX_UPLOAD_SIZE", defaultMaxUpload)
	pngData, err := encodePNGWithDPI(img, params.dpi)
	if err != nil {
		logError("PNG encoding failed: %v", err)
		http.Error(w, "failed to encode image", http.StatusInternalServerError)
		return
	}

	if len(pngData) > maxUpload {
		logError("image size %d bytes exceeds maximum %d bytes", len(pngData), maxUpload)
		http.Error(w, "image exceeds maximum upload size", http.StatusRequestEntityTooLarge)
		return
	}

	duration := time.Since(startTime)
	logInfo("generated %dx%d PNG (%d bytes, %.1f DPI) in %v",
		params.width, params.height, len(pngData), params.dpi, duration)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(pngData)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pngData)
}
