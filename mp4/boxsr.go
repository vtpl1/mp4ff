package mp4

import (
	"fmt"

	"github.com/vtpl1/mp4ff/internal/bits"
)

// DecodeHeaderSR - decode a box header (size + box type + possible largeSize) from sr.
func DecodeHeaderSR(sr bits.SliceReader) (BoxHeader, error) {
	size := uint64(sr.ReadUint32())
	boxType := sr.ReadFixedLengthString(4)
	headerLen := boxHeaderSize
	if size == 1 {
		size = sr.ReadUint64()
		headerLen += LargeSizeLen
	} else if size == 0 {
		return BoxHeader{}, ErrEndOfFileNotSupported
	}
	if uint64(headerLen) > size { //nolint:gosec
		return BoxHeader{}, ErrBoxHeaderSizeExceedsBoxSize
	}
	return BoxHeader{boxType, size, headerLen}, sr.AccError()
}

// DecodeBoxSR - decode a box from SliceReader.
func DecodeBoxSR(startPos uint64, sr bits.SliceReader) (Box, error) {
	var b Box

	h, err := DecodeHeaderSR(sr)
	if err != nil {
		return nil, err
	}

	maxSize := uint64(sr.NrRemainingBytes()) + uint64(h.Hdrlen) //nolint:gosec
	// In the following, we do not block mdat to allow for the case
	// that the first kiloBytes of a file are fetched and parsed to
	// get the init part of a file. In the future, a new decode option that
	// stops before the mdat starts is a better alternative.
	if h.Size > maxSize && h.Name != "mdat" {
		return nil, fmt.Errorf("decode box %q, size %d too big (max %d)", h.Name, h.Size, maxSize) //nolint:err113
	}

	d, ok := decodersSR[h.Name]

	if !ok {
		b, err = DecodeUnknownSR(h, startPos, sr)
	} else {
		b, err = d(h, startPos, sr)
	}
	if err != nil {
		return nil, fmt.Errorf("decode %s pos %d: %w", h.Name, startPos, err)
	}

	return b, err
}

//nolint:gochecknoglobals
var decodersSR map[string]BoxDecoderSR

// BoxDecoderSR is function signature of the Box DecodeSR method.
type BoxDecoderSR func(hdr BoxHeader, startPos uint64, sw bits.SliceReader) (Box, error)

//nolint:funlen,gochecknoinits
func init() {
	decodersSR = map[string]BoxDecoderSR{
		"\xa9ART": DecodeGenericContainerBoxSR,
		"\xa9cpy": DecodeGenericContainerBoxSR,
		"\xa9nam": DecodeGenericContainerBoxSR,
		"\xa9too": DecodeGenericContainerBoxSR,
		"ac-3":    DecodeAudioSampleEntrySR,
		// "alou":    DecodeAlouBoxSR,
		// "av01":    DecodeVisualSampleEntrySR,
		// "av1C":    DecodeAv1CSR,
		// "avc1":    DecodeVisualSampleEntrySR,
		// "avc3":    DecodeVisualSampleEntrySR,
		// "avcC":    DecodeAvcCSR,
		"btrt": DecodeBtrtSR,
		// "cdat":    DecodeCdatSR,
		// "cdsc":    DecodeTrefTypeSR,
		// "clap":    DecodeClapSR,
		// "co64":    DecodeCo64SR,
		// "CoLL":    DecodeCoLLSR,
		// "colr":    DecodeColrSR,
		// "cslg":    DecodeCslgSR,
		"ctim": DecodeCtimSR,
		// "ctts":    DecodeCttsSR,
		"dac3": DecodeDac3SR,
		// "data":    DecodeDataSR,
		"dec3": DecodeDec3SR,
		"desc": DecodeGenericContainerBoxSR,
		// "dinf":    DecodeDinfSR,
		// "dpnd":    DecodeTrefTypeSR,
		// "dref":    DecodeDrefSR,
		"ec-3": DecodeAudioSampleEntrySR,
		// "edts":    DecodeEdtsSR,
		// "elng":    DecodeElngSR,
		// "elst":    DecodeElstSR,
		// "emeb":    DecodeEmebSR,
		// "emib":    DecodeEmibSR,
		// "emsg":    DecodeEmsgSR,
		"enca": DecodeAudioSampleEntrySR,
		// "encv":    DecodeVisualSampleEntrySR,
		"esds": DecodeEsdsSR,
		// "evte":    DecodeEvteSR,
		// "font":    DecodeTrefTypeSR,
		// "free":    DecodeFreeSR,
		"frma": DecodeFrmaSR,
		// "ftyp":    DecodeFtypSR,
		// "hdlr":    DecodeHdlrSR,
		// "hev1":    DecodeVisualSampleEntrySR,
		// "hind":    DecodeTrefTypeSR,
		// "hint":    DecodeTrefTypeSR,
		// "hvc1":    DecodeVisualSampleEntrySR,
		// "hvcC":    DecodeHvcCSR,
		"iden": DecodeIdenSR,
		// "ilst":    DecodeIlstSR,
		// "iods":    DecodeUnknownSR,
		// "ipir":    DecodeTrefTypeSR,
		// "kind":    DecodeKindSR,
		// "leva":    DecodeLevaSR,
		// "ludt":    DecodeLudtSR,
		"mdat": DecodeMdatSR,
		// "mehd":    DecodeMehdSR,
		// "mdhd":    DecodeMdhdSR,
		// "mdia":    DecodeMdiaSR,
		// "meta":    DecodeMetaSR,
		// "mfhd":    DecodeMfhdSR,
		// "mfra":    DecodeMfraSR,
		// "mfro":    DecodeMfroSR,
		// "mime":    DecodeMimeSR,
		// "minf":    DecodeMinfSR,
		// "moof":    DecodeMoofSR,
		// "moov":    DecodeMoovSR,
		"mp4a": DecodeAudioSampleEntrySR,
		// "mpod":    DecodeTrefTypeSR,
		// "mvex":    DecodeMvexSR,
		// "mvhd":    DecodeMvhdSR,
		// "nmhd":    DecodeNmhdSR,
		// "pasp":    DecodePaspSR,
		"payl": DecodePaylSR,
		// "prft":    DecodePrftSR,
		// "pssh":    DecodePsshSR,
		// "saio":    DecodeSaioSR,
		"saiz": DecodeSaizSR,
		// "sbgp":    DecodeSbgpSR,
		"schi": DecodeSchiSR,
		"schm": DecodeSchmSR,
		// "sdtp":    DecodeSdtpSR,
		"senc": DecodeSencSR,
		// "sgpd":    DecodeSgpdSR,
		"sidx": DecodeSidxSR,
		// "silb":    DecodeSilbSR,
		"sinf": DecodeSinfSR,
		// "skip":    DecodeFreeSR,
		// "SmDm":    DecodeSmDmSR,
		// "smhd":    DecodeSmhdSR,
		// "ssix":    DecodeSsixSR,
		// "stbl":    DecodeStblSR,
		// "stco":    DecodeStcoSR,
		// "sthd":    DecodeSthdSR,
		// "stpp":    DecodeStppSR,
		// "stsc":    DecodeStscSR,
		// "stsd":    DecodeStsdSR,
		// "stss":    DecodeStssSR,
		// "stsz":    DecodeStszSR,
		"sttg": DecodeSttgSR,
		// "stts":    DecodeSttsSR,
		"styp": DecodeStypSR,
		// "subs":    DecodeSubsSR,
		// "subt":    DecodeTrefTypeSR,
		// "sync":    DecodeTrefTypeSR,
		"tenc": DecodeTencSR,
		// "tfdt":    DecodeTfdtSR,
		// "tfhd":    DecodeTfhdSR,
		// "tfra":    DecodeTfraSR,
		// "tkhd":    DecodeTkhdSR,
		// "tlou":    DecodeTlouSR,
		// "traf":    DecodeTrafSR,
		// "trak":    DecodeTrakSR,
		// "tref":    DecodeTrefSR,
		// "trep":    DecodeTrepSR,
		// "trex":    DecodeTrexSR,
		// "trun":    DecodeTrunSR,
		// "udta":    DecodeUdtaSR,
		// "url ":    DecodeURLBoxSR, //nolint:gocritic
		"uuid": DecodeUUIDBoxSR,
		// "vdep":    DecodeTrefTypeSR,
		"vlab": DecodeVlabSR,
		// "vmhd":    DecodeVmhdSR,
		// "vp08":    DecodeVisualSampleEntrySR,
		// "vp09":    DecodeVisualSampleEntrySR,
		// "vpcC":    DecodeVppCSR,
		// "vplx":    DecodeTrefTypeSR,
		"vsid": DecodeVsidSR,
		"vtta": DecodeVttaSR,
		"vttc": DecodeVttcSR,
		"vttC": DecodeVttCSR,
		"vtte": DecodeVtteSR,
		"wvtt": DecodeWvttSR,
	}
}
