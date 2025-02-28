package mp4_test

import (
	"testing"

	"github.com/vtpl1/mp4ff/mp4"
)

func TestEncDecTenc(t *testing.T) {
	kid := mp4.MustCreateUUID(mp4.UUIDPiffSenc)
	b := &mp4.TencBox{Version: 0, DefaultIsProtected: 1, DefaultPerSampleIVSize: 16, DefaultKID: kid}
	boxDiffAfterEncodeAndDecode(t, b)
}

func TestTencBox(t *testing.T) {
	var b mp4.Box = mp4.TencBox{}
	_ = b
}
