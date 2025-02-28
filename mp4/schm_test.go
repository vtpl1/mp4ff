package mp4_test

import (
	"testing"

	"github.com/vtpl1/mp4ff/mp4"
)

func TestSchm(t *testing.T) {
	b := &mp4.SchmBox{SchemeType: "cenc", SchemeVersion: 65536}
	boxDiffAfterEncodeAndDecode(t, b)
}

func TestSchmBox(t *testing.T) {
	var b mp4.Box = &mp4.SchmBox{}
	_ = b
}
