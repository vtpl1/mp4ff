package mp4_test

import (
	"testing"

	"github.com/vtpl1/mp4ff/mp4"
)

func TestStyp(t *testing.T) {
	b := mp4.CreateStyp()
	boxDiffAfterEncodeAndDecode(t, b)
}

func TestStypBox(t *testing.T) {
	var b mp4.Box = mp4.CreateStyp()
	_ = b
}
