package sei_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/go-test/deep"
	"github.com/vtpl1/mp4ff/sei"
)

func TestHEVCSETI1PicTiming(t *testing.T) {
	cases := []struct {
		name           string
		naluPayloadHex string
		extParams      sei.HEVCPicTimingParams
		expected       sei.PicTimingHevcSEI
		expNonFatalErr error
	}{
		{
			name:           "HEVC_SETI1PicTiming",
			naluPayloadHex: "01071000001a0000030180",
			extParams: sei.HEVCPicTimingParams{
				FrameFieldInfoPresentFlag:              true,
				CpbDpbDelaysPresentFlag:                true,
				SubPicHrdParamsPresentFlag:             false,
				SubPicCpbParamsInPicTimingSeiFlag:      false,
				AuCbpRemovalDelayLengthMinus1:          23,
				DpbOutputDelayLengthMinus1:             0,
				DpbOutputDelayDuLengthMinus1:           23,
				DuCpbRemovalDelayIncrementLengthMinus1: 0,
			},
			expected: sei.PicTimingHevcSEI{
				FrameFieldInfo: &sei.HEVCFrameFieldInfo{
					PicStruct:      1,
					SourceScanType: 0,
					DuplicateFlag:  false,
				},
			},
			expNonFatalErr: nil,
		},
	}

	for _, tc := range cases {
		seiNaluPayload, _ := hex.DecodeString(tc.naluPayloadHex)
		r := bytes.NewReader(seiNaluPayload)
		seis, err := sei.ExtractSEIData(r)
		if err != nil && err != tc.expNonFatalErr {
			t.Error(err)
		}
		if len(seis) != 1 {
			t.Errorf("%s: Not %d but %d sei messages found", tc.name, 1, len(seis))
		}
		seiMessage, err := sei.DecodePicTimingHevcSEI(&seis[0], tc.extParams)
		if err != nil {
			t.Error(err)
		}
		if seiMessage.Type() != sei.SEIPicTimingType {
			t.Errorf("%s: got SEI type %d instead of %d", tc.name, seiMessage.Type(), sei.SEIPicTimingType)
		}
		seiPT := seiMessage.(*sei.PicTimingHevcSEI)
		diff := deep.Equal(seiPT.FrameFieldInfo, tc.expected.FrameFieldInfo)
		if diff != nil {
			t.Errorf("%s: %v %s", tc.name, diff, "frame field info mismatch")
		}
	}
}
