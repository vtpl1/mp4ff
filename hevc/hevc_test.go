package hevc_test

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/vtpl1/mp4ff/hevc"
)

func TestGetNaluTypes(t *testing.T) {
	testCases := []struct {
		name                string
		input               []byte
		wanted              []hevc.NaluType
		nalusUpToFirstVideo []hevc.NaluType
		containsVPS         bool
		isRapSample         bool
		isIDRSample         bool
	}{
		{
			"AUD",
			[]byte{0, 0, 0, 2, 70, 0},
			[]hevc.NaluType{hevc.NALU_AUD},
			[]hevc.NaluType{hevc.NALU_AUD},
			false,
			false,
			false,
		},
		{
			"AUD, VPS, SPS, PPS, and IDR ",
			[]byte{
				0, 0, 0, 2, 70, 2,
				0, 0, 0, 3, 64, 1, 1,
				0, 0, 0, 3, 66, 2, 2,
				0, 0, 0, 3, 68, 3, 3,
				0, 0, 0, 3, 40, 4, 4,
			},
			[]hevc.NaluType{hevc.NALU_AUD, hevc.NALU_VPS, hevc.NALU_SPS, hevc.NALU_PPS, hevc.NALU_IDR_N_LP},
			[]hevc.NaluType{hevc.NALU_AUD, hevc.NALU_VPS, hevc.NALU_SPS, hevc.NALU_PPS, hevc.NALU_IDR_N_LP},
			true,
			true,
			true,
		},
		{
			"too short",
			[]byte{0, 0, 0},
			[]hevc.NaluType{},
			[]hevc.NaluType{},
			false,
			false,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := hevc.FindNaluTypes(tc.input)
			if diff := deep.Equal(got, tc.wanted); diff != nil {
				t.Errorf("nalulist diff: %v", diff)
			}
			got = hevc.FindNaluTypesUpToFirstVideoNalu(tc.input)
			if diff := deep.Equal(got, tc.nalusUpToFirstVideo); diff != nil {
				t.Errorf("nalus before first video diff: %v", diff)
			}
			hasVPS := hevc.ContainsNaluType(tc.input, hevc.NALU_VPS)
			if hasVPS != tc.containsVPS {
				t.Errorf("got %t instead of %t", hasVPS, tc.containsVPS)
			}
			isRAP := hevc.IsRAPSample(tc.input)
			if isRAP != tc.isRapSample {
				t.Errorf("got %t instead of %t", isRAP, tc.isRapSample)
			}
			isIDR := hevc.IsIDRSample(tc.input)
			if isIDR != tc.isIDRSample {
				t.Errorf("got %t instead of %t", isIDR, tc.isIDRSample)
			}
		})
	}
}

func TestHasParameterSets(t *testing.T) {
	testCases := []struct {
		name   string
		input  []byte
		wanted bool
	}{
		{
			"AUD",
			[]byte{0, 0, 0, 2, 70, 0},
			false,
		},
		{
			"AUD, VPS, SPS, PPS, and IDR ",
			[]byte{
				0, 0, 0, 2, 70, 2,
				0, 0, 0, 3, 64, 1, 1,
				0, 0, 0, 3, 66, 2, 2,
				0, 0, 0, 3, 68, 3, 3,
				0, 0, 0, 3, 40, 4, 4,
			},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := hevc.HasParameterSets(tc.input)
			if got != tc.wanted {
				t.Errorf("got %t instead of %t", got, tc.wanted)
			}
		})
	}
}

func TestGetParameterSets(t *testing.T) {
	testCases := []struct {
		name      string
		input     []byte
		wantedVPS [][]byte
		wantedSPS [][]byte
		wantedPPS [][]byte
	}{
		{
			"AUD",
			[]byte{0, 0, 0, 2, 70, 0},
			nil, nil, nil,
		},
		{
			"AUD, VPS, SPS, PPS, and IDR ",
			[]byte{
				0, 0, 0, 2, 70, 2,
				0, 0, 0, 3, 64, 1, 1,
				0, 0, 0, 3, 66, 2, 2,
				0, 0, 0, 3, 68, 3, 3,
				0, 0, 0, 3, 40, 4, 4,
			},
			[][]byte{{64, 1, 1}},
			[][]byte{{66, 2, 2}},
			[][]byte{{68, 3, 3}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotVPS, gotSPS, gotPPS := hevc.GetParameterSets(tc.input)
			if diff := deep.Equal(gotVPS, tc.wantedVPS); diff != nil {
				t.Errorf("VPS diff: %v", diff)
			}
			if diff := deep.Equal(gotSPS, tc.wantedSPS); diff != nil {
				t.Errorf("SPS diff: %v", diff)
			}
			if diff := deep.Equal(gotPPS, tc.wantedPPS); diff != nil {
				t.Errorf("PPS diff: %v", diff)
			}
		})
	}
}

func TestNaluTypeStrings(t *testing.T) {
	named := 0
	for n := hevc.NaluType(0); n < hevc.NaluType(64); n++ {
		desc := n.String()
		if !strings.HasPrefix(desc, "Other") {
			named++
		}
	}
	if named != 22 {
		t.Errorf("got %d named instead of 22", named)
	}
}

func TestIsVideoNaluType(t *testing.T) {
	testCases := []struct {
		name     string
		naluType hevc.NaluType
		want     bool
	}{
		{
			name:     "video type - NALU_TRAIL_N (0)",
			naluType: hevc.NALU_TRAIL_N,
			want:     true,
		},
		{
			name:     "video type - NALU_TRAIL_R (1)",
			naluType: hevc.NALU_TRAIL_R,
			want:     true,
		},
		{
			name:     "video type - NALU_IDR_W_RADL (19)",
			naluType: hevc.NALU_IDR_W_RADL,
			want:     true,
		},
		{
			name:     "video type - highest (31)",
			naluType: 31,
			want:     true,
		},
		{
			name:     "non-video type - VPS (32)",
			naluType: hevc.NALU_VPS,
			want:     false,
		},
		{
			name:     "non-video type - SPS (33)",
			naluType: hevc.NALU_SPS,
			want:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := hevc.IsVideoNaluType(tc.naluType)
			if got != tc.want {
				t.Errorf("IsVideoNaluType(%d) = %v; want %v", tc.naluType, got, tc.want)
			}
		})
	}
}
