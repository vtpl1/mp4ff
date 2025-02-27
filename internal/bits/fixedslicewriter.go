package bits

import "encoding/binary"

// FixedSliceWriter - write numbers to a fixed []byte slice.
type FixedSliceWriter struct {
	accError error
	buf      []byte
	off      int
	n        int  // current number of bits
	v        uint // current accumulated value for bits
}

// AccError implements SliceWriter.
func (f *FixedSliceWriter) AccError() error {
	return f.accError
}

// Bytes implements SliceWriter.
func (f *FixedSliceWriter) Bytes() []byte {
	return f.buf[:f.off]
}

// Capacity implements SliceWriter - max length of FixedSliceWriter buffer.
func (f *FixedSliceWriter) Capacity() int {
	return len(f.buf)
}

// FlushBits implements SliceWriter - write remaining bits to the underlying .Writer.
// bits will be left-shifted and zeros appended to fill up a byte.
func (f *FixedSliceWriter) FlushBits() {
	if f.accError != nil {
		return
	}
	if f.n != 0 {
		b := byte((f.v << (8 - uint(f.n))) & Mask(8)) //nolint:gosec
		f.WriteUint8(b)
	}
}

// Len implements SliceWriter - length of FixedSliceWriter buffer written. Same as Offset().
func (f *FixedSliceWriter) Len() int {
	return f.off
}

// Offset implements SliceWriter - offset for writing in FixedSliceWriter buffer.
func (f *FixedSliceWriter) Offset() int {
	return f.off
}

// WriteBits implements SliceWriter.
func (f *FixedSliceWriter) WriteBits(bits uint, n int) {
	if f.accError != nil {
		return
	}
	f.v <<= uint(n) //nolint:gosec
	f.v |= bits & Mask(n)
	f.n += n
	for f.n >= 8 {
		b := byte((f.v >> (uint(f.n) - 8)) & Mask(8))
		f.WriteUint8(b)
		f.n -= 8
	}
	f.v &= Mask(8)
}

// WriteBytes implements SliceWriter - write []byte.
func (f *FixedSliceWriter) WriteBytes(byteSlice []byte) {
	if f.off+len(byteSlice) > len(f.buf) {
		f.accError = ErrSliceWrite
		return
	}
	copy(f.buf[f.off:f.off+len(byteSlice)], byteSlice)
	f.off += len(byteSlice)
}

// WriteFlag implements SliceWriter - writes a flag as 1 bit.
func (f *FixedSliceWriter) WriteFlag(b bool) {
	bit := uint(0)
	if b {
		bit = 1
	}
	f.WriteBits(bit, 1)
}

// WriteInt16 implements SliceWriter - write int16 to slice.
func (f *FixedSliceWriter) WriteInt16(n int16) {
	if f.off+2 > len(f.buf) {
		f.accError = ErrSliceWrite
		return
	}
	binary.BigEndian.PutUint16(f.buf[f.off:], uint16(n)) //nolint:gosec
	f.off += 2
}

// WriteInt32 implements SliceWriter - write int32 to slice.
func (f *FixedSliceWriter) WriteInt32(n int32) {
	if f.off+4 > len(f.buf) {
		f.accError = ErrSliceWrite
		return
	}
	binary.BigEndian.PutUint32(f.buf[f.off:], uint32(n)) //nolint:gosec
	f.off += 4
}

// WriteInt64 implements SliceWriter - write int64 to slice.
func (f *FixedSliceWriter) WriteInt64(n int64) {
	if f.off+8 > len(f.buf) {
		f.accError = ErrSliceWrite
		return
	}
	binary.BigEndian.PutUint64(f.buf[f.off:], uint64(n)) //nolint:gosec
	f.off += 8
}

// WriteString implements SliceWriter - write string to slice with or without zero end.
func (f *FixedSliceWriter) WriteString(s string, addZeroEnd bool) {
	nrNew := len(s)
	if addZeroEnd {
		nrNew++
	}
	if f.off+nrNew > len(f.buf) {
		f.accError = ErrSliceWrite
		return
	}
	copy(f.buf[f.off:f.off+len(s)], s)
	f.off += len(s)
	if addZeroEnd {
		f.buf[f.off] = 0
		f.off++
	}
}

// WriteUint16 implements SliceWriter - write uint16 to slice.
func (f *FixedSliceWriter) WriteUint16(u uint16) {
	if f.off+2 > len(f.buf) {
		f.accError = ErrSliceWrite
		return
	}
	binary.BigEndian.PutUint16(f.buf[f.off:], u)
	f.off += 2
}

// WriteUint24 implements SliceWriter - write uint24 to slice.
func (f *FixedSliceWriter) WriteUint24(u uint32) {
	if f.off+3 > len(f.buf) {
		f.accError = ErrSliceWrite
		return
	}
	f.WriteUint8(byte(u >> 16))
	f.WriteUint16(uint16(u & 0xffff)) //nolint:gosec
}

// WriteUint32 implements SliceWriter - write uint32 to slice.
func (f *FixedSliceWriter) WriteUint32(u uint32) {
	if f.off+4 > len(f.buf) {
		f.accError = ErrSliceWrite
		return
	}
	binary.BigEndian.PutUint32(f.buf[f.off:], u)
	f.off += 4
}

// WriteUint48 implements SliceWriter - write uint48.
func (f *FixedSliceWriter) WriteUint48(u uint64) {
	if f.off+6 > len(f.buf) {
		f.accError = ErrSliceWrite
		return
	}
	msb := uint16(u >> 32) //nolint:gosec
	binary.BigEndian.PutUint16(f.buf[f.off:], msb)
	f.off += 2

	lsb := uint32(u & 0xffffffff) //nolint:gosec
	binary.BigEndian.PutUint32(f.buf[f.off:], lsb)
	f.off += 4
}

// WriteUint64 implements SliceWriter - write uint64 to slice.
func (f *FixedSliceWriter) WriteUint64(u uint64) {
	if f.off+8 > len(f.buf) {
		f.accError = ErrSliceWrite
		return
	}
	binary.BigEndian.PutUint64(f.buf[f.off:], u)
	f.off += 8
}

// WriteUint8 implements SliceWriter - write byte to slice.
func (f *FixedSliceWriter) WriteUint8(n byte) {
	if f.off+1 > len(f.buf) {
		f.accError = ErrSliceWrite
		return
	}
	f.buf[f.off] = n
	f.off++
}

// WriteUnityMatrix implements SliceWriter WriteUnityMatrix - write a unity matrix for mvhd or tkhd.
func (f *FixedSliceWriter) WriteUnityMatrix() {
	if f.off+36 > len(f.buf) {
		f.accError = ErrSliceWrite
		return
	}
	f.WriteUint32(0x00010000) // = 1 fixed 16.16
	f.WriteUint32(0)
	f.WriteUint32(0)
	f.WriteUint32(0)
	f.WriteUint32(0x00010000) // = 1 fixed 16.16
	f.WriteUint32(0)
	f.WriteUint32(0)
	f.WriteUint32(0)
	f.WriteUint32(0x40000000) // = 1 fixed 2.30
}

// WriteZeroBytes implements SliceWriter - write n byte of zeroes.
func (f *FixedSliceWriter) WriteZeroBytes(n int) {
	if f.off+n > len(f.buf) {
		f.accError = ErrSliceWrite
		return
	}
	for range n {
		f.buf[f.off] = 0
		f.off++
	}
}

// NewFixedSliceWriterFromSlice - create writer around slice.
// The slice will not grow, but stay the same size.
// If too much data is written, there will be
// an accumuluated error. Can be retrieved with AccError().
func NewFixedSliceWriterFromSlice(data []byte) *FixedSliceWriter {
	return &FixedSliceWriter{
		buf:      data,
		off:      0,
		n:        0,
		v:        0,
		accError: nil,
	}
}

// NewFixedSliceWriter - create slice writer with fixed size.
func NewFixedSliceWriter(size int) *FixedSliceWriter {
	return &FixedSliceWriter{
		buf:      make([]byte, size),
		off:      0,
		n:        0,
		v:        0,
		accError: nil,
	}
}
