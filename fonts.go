package main

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type textAlign int

const (
	alignLeft textAlign = iota
	alignCenter
	alignRight
)

func newFontFace(ttf []byte, pixelSize, dpi float64) (font.Face, error) {
	if dpi <= 0 {
		dpi = defaultDPI
	}
	points := pixelSize * 72.0 / dpi
	ft, err := opentype.Parse(ttf)
	if err != nil {
		return nil, err
	}
	return opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    points,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

func drawTextLines(drawer *font.Drawer, lines []string, x, topY, maxWidth int, align textAlign) int {
	if len(lines) == 0 {
		return 0
	}
	if maxWidth < 1 {
		maxWidth = 1
	}
	metrics := drawer.Face.Metrics()
	lineHeight := metrics.Height.Ceil()
	ascent := metrics.Ascent.Ceil()
	descent := metrics.Descent.Ceil()
	dstBounds := drawer.Dst.Bounds()

	for i, line := range lines {
		lineWidth := drawer.MeasureString(line).Ceil()
		alignedX := x
		switch align {
		case alignCenter:
			alignedX = x + (maxWidth-lineWidth)/2
		case alignRight:
			alignedX = x + maxWidth - lineWidth
		}
		y := topY + ascent + i*lineHeight
		if y+descent > dstBounds.Max.Y {
			continue
		}
		drawer.Dot = fixed.P(alignedX, y)
		drawer.DrawString(line)
	}

	return lineHeight * len(lines)
}

func textBlockHeight(face font.Face, lines int) int {
	if lines <= 0 {
		return 0
	}
	return face.Metrics().Height.Ceil() * lines
}
