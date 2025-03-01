package mp4

import (
	"fmt"
	"io"

	"github.com/vtpl1/mp4ff/internal/bits"
)

// EdtsBox - Edit Box (edts - optional)
//
// Contained in: Track Box ("trak")
//
// The edit box maps the presentation timeline to the media-time line
type EdtsBox struct {
	Elst     []*ElstBox
	Children []Box
}

// DecodeEdts - box-specific decode
func DecodeEdts(hdr BoxHeader, startPos uint64, r io.Reader) (Box, error) {
	l, err := DecodeContainerChildren(hdr, startPos+8, startPos+hdr.Size, r)
	if err != nil {
		return nil, err
	}
	e := &EdtsBox{}
	e.Children = l
	for _, b := range l {
		switch b.Type() {
		case "elst":
			e.Elst = append(e.Elst, b.(*ElstBox))
		default:
			return nil, fmt.Errorf("Box of type %s in edts", b.Type())
		}
	}
	return e, nil
}

// DecodeEdtsSR - box-specific decode
func DecodeEdtsSR(hdr BoxHeader, startPos uint64, sr bits.SliceReader) (Box, error) {
	children, err := DecodeContainerChildrenSR(hdr, startPos+8, startPos+hdr.Size, sr)
	if err != nil {
		return nil, err
	}
	e := &EdtsBox{}
	e.Children = children
	for _, b := range children {
		switch b.Type() {
		case "elst":
			e.Elst = append(e.Elst, b.(*ElstBox))
		default:
			return nil, fmt.Errorf("Box of type %s in edts", b.Type())
		}
	}
	return e, sr.AccError()
}

// AddChild - Add a child box and update EntryCount
func (e *EdtsBox) AddChild(child Box) {
	e.Children = append(e.Children, child)
}

// Type - box type
func (e *EdtsBox) Type() string {
	return "edts"
}

// Size - calculated size of box
func (e *EdtsBox) Size() uint64 {
	return containerSize(e.Children)
}

// GetChildren - list of child boxes
func (e *EdtsBox) GetChildren() []Box {
	return e.Children
}

// Encode - write edts container to w
func (e *EdtsBox) Encode(w io.Writer) error {
	return EncodeContainer(e, w)
}

// EncodeSW - write edts container to sw
func (e *EdtsBox) EncodeSW(sw bits.SliceWriter) error {
	return EncodeContainerSW(e, sw)
}

// Info - write box-specific information
func (e *EdtsBox) Info(w io.Writer, specificBoxLevels, indent, indentStep string) error {
	return ContainerInfo(e, w, specificBoxLevels, indent, indentStep)
}
