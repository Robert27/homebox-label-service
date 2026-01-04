# HomeBox Label Service

HTTP service that renders a single label image for Homebox Label Maker. It exposes a GET endpoint that accepts query parameters and returns a PNG label image.

## Requirements

- Go 1.25+
- No external state, no auth

## Run

```sh
go run ./src
```

The service listens on `:8080` by default.

## Environment Variables

- `PORT`: HTTP port (default `8080`)
- `HBOX_LABEL_MAKER_LABEL_SERVICE_TIMEOUT`: request timeout in seconds or Go duration (default `30s`)
- `HBOX_WEB_MAX_UPLOAD_SIZE`: max response size in bytes (default `10485760`)
- `HBOX_LABEL_MAKER_LABEL_SERVICE_URL`: set this in Homebox to the service URL
- `LOG_LEVEL`: logging verbosity - `INFO` (default) or `DEBUG` for detailed logs

## Endpoint

`GET /`

Headers:
- `User-Agent: Homebox-LabelMaker/1.0`
- `Accept: image/*`

Response:
- `200 OK`
- `Content-Type: image/png`
- Body: PNG binary

## Query Parameters

Unused parameters are ignored safely.

- `Width` (int): label width in pixels
- `Height` (int): label height in pixels
- `Dpi` (float): rendering DPI
- `Margin` (int): outer margin in pixels
- `ComponentPadding` (int): padding between components in pixels
- `QrSize` (int): QR code size in pixels
- `URL` (string): URL to encode into the QR code
- `TitleText` (string): primary label text
- `TitleFontSize` (float): font size for title text
- `DescriptionText` (string): secondary text (also used for domain display)
- `DescriptionFontSize` (float): font size for secondary text
- `AdditionalInformation` (string): optional ID value (fallback keys: `ID`, `Id`)
- `DynamicLength` (bool): accepted but ignored

## Layout

Landscape label with:
- Top-left: bold title
- Under title: secondary URL/domain
- Bottom-left: QR code for `URL`
- Right side: open-box icon centered vertically
- Bottom-right: ID block ("ID" + value)

## Example

```sh
curl -o label.png "http://localhost:8080/?Width=320&Height=240&Dpi=203&Margin=8&ComponentPadding=6&QrSize=140&URL=https%3A%2F%2Finv.eggl.one%2Fitem%2F000-029&TitleText=Zahnstange&AdditionalInformation=inv.eggl.one"
```
