package bits

import (
	"encoding/binary"
	"errors"
)

var (
	ErrDidNotFindTerminatingZero          = errors.New("did not find terminating zero")
	ErrAttemptToReadNegativeNumberOfBytes = errors.New("attempt to read negative number of bytes")
	ErrAttemptToSkipBytesBeyondSliceLen   = errors.New("attempt to skip bytes beyond slice len")
	ErrOutOfBounds                        = errors.New("out of bounds")
)

// FixedSliceReader - read integers and other data from a fixed slice.
// Accumulates error, and the first error can be retrieved.
// If err != nil, 0 or empty string is returned.
type FixedSliceReader struct {
	err   error
	slice []byte
	pos   int
	len   int
}

// AccError implements SliceReader - get accumulated error after read operations.
func (f *FixedSliceReader) AccError() error {
	return f.err
}

// GetPos implements SliceReader - get read position is slice.
func (f *FixedSliceReader) GetPos() int {
	return f.pos
}

// Length implements SliceReader - get length of slice.
func (f *FixedSliceReader) Length() int {
	return f.len
}

// LookAhead implements SliceReader.
func (f *FixedSliceReader) LookAhead(offset int, data []byte) error {
	if f.pos+offset+len(data) > f.len {
		return ErrOutOfBounds
	}
	copy(data, f.slice[f.pos+offset:])
	return nil
}

// NrRemainingBytes implements SliceReader - return number of bytes remaining.
func (f *FixedSliceReader) NrRemainingBytes() int {
	if f.err != nil {
		return 0
	}
	return f.Length() - f.GetPos()
}

// ReadBytes implements SliceReader - read a slice of n bytes
// Return empty slice if n bytes not available.
func (f *FixedSliceReader) ReadBytes(n int) []byte {
	if n < 0 {
		f.err = ErrAttemptToReadNegativeNumberOfBytes
		return []byte{}
	}
	if f.err != nil {
		return []byte{}
	}
	if f.pos > f.len-n {
		f.err = ErrSliceRead
		return []byte{}
	}
	res := f.slice[f.pos : f.pos+n]
	f.pos += n
	return res
}

// ReadFixedLengthString implements SliceReader - read string of specified length n.
// Sets err and returns empty string if full length not available.
func (f *FixedSliceReader) ReadFixedLengthString(n int) string {
	if f.err != nil {
		return ""
	}
	if f.pos > f.len-n {
		f.err = ErrSliceRead
		return ""
	}
	res := string(f.slice[f.pos : f.pos+n])
	f.pos += n
	return res
}

// ReadInt16 implements SliceReader - read int16 from slice.
func (f *FixedSliceReader) ReadInt16() int16 {
	if f.err != nil {
		return 0
	}
	if f.pos > f.len-2 {
		f.err = ErrSliceRead
		return 0
	}
	res := binary.BigEndian.Uint16(f.slice[f.pos : f.pos+2])
	f.pos += 2
	return int16(res) //nolint:gosec
}

// ReadInt32 implements SliceReader - read int32 from slice.
func (f *FixedSliceReader) ReadInt32() int32 {
	if f.err != nil {
		return 0
	}
	if f.pos > f.len-4 {
		f.err = ErrSliceRead
		return 0
	}
	res := binary.BigEndian.Uint32(f.slice[f.pos : f.pos+4])
	f.pos += 4
	return int32(res) //nolint:gosec
}

// ReadInt64 implements SliceReader - read int64 from slice.
func (f *FixedSliceReader) ReadInt64() int64 {
	if f.err != nil {
		return 0
	}
	if f.pos > f.len-8 {
		f.err = ErrSliceRead
		return 0
	}
	res := binary.BigEndian.Uint64(f.slice[f.pos : f.pos+8])
	f.pos += 8
	return int64(res) //nolint:gosec
}

// ReadPossiblyZeroTerminatedString implements SliceReader - read string until zero byte but at most maxLen.
// If maxLen is reached and no zero-byte, return string and ok = false.
func (f *FixedSliceReader) ReadPossiblyZeroTerminatedString(maxLen int) (string, bool) {
	startPos := f.pos
	maxPos := startPos + maxLen
	for {
		if f.pos == maxPos {
			return string(f.slice[startPos:f.pos]), true
		}
		if f.pos > maxPos {
			f.err = ErrDidNotFindTerminatingZero
			return "", false
		}
		c := f.slice[f.pos]
		if c == 0 {
			str := string(f.slice[startPos:f.pos])
			f.pos++ // Next position to read
			return str, true
		}
		f.pos++
	}
}

// ReadUint16 implements SliceReader - read uint16 from slice.
func (f *FixedSliceReader) ReadUint16() uint16 {
	if f.err != nil {
		return 0
	}
	if f.pos > f.len-2 {
		f.err = ErrSliceRead
		return 0
	}
	res := binary.BigEndian.Uint16(f.slice[f.pos : f.pos+2])
	f.pos += 2
	return res
}

// ReadUint24 implements SliceReader - read uint24 from slice.
func (f *FixedSliceReader) ReadUint24() uint32 {
	if f.err != nil {
		return 0
	}
	if f.pos > f.len-3 {
		f.err = ErrSliceRead
		return 0
	}
	res := uint32(binary.BigEndian.Uint16(f.slice[f.pos : f.pos+2]))
	res = res<<8 | uint32(f.slice[f.pos+2])
	f.pos += 3
	return res
}

// ReadUint32 implements SliceReader - read uint32 from slice.
func (f *FixedSliceReader) ReadUint32() uint32 {
	if f.err != nil {
		return 0
	}
	if f.pos > f.len-4 {
		f.err = ErrSliceRead
		return 0
	}
	res := binary.BigEndian.Uint32(f.slice[f.pos : f.pos+4])
	f.pos += 4
	return res
}

// ReadUint64 implements SliceReader - read uint64 from slice.
func (f *FixedSliceReader) ReadUint64() uint64 {
	if f.err != nil {
		return 0
	}
	if f.pos > f.len-8 {
		f.err = ErrSliceRead
		return 0
	}
	res := binary.BigEndian.Uint64(f.slice[f.pos : f.pos+8])
	f.pos += 8
	return res
}

// ReadUint8 implements SliceReader - read uint8 from slice.
func (f *FixedSliceReader) ReadUint8() byte {
	if f.err != nil {
		return 0
	}
	if f.pos > f.len-1 {
		f.err = ErrSliceRead
		return 0
	}
	res := f.slice[f.pos]
	f.pos++
	return res
}

// ReadZeroTerminatedString implements SliceReader - read string until zero byte but at most maxLen
// Set err and return empty string if no zero byte found.
func (f *FixedSliceReader) ReadZeroTerminatedString(maxLen int) string {
	if f.err != nil {
		return ""
	}
	startPos := f.pos
	maxPos := startPos + maxLen
	if maxPos > f.len {
		maxPos = f.len
	}
	for {
		if f.pos >= maxPos {
			f.err = ErrDidNotFindTerminatingZero
			return ""
		}
		c := f.slice[f.pos]
		if c == 0 {
			str := string(f.slice[startPos:f.pos])
			f.pos++ // Next position to read
			return str
		}
		f.pos++
	}
}

// RemainingBytes implements SliceReader - return remaining bytes of this slice.
func (f *FixedSliceReader) RemainingBytes() []byte {
	if f.err != nil {
		return []byte{}
	}
	res := f.slice[f.pos:]
	f.pos = f.Length()
	return res
}

// SetPos implements SliceReader - set read position is slice.
func (f *FixedSliceReader) SetPos(pos int) {
	if pos > f.len {
		f.err = ErrAttemptToSkipBytesBeyondSliceLen
		return
	}
	f.pos = pos
}

// SkipBytes implements SliceReader - skip passed n bytes.
func (f *FixedSliceReader) SkipBytes(n int) {
	if f.err != nil {
		return
	}
	if f.pos+n > f.Length() {
		f.err = ErrAttemptToSkipBytesBeyondSliceLen
		return
	}
	f.pos += n
}

// NewFixedSliceReader creates a new slice reader reading from data.
func NewFixedSliceReader(data []byte) *FixedSliceReader {
	return &FixedSliceReader{
		slice: data,
		pos:   0,
		len:   len(data),
		err:   nil,
	}
}
