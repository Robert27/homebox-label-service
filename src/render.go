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
	if params.width <= 0 || params.height <= 0 {
		return nil, errors.New("invalid label size")
	}

	img := image.NewRGBA(image.Rect(0, 0, params.width, params.height))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	innerWidth := params.width - 2*params.margin
	innerHeight := params.height - 2*params.margin
	if innerWidth < 1 || innerHeight < 1 {
		return nil, errors.New("invalid label size")
	}

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
	}
	leftColX := params.margin
	rightColX := params.margin + leftColWidth + colGap
	if singleColumn {
		rightColX = leftColX
	}

	// Header text should use left column width to leave space for icon
	headerWidth := leftColWidth
	headerX := params.margin
	cursorY := params.margin
	titleLines := wrapTwoLines(params.titleText, headerWidth, titleDrawer)
	if len(titleLines) > 0 {
		drawTextLines(titleDrawer, titleLines, headerX, cursorY, headerWidth, alignLeft)
		cursorY += textBlockHeight(titleDrawer.Face, len(titleLines))
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

	if len(titleLines) > 0 || secondaryText != "" {
		cursorY += params.padding
	}

	contentTop := cursorY
	qrY := contentTop

	qr, err := qrcode.New(params.url, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	availableHeight := params.height - params.margin - qrY
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
		qrImg := qr.Image(qrSize)
		qrX := leftColX
		qrRect := image.Rect(qrX, qrY, qrX+qrSize, qrY+qrSize)
		draw.Draw(img, qrRect, qrImg, image.Point{}, draw.Src)
	}

	// Show ID label with extracted ID in bottom right
	idText := strings.TrimSpace(params.idText)
	idBlockHeight := 0
	if idText != "" {
		idLabelDrawer := &font.Drawer{Dst: img, Src: image.Black, Face: idLabelFace}
		idValueDrawer := &font.Drawer{Dst: img, Src: image.Black, Face: idValueFace}
		// Only truncate if text doesn't fit
		if idValueDrawer.MeasureString(idText).Ceil() > rightColWidth {
			idText = truncateWithEllipsis(idText, rightColWidth, idValueDrawer)
		}
		idGap := maxInt(2, params.padding/2)
		idLabelHeight := textBlockHeight(idLabelFace, 1)
		idValueHeight := textBlockHeight(idValueFace, 1)
		idBlockHeight = idLabelHeight + idGap + idValueHeight
		idTop := params.height - params.margin - idBlockHeight
		drawTextLines(idLabelDrawer, []string{"ID"}, rightColX, idTop, rightColWidth, alignRight)
		drawTextLines(idValueDrawer, []string{idText}, rightColX, idTop+idLabelHeight+idGap, rightColWidth, alignRight)
	}

	iconAreaTop := contentTop
	iconAreaBottom := params.height - params.margin
	if idBlockHeight > 0 {
		iconAreaBottom = params.height - params.margin - idBlockHeight - params.padding
	}
	iconAreaHeight := iconAreaBottom - iconAreaTop
	if iconAreaHeight > 0 && rightColWidth > 0 {
		iconSize := minInt(rightColWidth, iconAreaHeight)
		iconSize = int(float64(iconSize) * 0.65)
		if iconSize >= 12 {
			iconX := rightColX + (rightColWidth-iconSize)/2
			iconY := iconAreaTop + (iconAreaHeight-iconSize)/2
			drawOpenBoxIcon(img, iconX, iconY, iconSize, iconSize)
		}
	}

	return img, nil
}
