package mp4

import (
	"fmt"
	"io"

	"github.com/vtpl1/mp4ff/bits"
)

// SaizBox - Sample Auxiliary Information Sizes Box (saiz)  (in stbl or traf box)
type SaizBox struct {
	Version               byte
	Flags                 uint32
	AuxInfoType           string // Used for Common Encryption Scheme (4-bytes uint32 according to spec)
	AuxInfoTypeParameter  uint32
	SampleCount           uint32
	SampleInfo            []byte
	DefaultSampleInfoSize byte
}

// DecodeSaiz - box-specific decode
func DecodeSaiz(hdr BoxHeader, startPos uint64, r io.Reader) (Box, error) {
	data, err := readBoxBody(r, hdr)
	if err != nil {
		return nil, err
	}
	sr := bits.NewFixedSliceReader(data)
	return DecodeSaizSR(hdr, startPos, sr)
}

// DecodeSaizSR - box-specific decode
func DecodeSaizSR(hdr BoxHeader, startPos uint64, sr bits.SliceReader) (Box, error) {
	versionAndFlags := sr.ReadUint32()
	version := byte(versionAndFlags >> 24)
	b := SaizBox{
		Version: version,
		Flags:   versionAndFlags & flagsMask,
	}
	if b.Flags&0x01 != 0 {
		b.AuxInfoType = sr.ReadFixedLengthString(4)
		b.AuxInfoTypeParameter = sr.ReadUint32()
	}
	b.DefaultSampleInfoSize = sr.ReadUint8()
	b.SampleCount = sr.ReadUint32()

	if hdr.Size != b.expectedSize() {
		return nil, fmt.Errorf("saiz: expected size %d, got %d", b.expectedSize(), hdr.Size)
	}

	if b.DefaultSampleInfoSize == 0 {
		b.SampleInfo = make([]byte, 0, b.SampleCount)
		for i := uint32(0); i < b.SampleCount; i++ {
			b.SampleInfo = append(b.SampleInfo, sr.ReadUint8())
		}
	}
	return &b, sr.AccError()
}

// NewSaizBox creates a SaizBox with appropriate size allocated.
func NewSaizBox(capacity int) *SaizBox {
	return &SaizBox{
		SampleInfo: make([]byte, 0, capacity),
	}
}

// AddSampleInfo adds a sampleinfo info based on parameters provided.
// If no length field, don't update the sample field (typicall audio cbcs)
func (s *SaizBox) AddSampleInfo(iv []byte, subsamplePatterns []SubSamplePattern) {
	size := len(iv)
	if len(subsamplePatterns) > 0 {
		size += 2 + len(subsamplePatterns)*6
		s.SampleInfo = append(s.SampleInfo, byte(size))
	} else if size > 0 {
		switch s.DefaultSampleInfoSize {
		case 0:
			s.DefaultSampleInfoSize = byte(size)
		default:
			if byte(size) != s.DefaultSampleInfoSize {
				panic("inconsistent sample info size")
			}
		}
	}
	if size > 0 {
		s.SampleCount++
	}
}

// Type - return box type
func (s *SaizBox) Type() string {
	return "saiz"
}

// Size - return calculated size
func (s *SaizBox) Size() uint64 {
	return s.expectedSize()
}

// expectedSize - calculate size based on flags and sample count
func (s *SaizBox) expectedSize() uint64 {
	size := uint64(boxHeaderSize + 9) // 9 = version + flags(4) + defaultSampleInfoSize(1) + sampleCount(4)
	if s.Flags&0x01 != 0 {
		size += 8 // auxInfoType(4) + auxInfoTypeParameter(4)
	}
	if s.DefaultSampleInfoSize == 0 {
		size += uint64(s.SampleCount) // 1 byte per sample info when default size is 0
	}
	return size
}

// Encode - write box to w
func (s *SaizBox) Encode(w io.Writer) error {
	sw := bits.NewFixedSliceWriter(int(s.Size()))
	err := s.EncodeSW(sw)
	if err != nil {
		return err
	}
	_, err = w.Write(sw.Bytes())
	return err
}

// EncodeSW - box-specific encode to slicewriter
func (s *SaizBox) EncodeSW(sw bits.SliceWriter) error {
	err := EncodeHeaderSW(s, sw)
	if err != nil {
		return err
	}
	versionAndFlags := (uint32(s.Version) << 24) + s.Flags
	sw.WriteUint32(versionAndFlags)
	if s.Flags&0x01 != 0 {
		sw.WriteString(s.AuxInfoType, false)
		sw.WriteUint32(s.AuxInfoTypeParameter)
	}
	sw.WriteUint8(s.DefaultSampleInfoSize)
	sw.WriteUint32(s.SampleCount)
	if s.DefaultSampleInfoSize == 0 {
		for i := uint32(0); i < s.SampleCount; i++ {
			sw.WriteUint8(s.SampleInfo[i])
		}
	}
	return sw.AccError()
}

// Info - write SaizBox details. Get sampleInfo list with level >= 1
func (s *SaizBox) Info(w io.Writer, specificBoxLevels, indent, indentStep string) (err error) {
	bd := newInfoDumper(w, indent, s, int(s.Version), s.Flags)
	if s.Flags&0x01 != 0 {
		bd.writef(" - auxInfoType: %s", s.AuxInfoType)
		bd.writef(" - auxInfoTypeParameter: %d", s.AuxInfoTypeParameter)
	}
	bd.writef(" - defaultSampleInfoSize: %d", s.DefaultSampleInfoSize)
	bd.writef(" - sampleCount: %d", s.SampleCount)
	level := getInfoLevel(s, specificBoxLevels)
	if level > 0 {
		if s.DefaultSampleInfoSize == 0 {
			for i := uint32(0); i < s.SampleCount; i++ {
				bd.writef(" - sampleInfo[%d]=%d", i+1, s.SampleInfo[i])
			}
		}
	}
	return bd.err
}
