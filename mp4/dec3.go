package mp4

import (
	"bytes"
	"io"
	"strings"

	"github.com/vtpl1/mp4ff/bits"
)

// dec3

// ETSI TS 102 366 V1.4.1 (2017) Table E.1.4
// chanmap - Custom channel map - 16 bits
var CustomChannelMapLocations = map[string]uint16{
	"L":       1 << 15, // Left (MSB)
	"C":       1 << 14, // Center
	"R":       1 << 13, // Right
	"Ls":      1 << 12, // Left Surround
	"Rs":      1 << 11, // Right Surround
	"Lc/Rc":   1 << 10, // Front Left/Right of Center
	"Lrs/Rrs": 1 << 9,  // Left/Right Rear Surround
	"Cs":      1 << 8,  // Back Center
	"Ts":      1 << 7,  // Top Center
	"Lsd/Rsd": 1 << 6,  // Left/Right Surround Direct
	"Lw/Rw":   1 << 5,  // Left/Right Wide
	"Vhl/Vhr": 1 << 4,  // Top Front Left/Right
	"Vhc":     1 << 3,  // Top Front Center
	"Lts/Rts": 1 << 2,  // Left/Right Top Surround
	"LFE2":    1 << 1,  // Low Frequency 2
	"LFE":     1 << 0,  // Low Frequency
}

// EC3ChannelLocationBits - channel location signal in 9bits Table F.6.1
var EC3ChannelLocationBits = []string{
	"Lc/Rc",
	"Lrs/Rrs",
	"Cs",
	"Ts",
	"Lsd/Rsd",
	"Lw/Rw",
	"Lvh/Rvh",
	"Cvh",
	"LFE2", // MSB
}

// Dec3Box - AC3SpecificBox from ETSI TS 102 366 V1.4.1 F.4 (2017)
type Dec3Box struct {
	DataRate  uint16
	NumIndSub uint16
	EC3Subs   []EC3Sub
	Reserved  []byte
}

// EC3Sub - Enhanced AC-3 substream information
type EC3Sub struct {
	FSCod     byte
	BSID      byte
	ASVC      byte
	BSMod     byte
	ACMod     byte
	LFEOn     byte
	NumDepSub byte
	ChanLoc   uint16
}

// DecodeDec3 - box-specific decode
func DecodeDec3(hdr BoxHeader, startPos uint64, r io.Reader) (Box, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return decodeDec3FromData(data)
}

// DecodeDec3SR - box-specific decode
func DecodeDec3SR(hdr BoxHeader, startPos uint64, sr bits.SliceReader) (Box, error) {
	data := sr.ReadBytes(hdr.payloadLen())
	if sr.AccError() != nil {
		return nil, sr.AccError()
	}
	return decodeDec3FromData(data)
}

func decodeDec3FromData(data []byte) (Box, error) {
	buf := bytes.NewBuffer(data)
	br := bits.NewReader(buf)
	b := Dec3Box{}
	b.DataRate = uint16(br.ReadBits(13))
	nrSubs := br.ReadBits(3) + 1 // There must be one base stream
	for i := 0; i < int(nrSubs); i++ {
		es := EC3Sub{}
		es.FSCod = byte(br.ReadBits(2))
		es.BSID = byte(br.ReadBits(5))
		_ = br.ReadBits(1) // Reserved 0
		es.ASVC = byte(br.ReadBits(1))
		es.BSMod = byte(br.ReadBits(3))
		es.ACMod = byte(br.ReadBits(3))
		es.LFEOn = byte(br.ReadBits(1))
		_ = br.ReadBits(3) // Reserved 000
		es.NumDepSub = byte(br.ReadBits(4))
		if es.NumDepSub > 0 {
			es.ChanLoc = uint16(br.ReadBits(9))
		} else {
			_ = br.ReadBits(1) // Reserved 0
		}
		if br.AccError() != nil {
			return nil, br.AccError()
		}
		b.EC3Subs = append(b.EC3Subs, es)
	}
	b.Reserved = br.ReadRemainingBytes()
	return &b, br.AccError()
}

// Type - box type
func (d *Dec3Box) Type() string {
	return "dec3"
}

// Size - calculated size of box
func (d *Dec3Box) Size() uint64 {
	size := boxHeaderSize + 2
	for _, es := range d.EC3Subs {
		size += 3
		if es.NumDepSub > 0 {
			size += 1
		}
	}
	size += len(d.Reserved)
	return uint64(size)
}

// Encode - write box to w
func (d *Dec3Box) Encode(w io.Writer) error {
	sw := bits.NewFixedSliceWriter(int(d.Size()))
	err := d.EncodeSW(sw)
	if err != nil {
		return err
	}
	_, err = w.Write(sw.Bytes())
	return err
}

// EncodeSW - write box to sw
func (d *Dec3Box) EncodeSW(sw bits.SliceWriter) error {
	err := EncodeHeaderSW(d, sw)
	if err != nil {
		return err
	}
	sw.WriteBits(uint(d.DataRate), 13)
	sw.WriteBits(uint(len(d.EC3Subs))-1, 3)
	for _, es := range d.EC3Subs {
		sw.WriteBits(uint(es.FSCod), 2)
		sw.WriteBits(uint(es.BSID), 5)
		sw.WriteBits(0, 1) // reserved 0
		sw.WriteBits(uint(es.ASVC), 1)
		sw.WriteBits(uint(es.BSMod), 3)
		sw.WriteBits(uint(es.ACMod), 3)
		sw.WriteBits(uint(es.LFEOn), 1)
		sw.WriteBits(0, 3) // reserved 000
		sw.WriteBits(uint(es.NumDepSub), 4)
		if es.NumDepSub > 0 {
			sw.WriteBits(uint(es.ChanLoc), 9)
		} else {
			sw.WriteBits(0, 1) // Reserved 0d
		}
	}
	if len(d.Reserved) > 0 {
		sw.WriteBytes(d.Reserved)
	}
	return sw.AccError()
}

func (d *Dec3Box) Info(w io.Writer, specificBoxLevels, indent, indentStep string) error {
	bd := newInfoDumper(w, indent, d, -1, 0)
	bd.writef(" - bitrate=%dkbps", d.DataRate)
	fscod := d.EC3Subs[0].FSCod
	bd.writef(" - sampleRateCode=%d => sampleRate=%d", fscod, AC3SampleRates[fscod])
	nrChannels, chanmap := d.ChannelInfo()
	bd.writef(" - nrChannels=%d, chanmap=%04x", nrChannels, chanmap)
	bd.writef(" - nrSubstreams=%d", len(d.EC3Subs))
	for i, es := range d.EC3Subs {
		bd.writef("   - %d fscod=%d bsid=%d asvc=%d bsmod=%d acmod=%d lfeon=%d num_dep_sub=%d chan_loc=%x",
			i+1, es.FSCod, es.BSID, es.ASVC, es.BSMod, es.ACMod, es.LFEOn, es.NumDepSub, es.ChanLoc)
	}
	return bd.err
}

func (d *Dec3Box) ChannelInfo() (nrChannels int, chanmap uint16) {
	// All Enhanced AC-3 bit streams shall contain an independent substream
	// assigned substream ID 0 (E.1.3.1.2)
	substream := d.EC3Subs[0]

	// Get base channel configuration according to acmod
	channels := GetChannelListFromACMod(substream.ACMod)
	if substream.LFEOn == 1 {
		channels = append(channels, "LFE")
	}

	// Dependent substreams associated with this independent substream
	if substream.NumDepSub > 0 {
		for i := 0; i < 9; i++ {
			if substream.ChanLoc&(1<<i) != 0 {
				channels = append(channels, EC3ChannelLocationBits[i])
			}
		}
	}
	for _, channel := range channels {
		// Check if a channel pair (contains /)
		if strings.Contains(channel, "/") {
			nrChannels += 2
		} else {
			nrChannels++
		}

		chanmap |= CustomChannelMapLocations[channel]
	}

	return nrChannels, chanmap
}
