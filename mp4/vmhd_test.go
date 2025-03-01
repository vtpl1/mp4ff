package mp4_test

import (
	"testing"

	"github.com/vtpl1/mp4ff/mp4"
)

func TestVmhd(t *testing.T) {
	vmhd := mp4.CreateVmhd()

	boxDiffAfterEncodeAndDecode(t, vmhd)
}
