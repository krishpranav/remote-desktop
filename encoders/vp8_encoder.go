//go:build vp8enc
// +build vp8enc

package encoders

import (
	"bytes"
	"fmt"
	"image"
	"unsafe"
)
import "C"

const keyFrameInterval = 10

type VP8Encoder struct {
	buffer     *bytes.Buffer
	realSize   image.Point
	codecCtx   C.vpx_codec_ctx_t
	vpxImage   C.vpx_image_t
	yuvBuffer  []byte
	frameCount uint
}

func newVP8Encoder(size image.Point, frameRate int) (Encoder, error) {
	buffer := bytes.NewBuffer(make([]byte, 0))

	var cfg C.vpx_codec_enc_cfg_t
	if C.codec_enc_config_default(&cfg) != 0 {
		return nil, fmt.Errorf("Can't init default enc. config")
	}
	cfg.g_w = C.uint(size.X)
	cfg.g_h = C.uint(size.Y)
	cfg.g_timebase.num = 1
	cfg.g_timebase.den = C.int(frameRate)
	cfg.rc_target_bitrate = 90000
	cfg.g_error_resilient = 1

	var vpxCodecCtx C.vpx_codec_ctx_t
	if C.codec_enc_init(&vpxCodecCtx, &cfg) != 0 {
		return nil, fmt.Errorf("Failed to initialize enc ctx")
	}
	var vpxImage C.vpx_image_t
	if C.vpx_img_alloc(&vpxImage, C.VPX_IMG_FMT_I420, C.uint(size.X), C.uint(size.Y), 0) == nil {
		return nil, fmt.Errorf("Can't alloc. vpx image")
	}

	return &VP8Encoder{
		buffer:     buffer,
		realSize:   size,
		codecCtx:   vpxCodecCtx,
		vpxImage:   vpxImage,
		yuvBuffer:  make([]byte, size.X*size.Y*2),
		frameCount: 0,
	}, nil
}

func (e *VP8Encoder) Encode(frame *image.RGBA) ([]byte, error) {

	encodedData := unsafe.Pointer(nil)
	var flags C.int
	if e.frameCount%keyFrameInterval == 0 {
		flags |= C.VPX_EFLAG_FORCE_KF
	}
	frameSize := C.encode_frame(
		&e.codecCtx,
		&e.vpxImage,
		C.int(e.frameCount),
		flags,
		unsafe.Pointer(&frame.Pix[0]),
		unsafe.Pointer(&e.yuvBuffer[0]),
		C.int(e.realSize.X),
		C.int(e.realSize.Y),
		&encodedData,
	)
	e.frameCount++
	if int(frameSize) > 0 {
		encoded := C.GoBytes(encodedData, frameSize)
		return encoded, nil
		return nil, nil
	}
	return nil, nil
}

func (e *VP8Encoder) VideoSize() (image.Point, error) {
	return e.realSize, nil
}

func (e *VP8Encoder) Close() error {
	C.vpx_img_free(&e.vpxImage)
	C.vpx_codec_destroy(&e.codecCtx)
	return nil
}

func init() {
	registeredEncoders[VP8Codec] = newVP8Encoder
}
