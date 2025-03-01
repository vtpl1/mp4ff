package mp4

import "strings"

// EncFragFileMode - mode for writing file
type EncFragFileMode byte

const (
	// EncModeSegment - only encode boxes that are part of Init and MediaSegments
	EncModeSegment = EncFragFileMode(0)
	// EncModeBoxTree - encode all boxes in file tree
	EncModeBoxTree = EncFragFileMode(1)
)

// DecFileMode - mode for decoding file
type DecFileMode byte

const (
	// DecModeNormal - read Mdat data into memory during file decoding.
	DecModeNormal DecFileMode = iota
	// DecModeLazyMdat - do not read mdat data into memory.
	// Thus, decode process requires less memory and faster.
	DecModeLazyMdat
)

// DecFileFlags can be combined for special decoding options
type DecFileFlags uint32

const (
	DecNoFlags DecFileFlags = 0
	// DecISMFlag tries to read mfra box at end to find segment boundaries (for ISM files)
	DecISMFlag DecFileFlags = (1 << 0)
	// DecStartOnMoof starts a segment at each moof boundary
	// This is provided no styp, or sidx/mfra box gives other information
	DecStartOnMoof = (1 << 1)
	// if no styp box, or sidx/mfra strudture
)

// EncOptimize - encoder optimization mode
type EncOptimize uint32

const (
	// OptimizeNone - no optimization
	OptimizeNone = EncOptimize(0)
	// OptimizeTrun - optimize trun box by moving default values to tfhd
	OptimizeTrun = EncOptimize(1 << 0)
)

func (eo EncOptimize) String() string {
	var optList []string
	msg := "OptimizeNone"
	if eo&OptimizeTrun != 0 {
		optList = append(optList, "OptimizeTrun")
	}
	if len(optList) > 0 {
		msg = strings.Join(optList, " | ")
	}
	return msg
}
