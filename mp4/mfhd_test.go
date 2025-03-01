package mp4_test

import (
	"testing"

	"github.com/vtpl1/mp4ff/mp4"
)

func TestEncodeDecodeMfhd(t *testing.T) {
	mfhd := &mp4.MfhdBox{
		Version:        0,
		Flags:          0,
		SequenceNumber: 1,
	}
	boxDiffAfterEncodeAndDecode(t, mfhd)
}
