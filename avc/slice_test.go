package avc_test

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/vtpl1/mp4ff/avc"
)

const videoNaluStart = "25888040ffde08e47a7bff05ab"

func TestSliceTypeParser(t *testing.T) {
	byteData, _ := hex.DecodeString(videoNaluStart)
	want := avc.SLICE_I
	got, err := avc.GetSliceTypeFromNALU(byteData)
	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestSliceTypeStrings(t *testing.T) {
	cases := []struct {
		sliceType avc.SliceType
		want      string
	}{
		{avc.SLICE_P, "P"},
		{avc.SLICE_B, "B"},
		{avc.SLICE_I, "I"},
		{avc.SLICE_SP, "SP"},
		{avc.SLICE_SI, "SI"},
		{avc.SliceType(12), ""},
	}
	for _, c := range cases {
		got := c.sliceType.String()
		if got != c.want {
			t.Errorf("got %s want %s", got, c.want)
		}
	}
}

func TestParseSliceHeader_BlackFrame(t *testing.T) {
	wantedHdr := avc.SliceHeader{
		SliceType:              7,
		SliceQPDelta:           6,
		SliceAlphaC0OffsetDiv2: -3,
		SliceBetaOffsetDiv2:    -3,
		Size:                   7,
	}
	data, err := os.ReadFile("testdata/blackframe.264")
	if err != nil {
		t.Error(err)
	}
	nalus := avc.ExtractNalusFromByteStream(data)
	spsMap := make(map[uint32]*avc.SPS, 1)
	ppsMap := make(map[uint32]*avc.PPS, 1)
	var gotHdr *avc.SliceHeader
	for _, nalu := range nalus {
		switch avc.GetNaluType(nalu[0]) {
		case avc.NALU_SPS:
			sps, err := avc.ParseSPSNALUnit(nalu, true)
			if err != nil {
				t.Error(err)
			}
			spsMap[uint32(sps.ParameterID)] = sps
		case avc.NALU_PPS:
			pps, err := avc.ParsePPSNALUnit(nalu, spsMap)
			if err != nil {
				t.Error(err)
			}
			ppsMap[uint32(pps.PicParameterSetID)] = pps
		case avc.NALU_IDR:
			gotHdr, err = avc.ParseSliceHeader(nalu, spsMap, ppsMap)
			if err != nil {
				t.Error(err)
			}
		}
	}
	if diff := deep.Equal(wantedHdr, *gotHdr); diff != nil {
		t.Errorf("Got slice header %+v. Diff=%v", *gotHdr, diff)
	}
}

func TestParseSliceHeader_TwoFrames(t *testing.T) {
	wantedIdrHdr := avc.SliceHeader{SliceType: avc.SLICE_I, IDRPicID: 1, SliceQPDelta: 8, Size: 5}
	wantedNonIdrHdr := avc.SliceHeader{
		SliceType: avc.SLICE_P, FrameNum: 1, ModificationOfPicNumsIDC: 3, SliceQPDelta: 13,
		Size: 5, NumRefIdxActiveOverrideFlag: true, RefPicListModificationL0Flag: true,
	}

	data, err := os.ReadFile("testdata/two-frames.264")
	if err != nil {
		t.Error(err)
	}
	nalus, err := avc.GetNalusFromSample(data)
	if err != nil {
		t.Error(err)
	}
	spsMap := make(map[uint32]*avc.SPS, 1)
	ppsMap := make(map[uint32]*avc.PPS, 1)
	var gotIdrHdr *avc.SliceHeader
	var gotNonIdrHdr *avc.SliceHeader
	for _, nalu := range nalus {
		switch avc.GetNaluType(nalu[0]) {
		case avc.NALU_SPS:
			sps, err := avc.ParseSPSNALUnit(nalu, true)
			if err != nil {
				t.Error(err)
			}
			spsMap[uint32(sps.ParameterID)] = sps
		case avc.NALU_PPS:
			pps, err := avc.ParsePPSNALUnit(nalu, spsMap)
			if err != nil {
				t.Error(err)
			}
			ppsMap[uint32(pps.PicParameterSetID)] = pps
		case avc.NALU_IDR:
			gotIdrHdr, err = avc.ParseSliceHeader(nalu, spsMap, ppsMap)
			if err != nil {
				t.Error(err)
			}
		case avc.NALU_NON_IDR:
			gotNonIdrHdr, err = avc.ParseSliceHeader(nalu, spsMap, ppsMap)
			if err != nil {
				t.Error(err)
			}
		}
	}
	if diff := deep.Equal(wantedIdrHdr, *gotIdrHdr); diff != nil {
		t.Errorf("Got IDR Slice Header: %+v\n Diff is: %v", *gotIdrHdr, diff)
	}
	if diff := deep.Equal(wantedNonIdrHdr, *gotNonIdrHdr); diff != nil {
		t.Errorf("Got NON_IDR Slice Header: %+v\n Diff is: %v", *gotNonIdrHdr, diff)
	}
}

func TestParseSliceHeaderLength(t *testing.T) {
	spsHex := "6764001eacd940a02ff9610000030001000003003c8f162d96"
	ppsHex := "68ebecb22c"
	naluStartHex := "419a6649e10f2653022fff8700000302c8a32d32"
	spsData, _ := hex.DecodeString(spsHex)
	sps, err := avc.ParseSPSNALUnit(spsData, true)
	if err != nil {
		t.Error(err)
	}
	spsMap := make(map[uint32]*avc.SPS, 1)
	spsMap[uint32(sps.ParameterID)] = sps
	ppsData, _ := hex.DecodeString(ppsHex)
	pps, err := avc.ParsePPSNALUnit(ppsData, spsMap)
	if err != nil {
		t.Error(err)
	}
	ppsMap := make(map[uint32]*avc.PPS, 1)
	ppsMap[uint32(pps.PicParameterSetID)] = pps
	naluStart, _ := hex.DecodeString(naluStartHex)
	sh, err := avc.ParseSliceHeader(naluStart, spsMap, ppsMap)
	if err != nil {
		t.Error(err)
	}
	wantedSliceHeaderSize := uint32(11)
	if sh.Size != wantedSliceHeaderSize {
		t.Errorf("got %d want %d", sh.Size, wantedSliceHeaderSize)
	}
}
