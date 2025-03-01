package mp4

import (
	"encoding/hex"
	"fmt"
	"io"

	"github.com/vtpl1/mp4ff/hevc"
	"github.com/vtpl1/mp4ff/internal/bits"
)

// HvcCBox - HEVCConfigurationBox (ISO/IEC 14496-15 8.4.1.1.2)
// Contains one HEVCDecoderConfigurationRecord
type HvcCBox struct {
	hevc.DecConfRec
}

// CreateHvcC - create an hvcC box based on VPS, SPS and PPS and signal completeness
// If includePS is false, the nalus are not included, but information from sps is extracted.
func CreateHvcC(vpsNalus, spsNalus, ppsNalus [][]byte, vpsComplete, spsComplete, ppsComplete, includePS bool) (*HvcCBox, error) {
	hevcDecConfRec, err := hevc.CreateHEVCDecConfRec(vpsNalus, spsNalus, ppsNalus,
		vpsComplete, spsComplete, ppsComplete, includePS)
	if err != nil {
		return nil, fmt.Errorf("CreateHEVCDecConfRec: %w", err)
	}

	return &HvcCBox{hevcDecConfRec}, nil
}

// DecodeHvcC - box-specific decode
func DecodeHvcC(hdr BoxHeader, startPos uint64, r io.Reader) (Box, error) {
	data, err := readBoxBody(r, hdr)
	if err != nil {
		return nil, err
	}
	hevcDecConfRec, err := hevc.DecodeHEVCDecConfRec(data)
	if err != nil {
		return nil, err
	}
	return &HvcCBox{hevcDecConfRec}, nil
}

// DecodeHvcCSR - box-specific decode
func DecodeHvcCSR(hdr BoxHeader, startPos uint64, sr bits.SliceReader) (Box, error) {
	hevcDecConfRec, err := hevc.DecodeHEVCDecConfRec(sr.ReadBytes(hdr.payloadLen()))
	return &HvcCBox{hevcDecConfRec}, err
}

// Type - return box type
func (b *HvcCBox) Type() string {
	return "hvcC"
}

// Size - return calculated size
func (b *HvcCBox) Size() uint64 {
	return uint64(boxHeaderSize + b.DecConfRec.Size())
}

// Encode - write box to w
func (b *HvcCBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	return b.DecConfRec.Encode(w)
}

// Encode - write box to w
func (b *HvcCBox) EncodeSW(sw bits.SliceWriter) error {
	err := EncodeHeaderSW(b, sw)
	if err != nil {
		return err
	}
	return b.DecConfRec.EncodeSW(sw)
}

// Info - box-specific Info
func (b *HvcCBox) Info(w io.Writer, specificBoxLevels, indent, indentStep string) error {
	bd := newInfoDumper(w, indent, b, -1, 0)
	hdcr := b.DecConfRec
	bd.writef(" - GeneralProfileSpace: %d", hdcr.GeneralProfileSpace)
	bd.writef(" - GeneralTierFlag: %t", hdcr.GeneralTierFlag)
	bd.writef(" - GeneralProfileIDC: %d", hdcr.GeneralProfileIDC)
	bd.writef(" - GeneralProfileCompatibilityFlags: %08x", hdcr.GeneralProfileCompatibilityFlags)
	bd.writef(" - GeneralConstraintIndicatorFlags: %012x", hdcr.GeneralConstraintIndicatorFlags)
	bd.writef(" - GeneralLevelIDC: %d", hdcr.GeneralLevelIDC)
	bd.writef(" - MinSpatialSegmentationIDC: %d", hdcr.MinSpatialSegmentationIDC)
	bd.writef(" - ParallellismType: %d", hdcr.ParallellismType)
	bd.writef(" - ChromaFormatIDC: %d", hdcr.ChromaFormatIDC)
	bd.writef(" - BitDepthLuma: %d", hdcr.BitDepthLumaMinus8+8)
	bd.writef(" - BitDepthChroma: %d", hdcr.BitDepthChromaMinus8+8)
	bd.writef(" - AvgFrameRate/256: %d", hdcr.AvgFrameRate)
	bd.writef(" - ConstantFrameRate: %d", hdcr.ConstantFrameRate)
	bd.writef(" - NumTemporalLayers: %d", hdcr.NumTemporalLayers)
	bd.writef(" - temporalIDNested: %d", hdcr.TemporalIDNested)
	for _, array := range hdcr.NaluArrays {
		bd.writef("   - %s complete: %d", array.NaluType(), array.Complete())
		for _, nalu := range array.Nalus {
			bd.writef("    %s", hex.EncodeToString(nalu))
		}
	}
	return bd.err
}
