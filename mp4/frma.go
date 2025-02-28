package mp4

import (
	"fmt"
	"io"

	"github.com/vtpl1/mp4ff/bits"
)

// FrmaBox - Original Format Box
type FrmaBox struct {
	DataFormat string // uint32 - original box type
}

// DecodeFrma - box-specific decode
func DecodeFrma(hdr BoxHeader, startPos uint64, r io.Reader) (Box, error) {
	data, err := readBoxBody(r, hdr)
	if err != nil {
		return nil, err
	}
	sr := bits.NewFixedSliceReader(data)
	return DecodeFrmaSR(hdr, startPos, sr)
}

// DecodeFrmaSR - box-specific decode
func DecodeFrmaSR(hdr BoxHeader, startPos uint64, sr bits.SliceReader) (Box, error) {
	if hdr.payloadLen() != 4 {
		return nil, fmt.Errorf("frma content length is not 4")
	}
	return &FrmaBox{DataFormat: string(sr.ReadFixedLengthString(4))}, sr.AccError()
}

// Type - return box type
func (b *FrmaBox) Type() string {
	return "frma"
}

// Size - return calculated size
func (b *FrmaBox) Size() uint64 {
	return 12
}

// Encode - write box to w
func (f *FrmaBox) Encode(w io.Writer) error {
	sw := bits.NewFixedSliceWriter(int(f.Size()))
	err := f.EncodeSW(sw)
	if err != nil {
		return err
	}
	_, err = w.Write(sw.Bytes())
	return err
}

// EncodeSW - box-specific encode to slicewriter
func (f *FrmaBox) EncodeSW(sw bits.SliceWriter) error {
	err := EncodeHeaderSW(f, sw)
	if err != nil {
		return err
	}
	sw.WriteString(f.DataFormat, false)
	return sw.AccError()
}

// Info - write box info to w
func (f *FrmaBox) Info(w io.Writer, specificBoxLevels, indent, indentStep string) (err error) {
	bd := newInfoDumper(w, indent, f, -1, 0)
	bd.writef(" - dataFormat: %s", f.DataFormat)
	return bd.err
}
