package main

import (
	"image"
	"image/color"
	"image/draw"
)

func drawOpenBoxIcon(img *image.RGBA, x, y, w, h int) {
	if w <= 0 || h <= 0 {
		return
	}
	thickness := w / 14
	if thickness < 2 {
		thickness = 2
	}
	if thickness > 5 {
		thickness = 5
	}

	frontLeftX := x + int(float64(w)*0.18)
	frontRightX := x + int(float64(w)*0.82)
	frontTopY := y + int(float64(h)*0.45)
	frontBottomY := y + int(float64(h)*0.85)

	flapLeftX := x + int(float64(w)*0.05)
	flapRightX := x + int(float64(w)*0.95)
	flapSideY := y + int(float64(h)*0.3)
	flapPeakX := x + w/2
	flapPeakY := y + int(float64(h)*0.1)

	drawLine(img, frontLeftX, frontTopY, frontRightX, frontTopY, thickness)
	drawLine(img, frontRightX, frontTopY, frontRightX, frontBottomY, thickness)
	drawLine(img, frontRightX, frontBottomY, frontLeftX, frontBottomY, thickness)
	drawLine(img, frontLeftX, frontBottomY, frontLeftX, frontTopY, thickness)

	drawLine(img, frontLeftX, frontTopY, flapLeftX, flapSideY, thickness)
	drawLine(img, flapLeftX, flapSideY, flapPeakX, flapPeakY, thickness)
	drawLine(img, flapPeakX, flapPeakY, flapRightX, flapSideY, thickness)
	drawLine(img, flapRightX, flapSideY, frontRightX, frontTopY, thickness)

	drawLine(img, flapPeakX, flapPeakY, flapPeakX, frontTopY, thickness)
}

func drawLine(img *image.RGBA, x0, y0, x1, y1, thickness int) {
	dx := absInt(x1 - x0)
	dy := -absInt(y1 - y0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy

	for {
		drawThickPoint(img, x0, y0, thickness)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func drawThickPoint(img *image.RGBA, x, y, thickness int) {
	half := thickness / 2
	fillRect(img, x-half, y-half, thickness, thickness, color.Black)
}

func fillRect(img *image.RGBA, x, y, w, h int, c color.Color) {
	if w <= 0 || h <= 0 {
		return
	}
	rect := image.Rect(x, y, x+w, y+h).Intersect(img.Bounds())
	if rect.Empty() {
		return
	}
	draw.Draw(img, rect, &image.Uniform{C: c}, image.Point{}, draw.Src)
}
