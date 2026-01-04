package main

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"strings"

	"github.com/skip2/go-qrcode"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
)

func renderLabel(params labelParams) (image.Image, error) {
	logDebug("starting label rendering: %dx%d", params.width, params.height)

	if params.width <= 0 || params.height <= 0 {
		logError("invalid label size: %dx%d", params.width, params.height)
		return nil, errors.New("invalid label size")
	}

	img := image.NewRGBA(image.Rect(0, 0, params.width, params.height))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	innerWidth := params.width - 2*params.margin
	innerHeight := params.height - 2*params.margin
	if innerWidth < 1 || innerHeight < 1 {
		logError("invalid inner size after margins: %dx%d (margin=%d)", innerWidth, innerHeight, params.margin)
		return nil, errors.New("invalid label size")
	}

	logDebug("inner dimensions: %dx%d (margins: %d)", innerWidth, innerHeight, params.margin)

	titleFace, err := newFontFace(gobold.TTF, params.titleFontSize, params.dpi)
	if err != nil {
		return nil, err
	}
	descFace, err := newFontFace(goregular.TTF, params.descriptionFontSize, params.dpi)
	if err != nil {
		return nil, err
	}
	idLabelSize := maxFloat(params.descriptionFontSize*0.85, 11.0)
	idValueSize := maxFloat(params.descriptionFontSize*1.4, params.descriptionFontSize+4.0)
	idLabelFace, err := newFontFace(goregular.TTF, idLabelSize, params.dpi)
	if err != nil {
		return nil, err
	}
	idValueFace, err := newFontFace(gobold.TTF, idValueSize, params.dpi)
	if err != nil {
		return nil, err
	}

	titleDrawer := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: titleFace,
	}
	descDrawer := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: descFace,
	}

	// Calculate column layout first to determine header width
	colGap := maxInt(params.padding, 4)
	if innerWidth-colGap < 2 {
		colGap = 0
	}
	minRightWidth := maxInt(80, int(float64(innerWidth)*0.32))
	leftTarget := maxInt(params.qrSize, int(float64(innerWidth)*0.6))
	leftColWidth := leftTarget
	if innerWidth-colGap-leftColWidth < minRightWidth {
		leftColWidth = innerWidth - colGap - minRightWidth
	}
	rightColWidth := innerWidth - colGap - leftColWidth
	singleColumn := false
	if leftColWidth < 1 || rightColWidth < 1 {
		leftColWidth = innerWidth
		rightColWidth = innerWidth
		colGap = 0
		singleColumn = true
		logDebug("using single column layout (left=%d, right=%d)", leftColWidth, rightColWidth)
	} else {
		logDebug("using two column layout (left=%d, right=%d, gap=%d)", leftColWidth, rightColWidth, colGap)
	}
	leftColX := params.margin
	rightColX := params.margin + leftColWidth + colGap
	if singleColumn {
		rightColX = leftColX
	}

	// Title uses full width and stays on one line to avoid shrinking QR space.
	headerWidth := innerWidth
	headerX := params.margin
	cursorY := params.margin
	titleBottom := cursorY
	titleText := strings.TrimSpace(params.titleText)
	if titleText != "" {
		if titleDrawer.MeasureString(titleText).Ceil() > headerWidth {
			titleText = truncateWithEllipsis(titleText, headerWidth, titleDrawer)
		}
		drawTextLines(titleDrawer, []string{titleText}, headerX, cursorY, headerWidth, alignLeft)
		cursorY += textBlockHeight(titleDrawer.Face, 1)
		titleBottom = cursorY
	}

	// Show AdditionalInformation/DescriptionText between title and QR code
	secondaryText := strings.TrimSpace(params.secondaryText)
	if secondaryText != "" {
		headerGap := maxInt(4, params.padding/2)
		if cursorY > params.margin {
			cursorY += headerGap
		}
		// Only truncate if text doesn't fit
		if descDrawer.MeasureString(secondaryText).Ceil() > leftColWidth {
			secondaryText = truncateWithEllipsis(secondaryText, leftColWidth, descDrawer)
		}
		drawTextLines(descDrawer, []string{secondaryText}, leftColX, cursorY, leftColWidth, alignLeft)
		cursorY += textBlockHeight(descDrawer.Face, 1)
	}

	if titleText != "" || secondaryText != "" {
		cursorY += params.padding
	}

	contentTop := cursorY

	qr, err := qrcode.New(params.url, qrcode.Medium)
	if err != nil {
		logError("QR code creation failed: %v", err)
		return nil, err
	}
	qr.DisableBorder = true

	availableHeight := params.height - params.margin - contentTop
	if availableHeight < 1 {
		availableHeight = 1
	}
	qrSize := params.qrSize
	if qrSize <= 0 {
		qrSize = minInt(leftColWidth, availableHeight)
	}
	qrSize = minInt(qrSize, leftColWidth)
	qrSize = minInt(qrSize, availableHeight)
	if qrSize > 0 {
		logDebug("rendering QR code: %dx%d at (%d,%d)", qrSize, qrSize, leftColX, params.height-params.margin-qrSize)
		qrImg := qr.Image(qrSize)
		qrX := leftColX
		qrY := params.height - params.margin - qrSize
		qrRect := image.Rect(qrX, qrY, qrX+qrSize, qrY+qrSize)
		draw.Draw(img, qrRect, qrImg, image.Point{}, draw.Src)
	} else {
		logDebug("skipping QR code (size would be 0)")
	}

	// Show ID label with extracted ID in bottom right
	idText := strings.TrimSpace(params.idText)
	idBlockHeight := 0
	if idText != "" {
		idLabelDrawer := &font.Drawer{Dst: img, Src: image.Black, Face: idLabelFace}
		idValueDrawer := &font.Drawer{Dst: img, Src: image.Black, Face: idValueFace}
		idGap := maxInt(2, params.padding/2)
		idLabelHeight := textBlockHeight(idLabelFace, 1)
		idValueHeight := textBlockHeight(idValueFace, 1)
		idBlockHeight = idLabelHeight + idGap + idValueHeight
		idTop := params.height - params.margin - idBlockHeight
		drawTextLines(idLabelDrawer, []string{"ID"}, rightColX, idTop, rightColWidth, alignRight)
		drawTextLines(idValueDrawer, []string{idText}, rightColX, idTop+idLabelHeight+idGap, rightColWidth, alignRight)
	}

	iconAreaTop := contentTop
	if !singleColumn {
		iconAreaTop = titleBottom
	}
	iconAreaBottom := params.height - params.margin
	if idBlockHeight > 0 {
		iconAreaBottom = params.height - params.margin - idBlockHeight - params.padding
	}
	iconAreaHeight := iconAreaBottom - iconAreaTop
	if iconAreaHeight > 0 && rightColWidth > 0 {
		iconSize := minInt(rightColWidth, iconAreaHeight)
		if iconSize >= 12 {
			iconX := rightColX + (rightColWidth-iconSize)/2
			iconY := iconAreaTop + (iconAreaHeight-iconSize)/2
			logDebug("rendering icon: %dx%d at (%d,%d)", iconSize, iconSize, iconX, iconY)
			drawOpenBoxIcon(img, iconX, iconY, iconSize, iconSize)
		} else {
			logDebug("skipping icon (size %d < minimum 12)", iconSize)
		}
	} else {
		logDebug("skipping icon (no available space: height=%d, width=%d)", iconAreaHeight, rightColWidth)
	}

	logDebug("label rendering completed successfully")
	return img, nil
}
