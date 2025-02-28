package mp4

import (
	"encoding/hex"
	"io"

	"github.com/vtpl1/mp4ff/bits"
)

// TencBox - Track Encryption Box
// Defined in ISO/IEC 23001-7 Secion 8.2
type TencBox struct {
	Version                byte
	Flags                  uint32
	DefaultCryptByteBlock  byte
	DefaultSkipByteBlock   byte
	DefaultIsProtected     byte
	DefaultPerSampleIVSize byte
	DefaultKID             UUIDType
	// DefaultConstantIVSize  byte given by len(DefaultConstantIV)
	DefaultConstantIV []byte
}

// DecodeTenc - box-specific decode
func DecodeTenc(hdr BoxHeader, startPos uint64, r io.Reader) (Box, error) {
	data, err := readBoxBody(r, hdr)
	if err != nil {
		return nil, err
	}
	sr := bits.NewFixedSliceReader(data)
	return DecodeTencSR(hdr, startPos, sr)
}

// DecodeTencSR - box-specific decode
func DecodeTencSR(hdr BoxHeader, startPos uint64, sr bits.SliceReader) (Box, error) {
	versionAndFlags := sr.ReadUint32()
	version := byte(versionAndFlags >> 24)

	b := TencBox{
		Version: version,
		Flags:   versionAndFlags & flagsMask,
	}
	_ = sr.ReadUint8() // Skip reserved == 0
	if version == 0 {
		_ = sr.ReadUint8() // Skip reserved == 0
	} else {
		infoByte := sr.ReadUint8()
		b.DefaultCryptByteBlock = infoByte >> 4
		b.DefaultSkipByteBlock = infoByte & 0x0f
	}
	b.DefaultIsProtected = sr.ReadUint8()
	b.DefaultPerSampleIVSize = sr.ReadUint8()
	b.DefaultKID = UUIDType(sr.ReadBytes(16))
	if b.DefaultIsProtected == 1 && b.DefaultPerSampleIVSize == 0 {
		defaultConstantIVSize := int(sr.ReadUint8())
		b.DefaultConstantIV = sr.ReadBytes(defaultConstantIVSize)
	}
	return &b, sr.AccError()
}

// Type - return box type
func (t *TencBox) Type() string {
	return "tenc"
}

// Size - return calculated size
func (t *TencBox) Size() uint64 {
	var size uint64 = 32
	if t.DefaultIsProtected == 1 && t.DefaultPerSampleIVSize == 0 {
		size += uint64(1 + len(t.DefaultConstantIV))
	}
	return size
}

// Encode - write box to w
func (t *TencBox) Encode(w io.Writer) error {
	sw := bits.NewFixedSliceWriter(int(t.Size()))
	err := t.EncodeSW(sw)
	if err != nil {
		return err
	}
	_, err = w.Write(sw.Bytes())
	return err
}

// EncodeSW - box-specific encode to slicewriter
func (t *TencBox) EncodeSW(sw bits.SliceWriter) error {
	err := EncodeHeaderSW(t, sw)
	if err != nil {
		return err
	}
	versionAndFlags := (uint32(t.Version) << 24) + t.Flags
	sw.WriteUint32(versionAndFlags)
	sw.WriteUint8(0) // reserved
	if t.Version == 0 {
		sw.WriteUint8(0) // reserved
	} else {
		sw.WriteUint8(t.DefaultCryptByteBlock<<4 | t.DefaultSkipByteBlock)
	}
	sw.WriteUint8(t.DefaultIsProtected)
	sw.WriteUint8(t.DefaultPerSampleIVSize)
	sw.WriteBytes(t.DefaultKID)
	if t.DefaultIsProtected == 1 && t.DefaultPerSampleIVSize == 0 {
		sw.WriteUint8(byte(len(t.DefaultConstantIV)))
		sw.WriteBytes(t.DefaultConstantIV)
	}
	return sw.AccError()
}

// Info - write box info to w
func (t *TencBox) Info(w io.Writer, specificBoxLevels, indent, indentStep string) (err error) {
	bd := newInfoDumper(w, indent, t, int(t.Version), t.Flags)
	if t.Version > 0 {
		bd.writef(" - defaultCryptByteBlock: %d", t.DefaultCryptByteBlock)
		bd.writef(" - defaultSkipByteBlock: %d", t.DefaultSkipByteBlock)
	}
	bd.writef(" - defaultIsProtected: %d", t.DefaultIsProtected)
	bd.writef(" - defaultPerSampleIVSize: %d", t.DefaultPerSampleIVSize)
	bd.writef(" - defaultKID: %s", t.DefaultKID)
	if t.DefaultIsProtected == 1 && t.DefaultPerSampleIVSize == 0 {
		bd.writef(" - defaultConstantIV: %s", hex.EncodeToString(t.DefaultConstantIV))
	}
	return bd.err
}
