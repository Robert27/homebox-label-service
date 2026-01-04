package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"image"
	"image/png"
	"math"
)

var pngSignature = []byte{137, 80, 78, 71, 13, 10, 26, 10}

func encodePNGWithDPI(img image.Image, dpi float64) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	if dpi <= 0 {
		return buf.Bytes(), nil
	}
	return addPhysChunk(buf.Bytes(), dpi)
}

func addPhysChunk(pngData []byte, dpi float64) ([]byte, error) {
	if len(pngData) < len(pngSignature)+12 {
		return nil, errors.New("invalid png data")
	}
	if !bytes.Equal(pngData[:len(pngSignature)], pngSignature) {
		return nil, errors.New("invalid png signature")
	}

	offset := len(pngSignature)
	if len(pngData[offset:]) < 8 {
		return nil, errors.New("invalid png header")
	}
	ihdrLen := int(binary.BigEndian.Uint32(pngData[offset : offset+4]))
	chunkType := string(pngData[offset+4 : offset+8])
	if chunkType != "IHDR" {
		return nil, errors.New("missing IHDR chunk")
	}
	ihdrEnd := offset + 8 + ihdrLen + 4
	if ihdrEnd > len(pngData) {
		return nil, errors.New("invalid IHDR length")
	}

	ppm := int(math.Round(dpi / 0.0254))
	if ppm <= 0 {
		return pngData, nil
	}

	physData := make([]byte, 9)
	binary.BigEndian.PutUint32(physData[0:4], uint32(ppm))
	binary.BigEndian.PutUint32(physData[4:8], uint32(ppm))
	physData[8] = 1

	var chunk bytes.Buffer
	if err := binary.Write(&chunk, binary.BigEndian, uint32(len(physData))); err != nil {
		return nil, err
	}
	if _, err := chunk.Write([]byte("pHYs")); err != nil {
		return nil, err
	}
	if _, err := chunk.Write(physData); err != nil {
		return nil, err
	}
	crc := crc32.ChecksumIEEE(append([]byte("pHYs"), physData...))
	if err := binary.Write(&chunk, binary.BigEndian, crc); err != nil {
		return nil, err
	}

	var out bytes.Buffer
	if _, err := out.Write(pngData[:ihdrEnd]); err != nil {
		return nil, err
	}
	if _, err := out.Write(chunk.Bytes()); err != nil {
		return nil, err
	}
	if _, err := out.Write(pngData[ihdrEnd:]); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
