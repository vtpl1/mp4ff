package mp4

import (
	"testing"
)

// TestUnknown including non-ascii character in name (box typs is uint32 according to spec)
func TestUnknown(t *testing.T) {
	unknownBox := &UnknownBox{
		Name:       "\xa9enc",
		SizeN:      12,
		NotDecoded: []byte{0, 0, 0, 0},
	}

	boxDiffAfterEncodeAndDecode(t, unknownBox)
}
