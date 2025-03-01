package mp4_test

import (
	"testing"

	"github.com/vtpl1/mp4ff/mp4"
)

func TestNmhd(t *testing.T) {
	encBox := &mp4.NmhdBox{}
	boxDiffAfterEncodeAndDecode(t, encBox)
}
