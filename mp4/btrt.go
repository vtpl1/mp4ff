package mp4

import (
	"io"

	"github.com/vtpl1/mp4ff/internal/bits"
)

// BtrtBox - BitRateBox - ISO/IEC 14496-12 Section 8.5.2.2.
type BtrtBox struct {
	BufferSizeDB uint32
	MaxBitrate   uint32
	AvgBitrate   uint32
}

// Encode implements Box.
func (b *BtrtBox) Encode(w io.Writer) error {
	sw := bits.NewFixedSliceWriter(int(b.Size())) //nolint:gosec
	err := b.EncodeSW(sw)
	if err != nil {
		return err
	}
	_, err = w.Write(sw.Bytes())
	return err
}

// EncodeSW implements Box.
func (b *BtrtBox) EncodeSW(sw bits.SliceWriter) error {
	err := EncodeHeaderSW(b, sw)
	if err != nil {
		return err
	}
	sw.WriteUint32(b.BufferSizeDB)
	sw.WriteUint32(b.MaxBitrate)
	sw.WriteUint32(b.AvgBitrate)
	return sw.AccError()
}

// Info implements Box.
func (b *BtrtBox) Info(w io.Writer, _ string, indent string, _ string) error {
	bd := newInfoDumper(w, indent, b, -1, 0)
	bd.writef(" - bufferSizeDB: %d", b.BufferSizeDB)
	bd.writef(" - maxBitrate: %d", b.MaxBitrate)
	bd.writef(" - AvgBitrate: %d", b.AvgBitrate)
	return bd.err
}

// Size implements Box.
func (b *BtrtBox) Size() uint64 {
	return 20
}

// Type implements Box.
func (b *BtrtBox) Type() string {
	return "btrt"
}

// DecodeBtrt - box-specific decode.
func DecodeBtrt(hdr BoxHeader, startPos uint64, r io.Reader) (Box, error) {
	data, err := readBoxBody(r, hdr)
	if err != nil {
		return nil, err
	}
	sr := bits.NewFixedSliceReader(data)
	return DecodeBtrtSR(hdr, startPos, sr)
}

// DecodeBtrtSR - box-specific decode.
func DecodeBtrtSR(_ BoxHeader, _ uint64, sr bits.SliceReader) (Box, error) {
	b := &BtrtBox{
		BufferSizeDB: sr.ReadUint32(),
		MaxBitrate:   sr.ReadUint32(),
		AvgBitrate:   sr.ReadUint32(),
	}
	return b, sr.AccError()
}
