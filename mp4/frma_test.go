package mp4_test

import (
	"testing"

	"github.com/vtpl1/mp4ff/mp4"
)

func TestFrma(t *testing.T) {
	frma := &mp4.FrmaBox{DataFormat: "avc1"}
	boxDiffAfterEncodeAndDecode(t, frma)
}

func TestFrmaBox(t *testing.T) {
	var b mp4.Box = &mp4.FrmaBox{}
	_ = b
}
