package camera

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"strings"

	"github.com/blackjack/webcam"
	"github.com/pion/codec"
)

var cam *webcam.Webcam

func Start(w io.Writer, encoder codec.Encoder, res image.Rectangle) {
	var err error
	width := res.Max.X - res.Min.X
	height := res.Max.Y - res.Min.Y

	cam, err = webcam.Open("/dev/video0")
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

	go func() {
		frames := 0

		for {
			err = cam.WaitForFrame(1)
			switch err.(type) {
			case nil:
			case *webcam.Timeout:
				fmt.Fprint(os.Stderr, err)
				continue
			default:
				fmt.Fprint(os.Stderr, err)
				return
			}

			frame, err := cam.ReadFrame()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			if len(frame) != 0 {
				img, err := jpeg.Decode(bytes.NewReader(frame))
				if err != nil {
					continue
				}

				encoded, err := encoder.Encode(img)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}

				frames++
				fmt.Fprintln(os.Stderr, frames)

				bytesLeft := len(encoded)
				for bytesLeft > 0 {
					n, err := w.Write(encoded)
					if err != nil {
						continue
					}

					bytesLeft -= n
				}
			}
		}
	}()
}

func Stop() {
	cam.Close()
}
