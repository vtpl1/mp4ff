package mp4

import (
	"io"

	"github.com/vtpl1/mp4ff/bits"
)

// SchmBox - Scheme Type Box
type SchmBox struct {
	Version       byte
	Flags         uint32
	SchemeType    string // 4CC represented as uint32
	SchemeVersion uint32
	SchemeURI     string // Absolute null-terminated URL
}

// DecodeSchm - box-specific decode
func DecodeSchm(hdr BoxHeader, startPos uint64, r io.Reader) (Box, error) {
	data, err := readBoxBody(r, hdr)
	if err != nil {
		return nil, err
	}
	sr := bits.NewFixedSliceReader(data)
	return DecodeSchmSR(hdr, startPos, sr)
}

// DecodeSchmSR - box-specific decode
func DecodeSchmSR(hdr BoxHeader, startPos uint64, sr bits.SliceReader) (Box, error) {
	versionAndFlags := sr.ReadUint32()
	version := byte(versionAndFlags >> 24)

	b := SchmBox{
		Version: version,
		Flags:   versionAndFlags & flagsMask,
	}
	b.SchemeType = sr.ReadFixedLengthString(4)
	b.SchemeVersion = sr.ReadUint32()
	if b.Flags&0x01 != 0 {
		b.SchemeURI = sr.ReadZeroTerminatedString(hdr.payloadLen())
	}
	return &b, sr.AccError()
}

// Type - return box type
func (s *SchmBox) Type() string {
	return "schm"
}

// Size - return calculated size
func (s *SchmBox) Size() uint64 {
	size := uint64(20)
	if s.Flags&0x01 != 0 {
		size += uint64(len(s.SchemeURI) + 1)
	}
	return size
}

// Encode - write box to w
func (s *SchmBox) Encode(w io.Writer) error {
	sw := bits.NewFixedSliceWriter(int(s.Size()))
	err := s.EncodeSW(sw)
	if err != nil {
		return err
	}
	_, err = w.Write(sw.Bytes())
	return err
}

// EncodeSW - box-specific encode to slicewriter
func (s *SchmBox) EncodeSW(sw bits.SliceWriter) error {
	err := EncodeHeaderSW(s, sw)
	if err != nil {
		return err
	}
	versionAndFlags := (uint32(s.Version) << 24) + s.Flags
	sw.WriteUint32(versionAndFlags)
	sw.WriteString(s.SchemeType, false)
	sw.WriteUint32(s.SchemeVersion)
	if s.Flags&0x01 != 0 {
		sw.WriteString(s.SchemeURI, true)
	}
	return sw.AccError()
}

// Info - write box info to w
func (s *SchmBox) Info(w io.Writer, specificBoxLevels, indent, indentStep string) (err error) {
	bd := newInfoDumper(w, indent, s, int(s.Version), s.Flags)
	bd.writef(" - schemeType: %s", s.SchemeType)
	bd.writef(" - schemeVersion: %d  (%d.%d)", s.SchemeVersion, s.SchemeVersion>>16, s.SchemeVersion&0xffff)
	if s.Flags&0x01 != 0 {
		bd.writef(" - schemeURI: %q", s.SchemeURI)
	}
	return bd.err
}
