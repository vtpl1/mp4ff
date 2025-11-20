package mp4

import (
	"bytes"
	"fmt"
	"io"

	"github.com/Eyevinn/mp4ff/bits"
)

// https://www.rfc-editor.org/rfc/rfc9639.html#name-streaminfo
type StreamInfo struct {
	BlockSizeMin  uint16   //u(16) 	The minimum block size (in samples) used in the stream, excluding the last block.
	BlockSizeMax  uint16   //u(16) 	The maximum block size (in samples) used in the stream.
	FrameSizeMin  uint32   //u(24) 	The minimum frame size (in bytes) used in the stream. 0 => value is not known.
	FrameSizeMax  uint32   //u(24) 	The maximum frame size (in bytes) used in the stream. 0 => value is not known.
	SampleRate    uint32   //u(20) 	Sample rate in Hz.
	Channels      uint8    //u(3) 	(number of channels)-1. FLAC supports from 1 to 8 channels.
	BitsPerSample uint8    //u(5) 	(bits per sample)-1. FLAC supports from 4 to 32 bits per sample.
	TotalSamples  uint64   //u(36) 	Total number of interchannel samples in the stream. 0 => no of total samples is unknown.
	MD5           [16]byte //u(128) MD5 checksum of the unencoded audio data. 0 => value is not known.
}

type Dfla struct {
	// data       []byte
	streamInfo StreamInfo
}

func CreateDfla(sampleRate uint32) *Dfla {
	return &Dfla{
		streamInfo: StreamInfo{
			BlockSizeMin:  32768,
			BlockSizeMax:  32768,
			FrameSizeMin:  0,
			FrameSizeMax:  0,
			Channels:      1,
			BitsPerSample: 16,
		},
	}

}

// DecodeDfla - box-specific decode
func DecodeDfla(hdr BoxHeader, startPos uint64, r io.Reader) (Box, error) {
	data, err := readBoxBody(r, hdr)
	if err != nil {
		return nil, err
	}
	sr := bits.NewFixedSliceReader(data)
	return DecodeDflaSR(hdr, startPos, sr)
}

// DecodeDflaSR - box-specific decode
func DecodeDflaSR(hdr BoxHeader, startPos uint64, sr bits.SliceReader) (Box, error) {
	data := sr.ReadBytes(hdr.payloadLen())
	rd := bytes.NewReader(data)
	br := bits.NewReader(rd)
	br.Read(64)
	blockSizeMin := br.Read(16)
	blockSizeMax := br.Read(16)
	frameSizeMin := br.Read(24)
	frameSizeMax := br.Read(24)
	sampleRate := br.Read(20)
	channels := br.Read(3) + 1
	bitsPerSample := br.Read(5) + 1
	totalSamples := br.Read(36)
	d := br.ReadRemainingBytes()
	// MD5
	var md5 [16]byte
	copy(md5[:], d[:])

	streamInfo := StreamInfo{
		BlockSizeMin:  uint16(blockSizeMin),
		BlockSizeMax:  uint16(blockSizeMax),
		FrameSizeMin:  uint32(frameSizeMin),
		FrameSizeMax:  uint32(frameSizeMax),
		SampleRate:    uint32(sampleRate),
		Channels:      uint8(channels),
		BitsPerSample: uint8(bitsPerSample),
		TotalSamples:  uint64(totalSamples),
		MD5:           md5,
	}

	return &Dfla{streamInfo: streamInfo}, sr.AccError()
}

// Type implements Box.
func (b *Dfla) Type() string {
	return "dfLa"
}

// Size implements Box.
func (b *Dfla) Size() uint64 {
	return uint64(boxHeaderSize + 42)
}

// Encode implements Box.
func (b *Dfla) Encode(w io.Writer) error {
	sw := bits.NewFixedSliceWriter(int(b.Size()))
	err := b.EncodeSW(sw)
	if err != nil {
		return err
	}
	_, err = w.Write(sw.Bytes())
	return err
}

// EncodeSW implements Box.
func (b *Dfla) EncodeSW(sw bits.SliceWriter) error {
	err := EncodeHeaderSW(b, sw)
	if err != nil {
		return err
	}
	si := b.streamInfo

	// ──────────────────────────────────────────────────────────
	// STREAMINFO block: 42 bytes
	// ──────────────────────────────────────────────────────────

	// 1. sync: first 8 bytes must be 0
	sw.WriteZeroBytes(8)

	// 2. Block sizes
	sw.WriteUint16(si.BlockSizeMin)
	sw.WriteUint16(si.BlockSizeMax)

	// 3. Frame sizes (24 bits each)
	sw.WriteBits(uint(si.FrameSizeMin), 24)
	sw.WriteBits(uint(si.FrameSizeMax), 24)

	// 4. Sample rate (20 bits)
	sw.WriteBits(uint(si.SampleRate), 20)

	// 5. Channels-1 (3 bits)
	if si.Channels == 0 {
		return fmt.Errorf("channels may not be zero")
	}
	sw.WriteBits(uint(si.Channels-1), 3)

	// 6. BitsPerSample-1 (5 bits)
	if si.BitsPerSample < 4 {
		return fmt.Errorf("bitsPerSample must be >= 4")
	}
	sw.WriteBits(uint(si.BitsPerSample-1), 5)

	// 7. TotalSamples (36 bits)
	sw.WriteBits(uint(si.TotalSamples), 36)

	// 8. MD5 (128 bits)
	sw.WriteBytes(si.MD5[:])

	return sw.AccError()
}

// Info implements Box.
func (b *Dfla) Info(w io.Writer, specificBoxLevels string, indent string, indentStep string) error {
	bd := newInfoDumper(w, indent, b, -1, 0)
	si := b.streamInfo
	bd.write(" - StreamInfo")
	bd.write(" 		- BlockSize: %d..%d", si.BlockSizeMin, si.BlockSizeMax)
	bd.write(" 		- FrameSize: %d..%d", si.FrameSizeMin, si.FrameSizeMax)
	bd.write(" 		- SampleRate: %d Hz", si.SampleRate)
	bd.write(" 		- Channels: %d", si.Channels)
	bd.write(" 		- BitsPerSample: %d", si.BitsPerSample)
	bd.write(" 		- TotalSamples: %d", si.TotalSamples)
	bd.write(" 		- MD5 (%d): %x", len(si.MD5), si.MD5)
	return bd.err
}
