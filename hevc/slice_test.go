package hevc_test

import (
	"os"
	"testing"

	"github.com/vtpl1/mp4ff/avc"
	"github.com/vtpl1/mp4ff/hevc"

	"github.com/go-test/deep"
)

func TestParseSliceHeader(t *testing.T) {
	wantedHdr := map[hevc.NaluType]hevc.SliceHeader{
		hevc.NALU_IDR_N_LP: {
			SliceType:                         hevc.SLICE_I,
			FirstSliceSegmentInPicFlag:        true,
			SaoLumaFlag:                       true,
			SaoChromaFlag:                     true,
			QpDelta:                           7,
			LoopFilterAcrossSlicesEnabledFlag: true,
			NumEntryPointOffsets:              1,
			OffsetLenMinus1:                   3,
			CollocatedFromL0Flag:              true,
			EntryPointOffsetMinus1:            []uint32{12},
			Size:                              6,
		},
		hevc.NALU_TRAIL_N: {
			SliceType:                  hevc.SLICE_B,
			FirstSliceSegmentInPicFlag: true,
			PicOrderCntLsb:             1,
			ShortTermRefPicSet: hevc.ShortTermRPS{
				DeltaPocS0:      []uint32{1},
				DeltaPocS1:      []uint32{2, 2},
				UsedByCurrPicS0: []bool{true},
				UsedByCurrPicS1: []bool{true, true},
				NumNegativePics: 1,
				NumPositivePics: 2,
				NumDeltaPocs:    3,
			},
			SaoLumaFlag:                       true,
			SaoChromaFlag:                     true,
			TemporalMvpEnabledFlag:            true,
			NumRefIdxActiveOverrideFlag:       true,
			NumRefIdxL0ActiveMinus1:           0,
			NumRefIdxL1ActiveMinus1:           1,
			FiveMinusMaxNumMergeCand:          2,
			QpDelta:                           10,
			LoopFilterAcrossSlicesEnabledFlag: true,
			NumEntryPointOffsets:              1,
			OffsetLenMinus1:                   1,
			CollocatedFromL0Flag:              false,
			EntryPointOffsetMinus1:            []uint32{1},
			Size:                              10,
		},
		hevc.NALU_TRAIL_R: {
			SliceType:                  hevc.SLICE_P,
			FirstSliceSegmentInPicFlag: true,
			PicOrderCntLsb:             5,
			ShortTermRefPicSet: hevc.ShortTermRPS{
				DeltaPocS0:      []uint32{5},
				DeltaPocS1:      []uint32{},
				UsedByCurrPicS0: []bool{true},
				UsedByCurrPicS1: []bool{},
				NumNegativePics: 1,
				NumDeltaPocs:    1,
			},
			SaoLumaFlag:            true,
			SaoChromaFlag:          true,
			TemporalMvpEnabledFlag: true,
			PredWeightTable: &hevc.PredWeightTable{
				LumaLog2WeightDenom:        7,
				DeltaChromaLog2WeightDenom: -1,
				WeightsL0: []hevc.WeightingFactors{
					{
						LumaWeightFlag:   false,
						ChromaWeightFlag: false,
					},
				},
			},
			CollocatedFromL0Flag:     true,
			FiveMinusMaxNumMergeCand: 2,
			QpDelta:                  7,
			NumEntryPointOffsets:     1,
			OffsetLenMinus1:          1,
			EntryPointOffsetMinus1:   []uint32{2},
			Size:                     10,
		},
	}
	data, err := os.ReadFile("testdata/blackframe.265")
	if err != nil {
		t.Error(err)
	}
	nalus := avc.ExtractNalusFromByteStream(data)
	spsMap := make(map[uint32]*hevc.SPS, 1)
	ppsMap := make(map[uint32]*hevc.PPS, 1)
	gotHdr := make(map[hevc.NaluType]hevc.SliceHeader)
	for _, nalu := range nalus {
		switch hevc.GetNaluType(nalu[0]) {
		case hevc.NALU_SPS:
			sps, err := hevc.ParseSPSNALUnit(nalu)
			if err != nil {
				t.Error(err)
			}
			spsMap[uint32(sps.SpsID)] = sps
		case hevc.NALU_PPS:
			pps, err := hevc.ParsePPSNALUnit(nalu, spsMap)
			if err != nil {
				t.Error(err)
			}
			ppsMap[pps.PicParameterSetID] = pps
		case hevc.NALU_IDR_N_LP, hevc.NALU_TRAIL_R, hevc.NALU_TRAIL_N:
			hdr, err := hevc.ParseSliceHeader(nalu, spsMap, ppsMap)
			if err != nil {
				t.Error(err)
			}
			gotHdr[hevc.GetNaluType(nalu[0])] = *hdr
		}
	}
	if diff := deep.Equal(wantedHdr, gotHdr); diff != nil {
		t.Errorf("Got Slice Headers: %+v\n Diff is %v", gotHdr, diff)
	}
}
