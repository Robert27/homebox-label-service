package main

const (
	defaultWidth         = 320
	defaultHeight        = 240
	defaultDPI           = 203.0
	defaultMargin        = 8
	defaultPadding       = 6
	defaultQRSize        = 170
	defaultTitleFontSize = 28.0
	defaultDescFontSize  = 16.0
	defaultMaxUpload     = 10 * 1024 * 1024
)

type labelParams struct {
	width               int
	height              int
	dpi                 float64
	margin              int
	padding             int
	qrSize              int
	url                 string
	titleText           string
	secondaryText       string
	idText              string
	titleFontSize       float64
	descriptionFontSize float64
}
