package av1_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/go-test/deep"
	"github.com/vtpl1/mp4ff/av1"
)

const (
	av1DecoderConfigRecord = "81094c000a0b0000004aabbfc377ffe701"
	configOBUs             = "0a0b0000004aabbfc377ffe701"
)

func TestDecodeAV1DecConfRec(t *testing.T) {
	byteData, _ := hex.DecodeString(av1DecoderConfigRecord)
	configOBUsBytes, _ := hex.DecodeString(configOBUs)

	wanted := av1.CodecConfRec{
		Version:                          1,
		SeqProfile:                       0,
		SeqLevelIdx0:                     9,
		SeqTier0:                         0,
		HighBitdepth:                     1,
		TwelveBit:                        0,
		MonoChrome:                       0,
		ChromaSubsamplingX:               1,
		ChromaSubsamplingY:               1,
		ChromaSamplePosition:             0,
		InitialPresentationDelayPresent:  0,
		InitialPresentationDelayMinusOne: 0,
		ConfigOBUs:                       configOBUsBytes,
	}

	got, err := av1.DecodeAV1CodecConfRec(byteData)
	if err != nil {
		t.Error("Error parsing Av1DecoderConfigRecord")
	}
	if diff := deep.Equal(got, wanted); diff != nil {
		t.Error(diff)
	}

	encBuf := bytes.Buffer{}
	err = got.Encode(&encBuf)
	if err != nil {
		t.Error("Error encoding Av1DecoderConfigRecord")
	}
	encBytes := encBuf.Bytes()
	if diff := deep.Equal(encBytes, byteData); diff != nil {
		t.Error(diff)
	}
}
