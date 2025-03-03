package mp4

import (
	"encoding/hex"
	"io"

	"github.com/vtpl1/mp4ff/av1"
	"github.com/vtpl1/mp4ff/bits"
)

type Av1CBox struct {
	av1.CodecConfRec
}

// DecodeAv1C - box-specific decode
func DecodeAv1C(hdr BoxHeader, startPos uint64, r io.Reader) (Box, error) {
	data, err := readBoxBody(r, hdr)
	if err != nil {
		return nil, err
	}
	av1DecConfRec, err := av1.DecodeAV1CodecConfRec(data)
	if err != nil {
		return nil, err
	}
	return &Av1CBox{av1DecConfRec}, nil
}

// DecodeAv1CSR - box-specific decode
func DecodeAv1CSR(hdr BoxHeader, startPos uint64, sr bits.SliceReader) (Box, error) {
	av1DecConfRec, err := av1.DecodeAV1CodecConfRec(sr.ReadBytes(hdr.payloadLen()))
	return &Av1CBox{av1DecConfRec}, err
}

// Type - return box type
func (b *Av1CBox) Type() string {
	return "av1C"
}

// Size - return calculated size
func (b *Av1CBox) Size() uint64 {
	return uint64(boxHeaderSize + b.CodecConfRec.Size())
}

// Encode - write box to w
func (b *Av1CBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	return b.CodecConfRec.Encode(w)
}

// Encode - write box to sw
func (b *Av1CBox) EncodeSW(sw bits.SliceWriter) error {
	err := EncodeHeaderSW(b, sw)
	if err != nil {
		return err
	}
	return b.CodecConfRec.EncodeSW(sw)
}

// Info - box-specific Info
func (b *Av1CBox) Info(w io.Writer, specificBoxLevels, indent, indentStep string) error {
	bd := newInfoDumper(w, indent, b, -1, 0)
	bd.writef(" - SeqProfile: %d", b.SeqProfile)
	bd.writef(" - SeqLevelIdx0: %d", b.SeqLevelIdx0)
	bd.writef(" - SeqTier0: %d", b.SeqTier0)
	bd.writef(" - HighBitdepth: %d", b.HighBitdepth)
	bd.writef(" - TwelveBit: %d", b.TwelveBit)
	bd.writef(" - MonoChrome: %d", b.MonoChrome)
	bd.writef(" - ChromaSubsamplingX: %d", b.ChromaSubsamplingX)
	bd.writef(" - ChromaSubsamplingY: %d", b.ChromaSubsamplingY)
	bd.writef(" - ChromaSamplePosition: %d", b.ChromaSamplePosition)
	if b.InitialPresentationDelayPresent == 1 {
		bd.writef(" - InitialPresentationDelayMinusOne: %d", b.InitialPresentationDelayMinusOne)
	}
	bd.writef("   - ConfigOBUs: %s", hex.EncodeToString(b.ConfigOBUs))
	return bd.err
}
