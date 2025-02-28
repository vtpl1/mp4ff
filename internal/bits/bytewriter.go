package bits

import (
	"encoding/binary"
	"io"
)

// ByteWriter - writer that wraps an io.Writer and accumulates error.
// Only the first error is saved, but any later calls will not panic.
type ByteWriter struct {
	w   io.Writer
	err error
}

// AccError implements SliceWriter - return accumulated error.
func (b *ByteWriter) AccError() error {
	return b.err
}

// Bytes implements SliceWriter.
func (b *ByteWriter) Bytes() []byte {
	panic("unimplemented")
}

// Capacity implements SliceWriter.
func (b *ByteWriter) Capacity() int {
	panic("unimplemented")
}

// FlushBits implements SliceWriter.
func (b *ByteWriter) FlushBits() {
	panic("unimplemented")
}

// Len implements SliceWriter.
func (b *ByteWriter) Len() int {
	panic("unimplemented")
}

// Offset implements SliceWriter.
func (b *ByteWriter) Offset() int {
	panic("unimplemented")
}

// WriteBits implements SliceWriter.
func (b *ByteWriter) WriteBits(_ uint, _ int) {
	panic("unimplemented")
}

// WriteBytes implements SliceWriter.
func (b *ByteWriter) WriteBytes(_ []byte) {
	panic("unimplemented")
}

// WriteFlag implements SliceWriter.
func (b *ByteWriter) WriteFlag(_ bool) {
	panic("unimplemented")
}

// WriteInt16 implements SliceWriter.
func (b *ByteWriter) WriteInt16(_ int16) {
	panic("unimplemented")
}

// WriteInt32 implements SliceWriter.
func (b *ByteWriter) WriteInt32(_ int32) {
	panic("unimplemented")
}

// WriteInt64 implements SliceWriter.
func (b *ByteWriter) WriteInt64(_ int64) {
	panic("unimplemented")
}

// WriteString implements SliceWriter.
func (b *ByteWriter) WriteString(_ string, _ bool) {
	panic("unimplemented")
}

// WriteUint16 implements SliceWriter - write uint16.
func (b *ByteWriter) WriteUint16(u uint16) {
	if b.err != nil {
		return
	}
	b.err = binary.Write(b.w, binary.BigEndian, u)
}

// WriteUint24 implements SliceWriter.
func (b *ByteWriter) WriteUint24(_ uint32) {
	panic("unimplemented")
}

// WriteUint32 implements SliceWriter - write uint32.
func (b *ByteWriter) WriteUint32(u uint32) {
	if b.err != nil {
		return
	}
	b.err = binary.Write(b.w, binary.BigEndian, u)
}

// WriteUint48 implements SliceWriter - write uint48.
func (b *ByteWriter) WriteUint48(u uint64) {
	if b.err != nil {
		return
	}
	msb := uint16(u >> 32) //nolint:gosec
	b.err = binary.Write(b.w, binary.BigEndian, msb)
	if b.err != nil {
		return
	}
	lsb := uint32(u & 0xffffffff) //nolint:gosec
	b.err = binary.Write(b.w, binary.BigEndian, lsb)
}

// WriteUint64 implements SliceWriter - write uint64.
func (b *ByteWriter) WriteUint64(u uint64) {
	if b.err != nil {
		return
	}
	b.err = binary.Write(b.w, binary.BigEndian, u)
}

// WriteUint8 implements SliceWriter - write a byte.
func (b *ByteWriter) WriteUint8(n byte) {
	if b.err != nil {
		return
	}
	b.err = binary.Write(b.w, binary.BigEndian, n)
}

// WriteUnityMatrix implements SliceWriter.
func (b *ByteWriter) WriteUnityMatrix() {
	panic("unimplemented")
}

// WriteZeroBytes implements SliceWriter.
func (b *ByteWriter) WriteZeroBytes(_ int) {
	panic("unimplemented")
}

// WriteSlice - write a slice.
func (b *ByteWriter) WriteSlice(s []byte) {
	if b.err != nil {
		return
	}
	_, b.err = b.w.Write(s)
}

// NewByteWriter creates accumulated error writer around io.Writer.
func NewByteWriter(w io.Writer) *ByteWriter {
	return &ByteWriter{
		w:   w,
		err: nil,
	}
}
