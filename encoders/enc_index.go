package encoders

import (
	"fmt"
	"image"
)

type encoderFactory = func(size image.Point, frameRate int) (Encoder, error)

var registeredEncoders = make(map[VideoCodec]encoderFactory, 2)

type EncoderService struct {
}

func NewEncoderService() Service {
	return &EncoderService{}
}

func (*EncoderService) NewEncoder(codec VideoCodec, size image.Point, frameRate int) (Encoder, error) {
	factory, found := registeredEncoders[codec]
	if !found {
		return nil, fmt.Errorf("Codec not supported")
	}
	return factory(size, frameRate)
}

func (*EncoderService) Supports(codec VideoCodec) bool {
	_, found := registeredEncoders[codec]
	return found
}
