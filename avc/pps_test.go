package avc_test

import (
	"encoding/hex"
	"testing"

	"github.com/go-test/deep"
	"github.com/vtpl1/mp4ff/avc"
)

const pps1 = "68e84332c8b0"

func TestPPSParser(t *testing.T) {
	byteData, _ := hex.DecodeString(pps1)

	wanted := &avc.PPS{
		PicParameterSetID:                     0,
		SeqParameterSetID:                     0,
		EntropyCodingModeFlag:                 true,
		BottomFieldPicOrderInFramePresentFlag: false,
		NumSliceGroupsMinus1:                  0,
		NumRefIdxI0DefaultActiveMinus1:        15,
		NumRefIdxI1DefaultActiveMinus1:        0,
		WeightedPredFlag:                      true,
		WeightedBipredIDC:                     0,
		PicInitQpMinus26:                      0,
		PicInitQsMinus26:                      0,
		ChromaQpIndexOffset:                   -2,
		DeblockingFilterControlPresentFlag:    true,
		ConstrainedIntraPredFlag:              false,
		RedundantPicCntPresentFlag:            false,
		Transform8x8ModeFlag:                  true,
		PicScalingMatrixPresentFlag:           false,
		PicScalingLists:                       nil,
		SecondChromaQpIndexOffset:             -2,
	}
	got, err := avc.ParsePPSNALUnit(byteData, nil)
	if err != nil {
		t.Error(err)
		return
	}
	if diff := deep.Equal(got, wanted); diff != nil {
		t.Error(diff)
	}
}
