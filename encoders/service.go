package encoders

import (
	"image"
	"io"
)

type Service interface {
	NewEncoder(codec VideoCodec, size image.Point, frameRate int) (Encoder, error)
	Supports(codec VideoCodec) bool
}

type Encoder interface {
	io.Closer
	Encode(*image.RGBA) ([]byte, error)
	VideoSize() (image.Point, error)
}

type VideoCodec = int

const (
	NoCodec VideoCodec = iota
	H264Codec
	VP8Codec
)
