package main

import (
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

	maxUpload := envInt("HBOX_WEB_MAX_UPLOAD_SIZE", defaultMaxUpload)
	pngData, err := encodePNGWithDPI(img, params.dpi)
	if err != nil {
		http.Error(w, "failed to encode image", http.StatusInternalServerError)
		return
	}
	if len(pngData) > maxUpload {
		http.Error(w, "image exceeds maximum upload size", http.StatusRequestEntityTooLarge)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(pngData)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pngData)
}
