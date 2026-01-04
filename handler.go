package main

import (
	"bytes"
	"image/png"
	"net/http"
	"strconv"
)

func labelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	params, err := parseLabelParams(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	img, err := renderLabel(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		http.Error(w, "failed to encode image", http.StatusInternalServerError)
		return
	}

	maxUpload := envInt("HBOX_WEB_MAX_UPLOAD_SIZE", defaultMaxUpload)
	if buf.Len() > maxUpload {
		http.Error(w, "image exceeds maximum upload size", http.StatusRequestEntityTooLarge)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}
