package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"strings"

	"github.com/blackjack/webcam"
	"github.com/pion/codec/h264"
)

func main() {
	height := 480
	width := 640
	enc, err := h264.NewEncoder(h264.H264Options{
		Height:       height,
		Width:        width,
		Bitrate:      2000000,
		MaxFrameRate: 60,
	})
	defer enc.Close()

	if err != nil {
		log.Fatal(err)
	}

	cam, err := webcam.Open("/dev/video0")
	if err != nil {
		panic(err)
	}

	var selectedFormat webcam.PixelFormat
	for v, k := range cam.GetSupportedFormats() {
		if strings.HasPrefix(k, "Motion-JPEG") {
			selectedFormat = v
		}
	}

	if selectedFormat == 0 {
		panic("Only Motion-JPEG supported")
	}

	if _, _, _, err = cam.SetImageFormat(selectedFormat, uint32(width), uint32(height)); err != nil {
		panic(err)
	}

	if err = cam.StartStreaming(); err != nil {
		panic(err)
	}

	for {
		err = cam.WaitForFrame(1)
		switch err.(type) {
		case nil:
		case *webcam.Timeout:
			fmt.Fprint(os.Stderr, err.Error())
			continue
		default:
			return
		}

		frame, err := cam.ReadFrame()
		if err != nil {
			panic(err)
		}
		if len(frame) != 0 {
			img, err := jpeg.Decode(bytes.NewReader(frame))
			if err != nil {
				continue
			}

			encoded, err := enc.Encode(img)
			if err != nil {
				panic(err)
			}

			os.Stdout.Write(encoded)
		}
	}
}
