package mp4

import (
	"io"

	"github.com/vtpl1/mp4ff/internal/bits"
)

// SchiBox -  Schema Information Box.
//
//nolint:recvcheck
type SchiBox struct {
	Tenc     *TencBox
	Children []Box
}

// GetChildren implements ContainerBox.
func (s SchiBox) GetChildren() []Box {
	return s.Children
}

// Encode implements Box.
func (s SchiBox) Encode(w io.Writer) error {
	return EncodeContainer(s, w)
}

// EncodeSW implements Box.
func (s SchiBox) EncodeSW(sw bits.SliceWriter) error {
	return EncodeContainerSW(s, sw)
}

// Info implements Box.
func (s SchiBox) Info(w io.Writer, specificBoxLevels string, indent string, indentStep string) error {
	return ContainerInfo(s, w, specificBoxLevels, indent, indentStep)
}

// Size implements Box.
func (s SchiBox) Size() uint64 {
	return containerSize(s.Children)
}

// Type implements Box.
func (s SchiBox) Type() string {
	return "schi"
}

// AddChild - Add a child box.
func (s *SchiBox) AddChild(child Box) {
	switch box := child.(type) { //nolint:gocritic
	case *TencBox:
		s.Tenc = box
	}
	s.Children = append(s.Children, child)
}

// DecodeSchi - box-specific decode.
func DecodeSchi(hdr BoxHeader, startPos uint64, r io.Reader) (Box, error) {
	children, err := DecodeContainerChildren(hdr, startPos+8, startPos+hdr.Size, r)
	if err != nil {
		return nil, err
	}
	b := &SchiBox{}
	for _, child := range children {
		b.AddChild(child)
	}
	return b, nil
}

// DecodeSchiSR - box-specific decode.
func DecodeSchiSR(hdr BoxHeader, startPos uint64, sr bits.SliceReader) (Box, error) {
	children, err := DecodeContainerChildrenSR(hdr, startPos+8, startPos+hdr.Size, sr)
	if err != nil {
		return nil, err
	}
	b := &SchiBox{}
	for _, child := range children {
		b.AddChild(child)
	}
	return b, nil
}
