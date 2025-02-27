package mp4

import (
	"io"

	"github.com/vtpl1/mp4ff/internal/bits"
)

// BoxHeader - 8 or 16 bytes depending on size.
type BoxHeader struct {
	Name   string
	Size   uint64
	Hdrlen int
}

func (b BoxHeader) payloadLen() int {
	return int(b.Size) - b.Hdrlen //nolint:gosec
}

// Box is the general interface to any ISOBMFF box or similar.
type Box interface {
	// Type of box, normally 4 asccii characters, but is uint32 according to spec
	Type() string
	// Size of box including header and all children if any
	Size() uint64
	// Encode box to writer
	Encode(w io.Writer) error
	// Encode box to SliceWriter
	EncodeSW(sw bits.SliceWriter) error
	// Info - write box details
	//   spedificBoxLevels is a comma-separated list box:level or all:level where level >= 0.
	//   Higher levels give more details. 0 is default
	//   indent is indent at this box level.
	//   indentStep is how much to indent at each level
	Info(w io.Writer, specificBoxLevels, indent, indentStep string) error
}
