package bits

import (
	"encoding/binary"
	"errors"
	"io"
)

// ESBPReader errors.
var (
	ErrNotReadSeeker                 = errors.New("reader does not support Seek")
	ErrRbspTrailingBitsDonTStartWith = errors.New("rbspTrailingBits don't start with 1")
	ErrAnotherRbspTrailingBits       = errors.New("another 1 in RbspTrailingBits")
)

const startCodeEmulationPreventionByte = 0x03

// EBSPReader reads an EBSP bitstream dropping start-code emulation bytes.
// It also supports checking for more rbsp data and reading rbsp_trailing_bits.
type EBSPReader struct {
	rd        io.Reader
	err       error
	n         int  // current number of bits
	v         uint // current accumulated value
	pos       int
	zeroCount int // Count number of zero bytes read
}

// NewEBSPReader return a new EBSP reader stopping reading at first error.
func NewEBSPReader(rd io.Reader) *EBSPReader {
	return &EBSPReader{
		rd:        rd,
		err:       nil,
		n:         0,
		v:         0,
		pos:       -1,
		zeroCount: 0,
	}
}

// AccError implements SliceReader - returns the accumulated error. If no error, returns nil.
func (e *EBSPReader) AccError() error {
	return e.err
}

// GetPos implements SliceReader.
func (e *EBSPReader) GetPos() int {
	panic("unimplemented")
}

// Length implements SliceReader.
func (e *EBSPReader) Length() int {
	panic("unimplemented")
}

// LookAhead implements SliceReader.
func (e *EBSPReader) LookAhead(offset int, data []byte) error { //nolint:revive
	panic("unimplemented")
}

// NrRemainingBytes implements SliceReader.
func (e *EBSPReader) NrRemainingBytes() int {
	panic("unimplemented")
}

// NrBytesRead returns how many bytes read into parser.
func (e *EBSPReader) NrBytesRead() int {
	return e.pos + 1 // Starts at -1
}

// NrBitsRead returns total number of bits read into parser.
func (e *EBSPReader) NrBitsRead() int {
	nrBits := e.NrBytesRead() * 8
	if e.NrBitsReadInCurrentByte() != 8 {
		nrBits += e.NrBitsReadInCurrentByte() - 8
	}
	return nrBits
}

// NrBitsReadInCurrentByte returns number of bits read in current byte.
func (e *EBSPReader) NrBitsReadInCurrentByte() int {
	return 8 - e.n
}

// ReadBits implements SliceReader - reads n bits and respects and accumulates errors. If error, returns 0.
func (e *EBSPReader) ReadBits(n int) uint {
	if e.err != nil {
		return 0
	}
	var err error
	for e.n < n {
		e.v <<= 8
		var b uint8
		err = binary.Read(e.rd, binary.BigEndian, &b)
		if err != nil {
			e.err = err
			return 0
		}
		e.pos++
		if e.zeroCount == 2 && b == startCodeEmulationPreventionByte {
			err = binary.Read(e.rd, binary.BigEndian, &b)
			if err != nil {
				e.err = err
				return 0
			}
			e.pos++
			e.zeroCount = 0
		}
		if b != 0 {
			e.zeroCount = 0
		} else {
			e.zeroCount++
		}
		e.v |= uint(b)

		e.n += 8
	}
	v := e.v >> uint(e.n-n) //nolint:gosec

	e.n -= n
	e.v &= Mask(e.n)

	return v
}

// ReadBytes implements SliceReader - read n bytes and return nil if new or accumulated error.
func (e *EBSPReader) ReadBytes(n int) []byte {
	if e.err != nil {
		return nil
	}
	payload := make([]byte, n)
	for i := range n {
		b := byte(e.ReadBits(8))
		payload[i] = b
	}
	if e.err != nil {
		return nil
	}
	return payload
}

// ReadFixedLengthString implements SliceReader.
func (e *EBSPReader) ReadFixedLengthString(n int) string { //nolint:revive
	panic("unimplemented")
}

// ReadInt16 implements SliceReader.
func (e *EBSPReader) ReadInt16() int16 {
	panic("unimplemented")
}

// ReadInt32 implements SliceReader.
func (e *EBSPReader) ReadInt32() int32 {
	panic("unimplemented")
}

// ReadInt64 implements SliceReader.
func (e *EBSPReader) ReadInt64() int64 {
	panic("unimplemented")
}

// ReadPossiblyZeroTerminatedString implements SliceReader.
func (e *EBSPReader) ReadPossiblyZeroTerminatedString(maxLen int) (str string, ok bool) { //nolint:nonamedreturns,revive
	panic("unimplemented")
}

// ReadUint16 implements SliceReader.
func (e *EBSPReader) ReadUint16() uint16 {
	panic("unimplemented")
}

// ReadUint24 implements SliceReader.
func (e *EBSPReader) ReadUint24() uint32 {
	panic("unimplemented")
}

// ReadUint32 implements SliceReader.
func (e *EBSPReader) ReadUint32() uint32 {
	panic("unimplemented")
}

// ReadUint64 implements SliceReader.
func (e *EBSPReader) ReadUint64() uint64 {
	panic("unimplemented")
}

// ReadUint8 implements SliceReader.
func (e *EBSPReader) ReadUint8() byte {
	panic("unimplemented")
}

// ReadZeroTerminatedString implements SliceReader.
func (e *EBSPReader) ReadZeroTerminatedString(maxLen int) string { //nolint:revive
	panic("unimplemented")
}

// RemainingBytes implements SliceReader.
func (e *EBSPReader) RemainingBytes() []byte {
	panic("unimplemented")
}

// SetPos implements SliceReader.
func (e *EBSPReader) SetPos(pos int) { //nolint:revive
	panic("unimplemented")
}

// SkipBytes implements SliceReader.
func (e *EBSPReader) SkipBytes(n int) { //nolint:revive
	panic("unimplemented")
}

// ReadFlag reads 1 bit and translates a bool.
func (e *EBSPReader) ReadFlag() bool {
	return e.ReadBits(1) == 1
}

// ReadExpGolomb reads one unsigned exponential Golomb code.
func (e *EBSPReader) ReadExpGolomb() uint {
	if e.err != nil {
		return 0
	}
	leadingZeroBits := 0
	for {
		b := e.ReadBits(1)
		if e.err != nil {
			return 0
		}
		if b == 1 {
			break
		}
		leadingZeroBits++
	}
	var res uint = (1 << leadingZeroBits) - 1
	endBits := e.ReadBits(leadingZeroBits)
	if e.err != nil {
		return 0
	}
	return res + endBits
}

// ReadSignedGolomb reads one signed exponential Golomb code.
func (e *EBSPReader) ReadSignedGolomb() int {
	if e.err != nil {
		return 0
	}
	unsignedGolomb := e.ReadExpGolomb()
	if e.err != nil {
		return 0
	}
	if unsignedGolomb%2 == 1 {
		return int((unsignedGolomb + 1) / 2) //nolint:gosec
	}
	return -int(unsignedGolomb / 2) //nolint:gosec
}

// IsSeeker returns tru if underluing reader supports Seek interface.
func (e *EBSPReader) IsSeeker() bool {
	_, ok := e.rd.(io.ReadSeeker)
	return ok
}

// MoreRbspData returns false if next bit is 1 and last 1-bit in fullSlice.
// Underlying reader must support ReadSeeker interface to reset after check.
// Return false, nil if underlying error.
func (e *EBSPReader) MoreRbspData() (bool, error) {
	if !e.IsSeeker() {
		return false, ErrNotReadSeeker
	}
	// Find out if next position is the last 1
	stateCopy := *e

	firstBit := e.ReadBits(1)
	if e.err != nil {
		return false, nil
	}
	if firstBit != 1 {
		err := e.reset(stateCopy)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	// If all remainging bits are zero, there is no more rbsp data
	more := false
	for {
		b := e.ReadBits(1)
		if errors.Is(e.err, io.EOF) {
			e.err = nil // Reset
			break
		}
		if e.err != nil {
			return false, nil
		}
		if b == 1 {
			more = true
			break
		}
	}
	err := e.reset(stateCopy)
	if err != nil {
		return false, err
	}
	return more, nil
}

// reset resets EBSPReader based on copy of previous state.
func (e *EBSPReader) reset(prevState EBSPReader) error {
	rdSeek, _ := e.rd.(io.ReadSeeker)
	_, err := rdSeek.Seek(int64(prevState.pos+1), 0)
	if err != nil {
		return err
	}
	e.n = prevState.n
	e.v = prevState.v
	e.pos = prevState.pos
	e.zeroCount = prevState.zeroCount
	return nil
}

// ReadRbspTrailingBits reads rbsp_traling_bits. Returns error if wrong pattern.
// If other error, returns nil and let AccError() provide that error.
func (e *EBSPReader) ReadRbspTrailingBits() error {
	if e.err != nil {
		return nil
	}
	firstBit := e.ReadBits(1)
	if e.err != nil {
		return nil
	}
	if firstBit != 1 {
		return ErrRbspTrailingBitsDonTStartWith
	}
	for {
		b := e.ReadBits(1)
		if errors.Is(e.err, io.EOF) {
			e.err = nil // Reset
			return nil
		}
		if e.err != nil {
			return nil
		}
		if b == 1 {
			return ErrAnotherRbspTrailingBits
		}
	}
}

// SetError sets an error if not already set.
func (e *EBSPReader) SetError(err error) {
	if e.err == nil {
		e.err = err
	}
}
