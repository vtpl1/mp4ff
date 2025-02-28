package mp4

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/vtpl1/mp4ff/bits"
)

// AC3SampleRates - Sample rates as defined in  ETSI TS 102 366 V1.4.1 (2017) section 4.4.1.3
// Signaled in fscod - Sample rate code - 2 bits
var AC3SampleRates = []int{48000, 44100, 32000}

// AX3acmodChanneTable - channel configurations from ETSI TS 102 366 V1.4.1 (2017) section 4.4.2.3A
// Signaled in acmod - audio coding mode - 3 bits
var AC3acmodChannelTable = []string{
	"L/R", // Ch1 Ch2 dual mono but name them L R
	"C",
	"L/R",
	"L/C/R",
	"L/R/Cs",
	"L/C/R/Cs",
	"L/R/Ls/Rs",
	"L/C/R/Ls/Rs",
}

// AC3BitrateCodesKbps - Bitrates in kbps ETSI TS 102 366 V1.4.1 Table F.4.1 (2017)
var AC3BitrateCodesKbps = []uint16{
	32,
	40,
	48,
	56,
	64,
	80,
	96,
	112,
	128,
	160,
	192,
	224,
	256,
	320,
	384,
	448,
	512,
	576,
	640,
}

// Dac3Box - AC3SpecificBox from ETSI TS 102 366 V1.4.1 F.4 (2017)
// Extra b
type Dac3Box struct {
	FSCod         byte
	BSID          byte
	BSMod         byte
	ACMod         byte
	LFEOn         byte
	BitRateCode   byte
	Reserved      byte
	InitialZeroes byte // Should be zero
}

// DecodeDac3 - box-specific decode
func DecodeDac3(hdr BoxHeader, startPos uint64, r io.Reader) (Box, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return decodeDac3FromData(data)
}

// DecodeDac3SR - box-specific decode
func DecodeDac3SR(hdr BoxHeader, startPos uint64, sr bits.SliceReader) (Box, error) {
	data := sr.ReadBytes(hdr.payloadLen())
	if sr.AccError() != nil {
		return nil, sr.AccError()
	}
	return decodeDac3FromData(data)
}

func decodeDac3FromData(data []byte) (Box, error) {
	d := Dac3Box{}
	if len(data) > 3 {
		d.InitialZeroes = byte(len(data) - 3)
	}
	buf := bytes.NewBuffer(data)
	br := bits.NewReader(buf)
	for i := 0; i < int(d.InitialZeroes); i++ {
		if zero := br.ReadBits(8); zero != 0 {
			return nil, fmt.Errorf("dac3 box, extra initial bytes are not zero")
		}
	}
	d.FSCod = byte(br.ReadBits(2))
	d.BSID = byte(br.ReadBits(5))
	d.BSMod = byte(br.ReadBits(3))
	d.ACMod = byte(br.ReadBits(3))
	d.LFEOn = byte(br.ReadBits(1))
	d.BitRateCode = byte(br.ReadBits(5))
	// 5 bits reserved follows
	d.Reserved = byte(br.ReadBits(5))
	return &d, nil
}

// Type - box type
func (d *Dac3Box) Type() string {
	return "dac3"
}

// Size - calculated size of box
func (d *Dac3Box) Size() uint64 {
	return uint64(boxHeaderSize + 3 + uint(d.InitialZeroes))
}

// Encode - write box to w
func (d *Dac3Box) Encode(w io.Writer) error {
	sw := bits.NewFixedSliceWriter(int(d.Size()))
	err := d.EncodeSW(sw)
	if err != nil {
		return err
	}
	_, err = w.Write(sw.Bytes())
	return err
}

// Encode - write box to sw
func (d *Dac3Box) EncodeSW(sw bits.SliceWriter) error {
	err := EncodeHeaderSW(d, sw)
	if err != nil {
		return err
	}
	for i := 0; i < int(d.InitialZeroes); i++ {
		sw.WriteBits(0, 8)
	}
	sw.WriteBits(uint(d.FSCod), 2)
	sw.WriteBits(uint(d.BSID), 5)
	sw.WriteBits(uint(d.BSMod), 3)
	sw.WriteBits(uint(d.ACMod), 3)
	sw.WriteBits(uint(d.LFEOn), 1)
	sw.WriteBits(uint(d.BitRateCode), 5)
	sw.WriteBits(uint(d.Reserved), 5) // 5-bits reserved
	return sw.AccError()
}

// ChannelInfo - number of channels and channelmap according to E.1.3.1.8
func (d *Dac3Box) ChannelInfo() (nrChannels int, chanmap uint16) {
	speakers := GetChannelListFromACMod(d.ACMod)
	if d.LFEOn == 1 {
		speakers = append(speakers, "LFE")
	}
	nrChannels = len(speakers)
	for _, speaker := range speakers {
		chanmap |= CustomChannelMapLocations[speaker]
	}
	return nrChannels, chanmap
}

func (d *Dac3Box) Info(w io.Writer, specificBoxLevels, indent, indentStep string) error {
	bd := newInfoDumper(w, indent, d, -1, 0)
	bd.writef(" - sampleRateCode=%d => sampleRate=%d", d.FSCod, AC3SampleRates[d.FSCod])
	bd.writef(" - bitStreamInformation=%d", d.BSID)
	bd.writef(" - audioCodingMode=%d => channelConfiguration=%q", d.ACMod, AC3acmodChannelTable[d.ACMod])
	bd.writef(" - lowFrequencyEffectsChannelOn=%d", d.LFEOn)
	bd.writef(" - bitRateCode=%d => bitrate=%dkbps", d.BitRateCode, AC3BitrateCodesKbps[d.BitRateCode])
	nrChannels, chanmap := d.ChannelInfo()
	bd.writef(" - nrChannels=%d, chanmap=%04x", nrChannels, chanmap)
	if d.Reserved != 0 {
		bd.writef(" - reserved=%d", d.Reserved)
	}
	if d.InitialZeroes > 0 {
		bd.writef(" - weird initial zero bytes=%d", d.InitialZeroes)
	}
	return bd.err
}

func (d *Dac3Box) BitrateBps() int {
	return int(AC3BitrateCodesKbps[d.BitRateCode]) * 1000
}

func (d *Dac3Box) SamplingFrequency() int {
	return int(AC3SampleRates[d.FSCod])
}

// GetChannelListFromACMod - get list of channels from acmod byte
func GetChannelListFromACMod(acmod byte) []string {
	return strings.Split(AC3acmodChannelTable[acmod], "/")
}
