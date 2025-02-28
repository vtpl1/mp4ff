package mp4_test

import (
	"testing"

	"github.com/vtpl1/mp4ff/mp4"
)

func TestSaiz(t *testing.T) {
	saiz := &mp4.SaizBox{DefaultSampleInfoSize: 1}
	boxDiffAfterEncodeAndDecode(t, saiz)
}

func TestSaizBox(t *testing.T) {
	var b mp4.Box = &mp4.SaizBox{DefaultSampleInfoSize: 1}
	_ = b
}
