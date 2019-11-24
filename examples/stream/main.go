package main

import (
	"image"
	"log"
	"os"

	"github.com/pion/codec/h264"
	"github.com/pion/codec/internal/camera"
)

func main() {
	width := 320
	height := 240
	res := image.Rect(0, 0, width, height)
	enc, err := h264.NewEncoder(h264.Options{
		Height:       height,
		Width:        width,
		Bitrate:      2000000,
		MaxFrameRate: 60,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer enc.Close()

	camera.Start(os.Stdout, enc, res)
	select {}
}
