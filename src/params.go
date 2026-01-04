package main

import (
	"net/url"
	"strconv"
	"strings"
)

func parseLabelParams(values url.Values) (labelParams, error) {
	params := labelParams{
		width:               parseInt(values, "Width", defaultWidth),
		height:              parseInt(values, "Height", defaultHeight),
		dpi:                 parseFloat(values, "Dpi", defaultDPI),
		margin:              parseInt(values, "Margin", defaultMargin),
		padding:             parseInt(values, "ComponentPadding", defaultPadding),
		qrSize:              parseInt(values, "QrSize", defaultQRSize),
		url:                 queryGet(values, "URL"),
		titleText:           queryGet(values, "TitleText"),
		secondaryText: firstNonEmpty(
			queryGet(values, "DescriptionText"),
			queryGet(values, "AdditionalInformation"),
			queryGet(values, "AdditiontalInformation"),
		),
		idText:              firstNonEmpty(queryGet(values, "ID"), queryGet(values, "Id")),
		titleFontSize:       parseFloat(values, "TitleFontSize", defaultTitleFontSize),
		descriptionFontSize: parseFloat(values, "DescriptionFontSize", defaultDescFontSize),
	}

	if params.width <= 0 {
		params.width = defaultWidth
	}
	if params.height <= 0 {
		params.height = defaultHeight
	}
	if params.dpi <= 0 {
		params.dpi = defaultDPI
	}
	if params.margin < 0 {
		params.margin = 0
	}
	minDim := minInt(params.width, params.height)
	maxMargin := (minDim - 1) / 2
	if maxMargin < 0 {
		maxMargin = 0
	}
	if params.margin > maxMargin {
		params.margin = maxMargin
	}
	if params.padding < 0 {
		params.padding = 0
	}
	if params.titleFontSize <= 0 {
		params.titleFontSize = defaultTitleFontSize
	}
	if params.descriptionFontSize <= 0 {
		params.descriptionFontSize = defaultDescFontSize
	}

	if params.qrSize <= 0 {
		params.qrSize = defaultQRSize
	}
	maxQR := minInt(params.width-2*params.margin, params.height-2*params.margin)
	if maxQR < 1 {
		maxQR = 1
	}
	if params.qrSize > maxQR {
		logDebug("QR size %d clamped to maximum %d (label size: %dx%d, margin: %d)",
			params.qrSize, maxQR, params.width, params.height, params.margin)
		params.qrSize = maxQR
	}

	if params.titleText == "" && params.secondaryText != "" {
		params.titleText = params.secondaryText
		params.secondaryText = ""
	}
	// Always extract item ID from URL for bottom right (no "ID" label)
	// Use explicit ID/Id parameters if provided, otherwise extract from URL
	if params.idText == "" {
		extractedID := extractItemIDFromURL(params.url)
		if extractedID != "" {
			logDebug("extracted ID '%s' from URL", extractedID)
			params.idText = extractedID
		}
	}
	if params.url == "" {
		params.url = " "
	}

	return params, nil
}

func queryGet(values url.Values, key string) string {
	if value := values.Get(key); value != "" {
		return value
	}
	for k, v := range values {
		if strings.EqualFold(k, key) && len(v) > 0 {
			return v[0]
		}
	}
	return ""
}

func parseInt(values url.Values, key string, fallback int) int {
	value := queryGet(values, key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseFloat(values url.Values, key string, fallback float64) float64 {
	value := queryGet(values, key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}
	return parsed
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func shortURLFrom(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	host := parsed.Host
	if host == "" {
		parts := strings.Split(raw, "/")
		host = parts[0]
	}
	host = strings.TrimPrefix(host, "www.")
	return host
}

func extractItemIDFromURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	path := parsed.Path
	if path == "" {
		return ""
	}
	// Look for /item/ pattern and extract the ID after it
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "item" && i+1 < len(parts) {
			return strings.TrimSpace(parts[i+1])
		}
	}
	return ""
}
