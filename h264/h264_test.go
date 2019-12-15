package h264

import (
	"image"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/pion/codec/internal/camera"
)

func TestNewH264Encoder(t *testing.T) {
	_, err := NewEncoder(Options{
		Height:       1080,
		Width:        1920,
		Bitrate:      5000000,
		MaxFrameRate: 60,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkLatency(b *testing.B) {
	width := 640
	height := 480
	res := image.Rect(0, 0, width, height)
	enc, err := NewEncoder(Options{
		Height:       height,
		Width:        width,
		Bitrate:      2000000,
		MaxFrameRate: 60,
	})
	defer enc.Close()

	if err != nil {
		log.Fatal(err)
	}

	camera.Start(ioutil.Discard, enc, res)
	defer camera.Stop()

	<-time.After(time.Second * 30)
}
