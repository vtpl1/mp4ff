package mp4

import (
	"io"

	"github.com/vtpl1/mp4ff/internal/bits"
)

// SinfBox -  Protection Scheme Information Box according to ISO/IEC 23001-7.
type SinfBox struct {
	Frma     *FrmaBox // Mandatory
	Schm     *SchmBox // Optional
	Schi     *SchiBox // Optional
	Children []Box
}

// GetChildren implements ContainerBox.
func (b *SinfBox) GetChildren() []Box {
	return b.Children
}

// Encode implements Box.
func (b *SinfBox) Encode(w io.Writer) error {
	return EncodeContainer(b, w)
}

// EncodeSW implements Box.
func (b *SinfBox) EncodeSW(sw bits.SliceWriter) error {
	return EncodeContainerSW(b, sw)
}

// Info implements Box.
func (b *SinfBox) Info(w io.Writer, specificBoxLevels string, indent string, indentStep string) error {
	return ContainerInfo(b, w, specificBoxLevels, indent, indentStep)
}

// Size implements Box.
func (b *SinfBox) Size() uint64 {
	return containerSize(b.Children)
}

// Type implements Box.
func (b *SinfBox) Type() string {
	return "sinf"
}

// AddChild - Add a child box.
func (b *SinfBox) AddChild(child Box) {
	switch box := child.(type) {
	case *FrmaBox:
		b.Frma = box
	case *SchmBox:
		b.Schm = box
	case *SchiBox:
		b.Schi = box
	}
	b.Children = append(b.Children, child)
}

// DecodeSinf - box-specific decode.
func DecodeSinf(hdr BoxHeader, startPos uint64, r io.Reader) (Box, error) {
	children, err := DecodeContainerChildren(hdr, startPos+8, startPos+hdr.Size, r)
	if err != nil {
		return nil, err
	}
	b := SinfBox{}
	for _, c := range children {
		b.AddChild(c)
	}
	return &b, nil
}

// DecodeSinfSR - box-specific decode.
func DecodeSinfSR(hdr BoxHeader, startPos uint64, sr bits.SliceReader) (Box, error) {
	children, err := DecodeContainerChildrenSR(hdr, startPos+8, startPos+hdr.Size, sr)
	if err != nil {
		return nil, err
	}
	b := SinfBox{}
	for _, c := range children {
		b.AddChild(c)
	}
	return &b, sr.AccError()
}
