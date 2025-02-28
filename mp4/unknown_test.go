package mp4_test

import (
	"testing"

	"github.com/vtpl1/mp4ff/mp4"
)

// TestUnknown including non-ascii character in name (box typs is uint32 according to spec)
func TestUnknown(t *testing.T) {
	unknownBox := &mp4.UnknownBox{
		Name:       "\xa9enc",
		SizeN:      12,
		NotDecoded: []byte{0, 0, 0, 0},
	}

	boxDiffAfterEncodeAndDecode(t, unknownBox)
}
