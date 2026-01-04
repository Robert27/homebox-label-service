package main

import (
	"strings"

	"golang.org/x/image/font"
)

func wrapTwoLines(text string, maxWidth int, drawer *font.Drawer) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	if maxWidth < 1 {
		return nil
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	line1, idx, overflow := buildLine(words, maxWidth, drawer)
	if idx >= len(words) || overflow && idx == 1 && line1 != "" && strings.HasSuffix(line1, "...") {
		return []string{line1}
	}

	line2, idx2, overflow2 := buildLine(words[idx:], maxWidth, drawer)
	if line2 == "" {
		return []string{line1}
	}
	if idx2 < len(words[idx:]) || overflow2 {
		line2 = ensureEllipsis(line2, maxWidth, drawer)
	}

	return []string{line1, line2}
}

func buildLine(words []string, maxWidth int, drawer *font.Drawer) (string, int, bool) {
	var line string
	for i, word := range words {
		candidate := word
		if line != "" {
			candidate = line + " " + word
		}
		if drawer.MeasureString(candidate).Ceil() <= maxWidth {
			line = candidate
			continue
		}
		if line == "" {
			return truncateWithEllipsis(word, maxWidth, drawer), i + 1, true
		}
		return line, i, true
	}
	return line, len(words), false
}

func ensureEllipsis(line string, maxWidth int, drawer *font.Drawer) string {
	if strings.HasSuffix(line, "...") {
		if drawer.MeasureString(line).Ceil() <= maxWidth {
			return line
		}
		return truncateWithEllipsis(strings.TrimSuffix(line, "..."), maxWidth, drawer)
	}
	return truncateWithEllipsis(line, maxWidth, drawer)
}

func truncateWithEllipsis(text string, maxWidth int, drawer *font.Drawer) string {
	ellipsis := "..."
	if maxWidth < 1 {
		return ""
	}
	if maxWidth <= drawer.MeasureString(ellipsis).Ceil() {
		return ellipsis
	}
	if drawer.MeasureString(text).Ceil() <= maxWidth {
		return text
	}

	runes := []rune(text)
	for len(runes) > 0 {
		candidate := string(runes) + ellipsis
		if drawer.MeasureString(candidate).Ceil() <= maxWidth {
			return candidate
		}
		runes = runes[:len(runes)-1]
	}
	return ellipsis
}
