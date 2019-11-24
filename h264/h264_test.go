package h264

import (
	"testing"
)

func TestNewH264Encoder(t *testing.T) {
	_, err := NewEncoder(H264Options{
		Height:       1080,
		Width:        1920,
		Bitrate:      5000000,
		MaxFrameRate: 60,
	})
	if err != nil {
		t.Fatal(err)
	}
}
