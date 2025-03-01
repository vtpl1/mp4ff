package avc_test

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/vtpl1/mp4ff/avc"
)

func TestGetNaluTypes(t *testing.T) {
	testCases := []struct {
		name   string
		input  []byte
		wanted []avc.NaluType
	}{
		{
			"IDR",
			[]byte{0, 0, 0, 2, 5, 0},
			[]avc.NaluType{avc.NALU_IDR},
		},
		{
			"AUD and SPS",
			[]byte{0, 0, 0, 2, 9, 2, 0, 0, 0, 3, 7, 5, 4},
			[]avc.NaluType{avc.NALU_AUD, avc.NALU_SPS},
		},
	}

	for _, tc := range testCases {
		got := avc.FindNaluTypes(tc.input)
		if diff := deep.Equal(got, tc.wanted); diff != nil {
			t.Errorf("%s: %v", tc.name, diff)
		}
		for _, w := range tc.wanted {
			if !avc.ContainsNaluType(tc.input, w) {
				t.Errorf("%s: wanted %v in %v", tc.name, w, tc.input)
			}
			if w == avc.NALU_IDR && !avc.IsIDRSample(tc.input) {
				t.Errorf("%s: wanted IDR in %v", tc.name, tc.input)
			}
		}
	}
}

func TestNaluTypeStrings(t *testing.T) {
	named := make([]int, 0, 9)
	for i := 0; i < 14; i++ {
		nt := avc.NaluType(i)
		if !strings.HasPrefix(nt.String(), "Other") {
			named = append(named, i)
		}
	}
	if len(named) != 9 {
		t.Errorf("Expected 9 named NaluTypes, got %d", len(named))
	}
}

func TestHasParameterSets(t *testing.T) {
	testCases := []struct {
		name   string
		input  []byte
		wanted bool
	}{
		{
			"Only IDR",
			[]byte{0, 0, 0, 2, 5, 0},
			false,
		},
		{
			"AUD, SPS, PPS, IDRx2",
			[]byte{
				0, 0, 0, 2, 9, 2,
				0, 0, 0, 3, 7, 5, 4,
				0, 0, 0, 3, 8, 1, 2,
				0, 0, 0, 2, 5, 0,
				0, 0, 0, 2, 5, 0,
			},
			true,
		},
	}

	for _, tc := range testCases {
		got := avc.HasParameterSets(tc.input)
		if got != tc.wanted {
			t.Errorf("%s: got %t instead of %t", tc.name, got, tc.wanted)
		}
	}
}

func TestGetParameterSets(t *testing.T) {
	testCases := []struct {
		name      string
		input     []byte
		wantedSPS [][]byte
		wantedPPS [][]byte
	}{
		{
			"Only IDR",
			[]byte{0, 0, 0, 2, 5, 0},
			nil, nil,
		},
		{
			"AUD, SPS, PPS, IDRx2",
			[]byte{
				0, 0, 0, 2, 9, 2,
				0, 0, 0, 3, 7, 5, 4,
				0, 0, 0, 3, 8, 1, 2,
				0, 0, 0, 2, 5, 0,
				0, 0, 0, 2, 5, 0,
			},
			[][]byte{{7, 5, 4}},
			[][]byte{{8, 1, 2}},
		},
	}

	for _, tc := range testCases {
		gotSPS, gotPPS := avc.GetParameterSets(tc.input)
		if diff := deep.Equal(gotSPS, tc.wantedSPS); diff != nil {
			t.Errorf("%s: %v", tc.name, diff)
		}
		if diff := deep.Equal(gotPPS, tc.wantedPPS); diff != nil {
			t.Errorf("%s: %v", tc.name, diff)
		}
	}
}

func TestIsVideoNaluType(t *testing.T) {
	testCases := []struct {
		name     string
		naluType avc.NaluType
		want     bool
	}{
		{
			name:     "video type - NALU_NON_IDR (1)",
			naluType: avc.NALU_NON_IDR,
			want:     true,
		},
		{
			name:     "video type - NALU_IDR (5)",
			naluType: avc.NALU_IDR,
			want:     true,
		},
		{
			name:     "non-video type - NALU_SEI (6)",
			naluType: avc.NALU_SEI,
			want:     false,
		},
		{
			name:     "non-video type - NALU_SPS (7)",
			naluType: avc.NALU_SPS,
			want:     false,
		},
		{
			name:     "non-video type - NALU_AUD (9)",
			naluType: avc.NALU_AUD,
			want:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := avc.IsVideoNaluType(tc.naluType)
			if got != tc.want {
				t.Errorf("IsVideoNaluType(%d) = %v; want %v", tc.naluType, got, tc.want)
			}
		})
	}
}
