package h264

// #cgo CFLAGS: -I${SRCDIR}/../vendor/include
// #cgo CXXFLAGS: -I${SRCDIR}/../vendor/include
// #cgo LDFLAGS: ${SRCDIR}/../vendor/lib/openh264/libopenh264.a
// #include <string.h>
// #include <openh264/codec_api.h>
// #include <errno.h>
// #include "bridge.hpp"
import "C"

import (
	"fmt"
	"unsafe"
	"image"
	"github.com/pion/codec"
)

type H264Options struct {
	Width        int
	Height       int
	Bitrate      int
	MaxFrameRate float32
}

// https://github.com/cisco/openh264/wiki/TypesAndStructures#sencparambase
func (h *H264Options) translate() C.SEncParamBase {
	return C.SEncParamBase{
		iUsageType: 	C.CAMERA_VIDEO_REAL_TIME,
		iRCMode: 		C.RC_BITRATE_MODE,
		iPicWidth:      C.int(h.Width),
		iPicHeight:     C.int(h.Height),
		iTargetBitrate: C.int(h.Bitrate),
		fMaxFrameRate:  C.float(h.MaxFrameRate),
	}
}

type h264Encoder struct {
	encoder *C.Encoder
}

var _ codec.Encoder = &h264Encoder{}


func NewEncoder(opts H264Options) (codec.Encoder, error) {
	encoder, err := C.enc_new(opts.translate())
	if (err != nil) {
		// TODO: better error message
		return nil, fmt.Errorf("failed in creating encoder")
	}

	e := h264Encoder{
		encoder: encoder,
	}
	return &e, nil
}

func (e *h264Encoder) Encode(img image.Image) ([]byte, error) {
	// TODO: Conver img to yuv since openh264 only accepts yuv
	yuvImg := img.(*image.YCbCr)
	bounds := yuvImg.Bounds()
	s, err := C.enc_encode(e.encoder, C.Frame{
		y: unsafe.Pointer(&yuvImg.Y[0]),
		u: unsafe.Pointer(&yuvImg.Cb[0]),
		v: unsafe.Pointer(&yuvImg.Cr[0]),
		height: C.int(bounds.Max.Y - bounds.Min.Y),
		width: C.int(bounds.Max.X - bounds.Min.X),
	})
	if (err != nil) {
		// TODO: better error message
		return nil, fmt.Errorf("failed in encoding")
	}

	return C.GoBytes(unsafe.Pointer(s.data), s.data_len), nil
}

func (e *h264Encoder) Close() error {
	C.enc_free(e.encoder)
	return nil
}