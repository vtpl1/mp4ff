package mp4_test

import (
	"testing"

	"github.com/vtpl1/mp4ff/mp4"
)

func TestEncDecSchi(t *testing.T) {
	b := &mp4.SchiBox{}
	boxDiffAfterEncodeAndDecode(t, b)
}

func TestSchiBox(t *testing.T) {
	var b mp4.Box = mp4.SchiBox{}
	_ = b
}
