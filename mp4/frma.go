package mp4

import (
	"errors"
	"io"

	"github.com/vtpl1/mp4ff/internal/bits"
)

var ErrFrmaContentLengthIsNot4 = errors.New("frma content length is not 4")

// FrmaBox - Original Format Box.
type FrmaBox struct {
	DataFormat string // uint32 - original box type
}

// Encode implements Box.
func (f *FrmaBox) Encode(w io.Writer) error {
	sw := bits.NewFixedSliceWriter(int(f.Size())) //nolint:gosec
	err := f.EncodeSW(sw)
	if err != nil {
		return err
	}
	_, err = w.Write(sw.Bytes())
	return err
}

// EncodeSW implements Box.
func (f *FrmaBox) EncodeSW(sw bits.SliceWriter) error {
	err := EncodeHeaderSW(f, sw)
	if err != nil {
		return err
	}
	sw.WriteString(f.DataFormat, false)
	return sw.AccError()
}

// Info implements Box.
func (f *FrmaBox) Info(w io.Writer, _ string, indent string, _ string) error {
	bd := newInfoDumper(w, indent, f, -1, 0)
	bd.writef(" - dataFormat: %s", f.DataFormat)
	return bd.err
}

// Size implements Box.
func (f *FrmaBox) Size() uint64 {
	return 12
}

// Type implements Box.
func (f *FrmaBox) Type() string {
	return "frma"
}

// DecodeFrma - box-specific decode.
func DecodeFrma(hdr BoxHeader, startPos uint64, r io.Reader) (Box, error) {
	data, err := readBoxBody(r, hdr)
	if err != nil {
		return nil, err
	}
	sr := bits.NewFixedSliceReader(data)
	return DecodeFrmaSR(hdr, startPos, sr)
}

// DecodeFrmaSR - box-specific decode.
func DecodeFrmaSR(hdr BoxHeader, _ uint64, sr bits.SliceReader) (Box, error) {
	if hdr.payloadLen() != 4 {
		return nil, ErrFrmaContentLengthIsNot4
	}
	return &FrmaBox{DataFormat: sr.ReadFixedLengthString(4)}, sr.AccError()
}
