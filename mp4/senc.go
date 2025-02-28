package mp4

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/vtpl1/mp4ff/internal/bits"
)

var (
	ErrSencBoxAlreadyParsed = errors.New("senc box already parsed")
	ErrMixOfIvLengths       = errors.New("mix of IV lengths")
	ErrCouldNotDecodeSenc   = errors.New("could not decode senc")
)

// SencBox - Sample Encryption Box (senc) (in trak or traf box)
// Should only be decoded after saio and saiz provide relevant offset and sizes
// Here we make a two-step decode, with first step reading, and other parsing.
// See ISO/IEC 23001-7 Section 7.2 and CMAF specification
// Full Box + SampleCount.
type SencBox struct {
	Version          byte
	ReadButNotParsed bool
	PerSampleIVSize  byte
	Flags            uint32
	SampleCount      uint32
	StartPos         uint64
	rawData          []byte                 // intermediate storage when reading
	IVs              []InitializationVector // 8 or 16 bytes if present
	SubSamples       [][]SubSamplePattern
	readBoxSize      uint64 // As read from box header
}

// Encode implements Box.
func (s *SencBox) Encode(w io.Writer) error {
	// First check if subsamplencryption is to be used since it influences the box size
	s.setSubSamplesUsedFlag()
	sw := bits.NewFixedSliceWriter(int(s.Size())) //nolint:gosec
	err := s.EncodeSW(sw)
	if err != nil {
		return err
	}
	_, err = w.Write(sw.Bytes())
	return err
}

// EncodeSW implements Box.
func (s *SencBox) EncodeSW(sw bits.SliceWriter) error {
	s.setSubSamplesUsedFlag()
	err := EncodeHeaderSW(s, sw)
	if err != nil {
		return err
	}
	err = s.EncodeSWNoHdr(sw)
	return err
}

// EncodeSWNoHdr encodes without header (useful for PIFF box).
func (s *SencBox) EncodeSWNoHdr(sw bits.SliceWriter) error {
	versionAndFlags := (uint32(s.Version) << 24) + s.Flags
	sw.WriteUint32(versionAndFlags)
	sw.WriteUint32(s.SampleCount)
	if s.ReadButNotParsed {
		sw.WriteBytes(s.rawData)
		return sw.AccError()
	}
	perSampleIVSize := s.GetPerSampleIVSize()
	for i := range s.SampleCount {
		if perSampleIVSize > 0 {
			sw.WriteBytes(s.IVs[i])
		}
		if s.Flags&UseSubSampleEncryption != 0 {
			sw.WriteUint16(uint16(len(s.SubSamples[i]))) //nolint:gosec
			for _, subSample := range s.SubSamples[i] {
				sw.WriteUint16(subSample.BytesOfClearData)
				sw.WriteUint32(subSample.BytesOfProtectedData)
			}
		}
	}
	return sw.AccError()
}

// Info implements Box.
func (s *SencBox) Info(w io.Writer, specificBoxLevels string, indent string, _ string) error {
	bd := newInfoDumper(w, indent, s, int(s.Version), s.Flags)
	bd.writef(" - sampleCount: %d", s.SampleCount)
	if s.ReadButNotParsed {
		bd.writef(" - NOT YET PARSED, call ParseReadBox to parse it")
		return nil
	}
	for _, subSamples := range s.SubSamples {
		if len(subSamples) > 0 {
			s.Flags |= UseSubSampleEncryption
		}
	}
	perSampleIVSize := s.GetPerSampleIVSize()
	bd.writef(" - perSampleIVSize: %d", perSampleIVSize)
	level := getInfoLevel(s, specificBoxLevels)
	if level > 0 && (perSampleIVSize > 0 || s.Flags&UseSubSampleEncryption != 0) {
		for i := range s.SampleCount {
			line := fmt.Sprintf(" - sample[%d]:", i+1)
			if perSampleIVSize > 0 {
				line += " iv=" + hex.EncodeToString(s.IVs[i])
			}
			bd.writef(line)
			if s.Flags&UseSubSampleEncryption != 0 {
				for j, subSample := range s.SubSamples[i] {
					bd.writef("   - subSample[%d]: nrBytesClear=%d nrBytesProtected=%d", j+1,
						subSample.BytesOfClearData, subSample.BytesOfProtectedData)
				}
			}
		}
	}
	return bd.err
}

// GetPerSampleIVSize - return perSampleIVSize.
func (s *SencBox) GetPerSampleIVSize() int {
	return int(s.PerSampleIVSize)
}

// Size implements Box.
func (s *SencBox) Size() uint64 {
	if s.readBoxSize > 0 {
		return s.readBoxSize
	}
	return s.calcSize()
}

func (s *SencBox) calcSize() uint64 {
	totalSize := uint64(boxHeaderSize + 8)
	perSampleIVSize := uint64(s.GetPerSampleIVSize()) //nolint:gosec
	for i := range s.SampleCount {
		totalSize += perSampleIVSize
		if s.Flags&UseSubSampleEncryption != 0 {
			totalSize += 2 + 6*uint64(len(s.SubSamples[i]))
		}
	}
	return totalSize
}

// Type implements Box.
func (s *SencBox) Type() string {
	return "senc"
}

// UseSubSampleEncryption - flag for subsample encryption.
const UseSubSampleEncryption = 0x2

// SubSamplePattern - pattern of subsample encryption.
type SubSamplePattern struct {
	BytesOfClearData     uint16
	BytesOfProtectedData uint32
}

// InitializationVector (8 or 16 bytes).
type InitializationVector []byte

// SencSample - sample in SencBox.
type SencSample struct {
	IV         InitializationVector // 0,8,16 byte length
	SubSamples []SubSamplePattern
}

// CreateSencBox - create an empty SencBox.
func CreateSencBox() *SencBox {
	return &SencBox{} //nolint:exhaustruct
}

// NewSencBox returns a SencBox with capacity for IVs and SubSamples.
func NewSencBox(ivCapacity, subSampleCapacity int) *SencBox {
	s := SencBox{} //nolint:exhaustruct
	if ivCapacity > 0 {
		s.IVs = make([]InitializationVector, 0, ivCapacity)
	}
	if subSampleCapacity > 0 {
		s.SubSamples = make([][]SubSamplePattern, 0, subSampleCapacity)
	}
	return &s
}

// AddSample - add a senc sample with possible IV and subsamples.
func (s *SencBox) AddSample(sample SencSample) error {
	if len(sample.IV) != 0 {
		if s.SampleCount == 0 {
			s.PerSampleIVSize = byte(len(sample.IV))
		} else if len(sample.IV) != int(s.PerSampleIVSize) {
			return ErrMixOfIvLengths
		}

		if len(sample.IV) != 0 {
			s.IVs = append(s.IVs, sample.IV)
		}
	}

	if len(sample.SubSamples) > 0 {
		s.SubSamples = append(s.SubSamples, sample.SubSamples)
		s.Flags |= UseSubSampleEncryption
	}
	s.SampleCount++
	return nil
}

// setSubSamplesUsedFlag - set flag if subsamples are used.
func (s *SencBox) setSubSamplesUsedFlag() {
	for _, subSamples := range s.SubSamples {
		if len(subSamples) > 0 {
			s.Flags |= UseSubSampleEncryption
			break
		}
	}
}

// ParseReadBox - second phase when perSampleIVSize should be known from tenc or sgpd boxes
// if perSampleIVSize is 0, we try to find the appropriate error given data length.
func (s *SencBox) ParseReadBox(perSampleIVSize byte, _ *SaizBox) error {
	if !s.ReadButNotParsed {
		return ErrSencBoxAlreadyParsed
	}
	if perSampleIVSize != 0 {
		s.PerSampleIVSize = perSampleIVSize
	}
	sr := bits.NewFixedSliceReader(s.rawData)
	nrBytesLeft := uint32(sr.NrRemainingBytes()) //nolint:gosec

	if s.Flags&UseSubSampleEncryption == 0 {
		// No subsamples
		if perSampleIVSize == 0 { // Infer the size
			perSampleIVSize = byte(nrBytesLeft / s.SampleCount)
			s.PerSampleIVSize = perSampleIVSize
		}

		s.IVs = make([]InitializationVector, 0, s.SampleCount)
		switch perSampleIVSize {
		case 0:
			// Nothing to do
		case 8:
			for range s.SampleCount {
				s.IVs = append(s.IVs, sr.ReadBytes(8))
			}
		case 16:
			for range s.SampleCount {
				s.IVs = append(s.IVs, sr.ReadBytes(16))
			}
		default:
			return fmt.Errorf("strange derived PerSampleIVSize: %d", perSampleIVSize) //nolint:err113
		}
		s.ReadButNotParsed = false
		return nil
	}
	// 6 bytes of subsamplecount per subsample and known perSampleIVSize
	// The total length for each sample should correspond to
	// sizes in saiz (defaultSampleInfoSize or SampleInfo value)
	// We don't check that though, but it could be implemented here.
	if perSampleIVSize != 0 {
		if ok := s.parseAndFillSamples(sr, perSampleIVSize); !ok {
			return fmt.Errorf("error decoding senc with perSampleIVSize = %d", perSampleIVSize) //nolint:err113
		}
		s.ReadButNotParsed = false
		return nil
	}

	// Finally, 6 bytes of subsamplecount per subsample and unknown perSampleIVSize
	startPos := sr.GetPos()
	ok := false
	for perSampleIVSize := byte(0); perSampleIVSize <= 16; perSampleIVSize += 8 {
		sr.SetPos(startPos)
		ok = s.parseAndFillSamples(sr, perSampleIVSize)
		if ok {
			break // We have found a working perSampleIVSize
		}
	}
	if !ok {
		return ErrCouldNotDecodeSenc
	}
	s.ReadButNotParsed = false
	return nil
}

// parseAndFillSamples - parse and fill senc samples given perSampleIVSize.
//
//nolint:nonamedreturns
func (s *SencBox) parseAndFillSamples(sr bits.SliceReader, perSampleIVSize byte) (ok bool) {
	ok = true
	s.SubSamples = make([][]SubSamplePattern, s.SampleCount)
	for i := range s.SampleCount {
		if perSampleIVSize > 0 {
			if sr.NrRemainingBytes() < int(perSampleIVSize) {
				ok = false
				break
			}
			s.IVs = append(s.IVs, sr.ReadBytes(int(perSampleIVSize)))
		}
		if sr.NrRemainingBytes() < 2 {
			ok = false
			break
		}
		subsampleCount := int(sr.ReadUint16())
		if sr.NrRemainingBytes() < subsampleCount*6 {
			ok = false
			break
		}
		s.SubSamples[i] = make([]SubSamplePattern, subsampleCount)
		for j := range subsampleCount {
			s.SubSamples[i][j].BytesOfClearData = sr.ReadUint16()
			s.SubSamples[i][j].BytesOfProtectedData = sr.ReadUint32()
		}
	}
	if !ok || sr.NrRemainingBytes() != 0 {
		// Cleanup the IVs and SubSamples which may have been partially set
		s.IVs = nil
		s.SubSamples = nil
		ok = false
	}
	s.PerSampleIVSize = perSampleIVSize
	return ok
}

// DecodeSenc - box-specific decode.
func DecodeSenc(hdr BoxHeader, startPos uint64, r io.Reader) (Box, error) {
	if hdr.Size < 16 {
		return nil, fmt.Errorf("box size %d less than min size 16", hdr.Size) //nolint:err113
	}
	data, err := readBoxBody(r, hdr)
	if err != nil {
		return nil, err
	}

	versionAndFlags := binary.BigEndian.Uint32(data[0:4])
	version := byte(versionAndFlags >> 24)
	flags := versionAndFlags & flagsMask
	if version > 0 {
		return nil, fmt.Errorf("version %d not supported", version) //nolint:err113
	}
	sampleCount := binary.BigEndian.Uint32(data[4:8])

	if len(data) < 8 {
		return nil, fmt.Errorf("senc: box size %d less than 16", hdr.Size) //nolint:err113
	}

	senc := SencBox{ //nolint:exhaustruct
		Version:          version,
		rawData:          data[8:], // After the first 8 bytes of box content
		Flags:            flags,
		StartPos:         startPos,
		SampleCount:      sampleCount,
		ReadButNotParsed: true,
		readBoxSize:      hdr.Size,
	}

	if flags&UseSubSampleEncryption != 0 && (len(senc.rawData) < 2*int(sampleCount)) {
		//nolint:err113
		return nil, fmt.Errorf("box size %d too small for %d samples and subSampleEncryption",
			hdr.Size, sampleCount)
	}

	if senc.SampleCount == 0 || len(senc.rawData) == 0 {
		senc.ReadButNotParsed = false
		return &senc, nil
	}
	return &senc, nil
}

// DecodeSencSR - box-specific decode.
func DecodeSencSR(hdr BoxHeader, startPos uint64, sr bits.SliceReader) (Box, error) {
	if hdr.Size < 16 {
		return nil, fmt.Errorf("box size %d less than min size 16", hdr.Size) //nolint:err113
	}

	versionAndFlags := sr.ReadUint32()
	version := byte(versionAndFlags >> 24)
	if version > 0 {
		return nil, fmt.Errorf("version %d not supported", version) //nolint:err113
	}
	flags := versionAndFlags & flagsMask
	sampleCount := sr.ReadUint32()

	if flags&UseSubSampleEncryption != 0 && ((hdr.Size - 16) < 2*uint64(sampleCount)) {
		//nolint:err113
		return nil, fmt.Errorf("box size %d too small for %d samples and subSampleEncryption",
			hdr.Size, sampleCount)
	}

	senc := SencBox{ //nolint:exhaustruct
		Version:          version,
		rawData:          sr.ReadBytes(hdr.payloadLen() - 8), // After the first 8 bytes of box content
		Flags:            flags,
		StartPos:         startPos,
		SampleCount:      sampleCount,
		ReadButNotParsed: true,
		readBoxSize:      hdr.Size,
	}

	if senc.SampleCount == 0 || len(senc.rawData) == 0 {
		senc.ReadButNotParsed = false
		return &senc, sr.AccError()
	}
	return &senc, sr.AccError()
}
