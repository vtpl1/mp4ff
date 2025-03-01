package mp4_test

import (
	"testing"

	"github.com/vtpl1/mp4ff/mp4"
)

func TestSinf(t *testing.T) {
	b := &mp4.SinfBox{}
	boxDiffAfterEncodeAndDecode(t, b)
}

func TestSinfBox(t *testing.T) {
	var b mp4.Box = &mp4.SchmBox{}
	_ = b
}
