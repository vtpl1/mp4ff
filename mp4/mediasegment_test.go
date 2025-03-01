package mp4_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/vtpl1/mp4ff/mp4"
)

func TestMediaSegmentFragmentation(t *testing.T) {
	trex := &mp4.TrexBox{
		TrackID: 2,
	}

	inFile := "testdata/1.m4s"
	inFileGoldenDumpPath := "testdata/golden_1_m4s_dump.txt"
	goldenFragPath := "testdata/golden_1_frag.m4s"
	goldenFragDumpPath := "testdata/golden_1_frag_m4s_dump.txt"
	fd, err := os.Open(inFile)
	if err != nil {
		t.Fatal(err)
	}
	defer fd.Close()

	f, err := mp4.DecodeFile(fd)
	if err != io.EOF && err != nil {
		t.Error(err)
	}
	if len(f.Segments) != 1 {
		t.Errorf("Not exactly one mediasegment")
	}

	var bufInSeg bytes.Buffer
	f.EncOptimize = mp4.OptimizeNone // Avoid trun optimization
	f.FragEncMode = mp4.EncModeBoxTree
	err = f.Encode(&bufInSeg)
	if err != nil {
		t.Error(err)
	}

	inSeg, err := os.ReadFile(inFile)
	if err != nil {
		t.Fatal(err)
	}

	diff := deep.Equal(inSeg, bufInSeg.Bytes())
	if diff != nil {
		t.Errorf("Written segment differs from %s", inFile)
	}

	err = compareOrUpdateInfo(t, f, inFileGoldenDumpPath)
	if err != nil {
		t.Error(err)
	}

	mediaSegment := f.Segments[0]
	var timeScale uint64 = 90000
	var duration uint32 = 45000

	fragments, err := mediaSegment.Fragmentify(timeScale, trex, duration)
	if err != nil {
		t.Errorf("Fragmentation went wrong")
	}
	if len(fragments) != 4 {
		t.Errorf("%d fragments instead of 4", len(fragments))
	}

	var bufFrag bytes.Buffer
	fragmentedSegment := mp4.NewMediaSegment()
	fragmentedSegment.EncOptimize = mp4.OptimizeTrun
	fragmentedSegment.Styp = f.Segments[0].Styp
	fragmentedSegment.Fragments = fragments

	err = fragmentedSegment.Encode(&bufFrag)
	if err != nil {
		t.Error(err)
	}

	err = compareOrUpdateInfo(t, fragmentedSegment, goldenFragDumpPath)
	if err != nil {
		t.Error(err)
	}

	if *update {
		err = writeGolden(t, goldenFragPath, bufFrag.Bytes())
		if err != nil {
			t.Error(err)
		}
	} else {
		goldenFrag, err := os.ReadFile(goldenFragPath)
		if err != nil {
			t.Error(err)
		}
		diff := deep.Equal(goldenFrag, bufFrag.Bytes())
		if diff != nil {
			t.Errorf("Generated dump different from %s", goldenFragPath)
		}
	}
}

func TestDoubleDecodeEncodeOptimize(t *testing.T) {
	inFile := "testdata/1.m4s"

	fd, err := os.Open(inFile)
	if err != nil {
		t.Fatal(err)
	}
	defer fd.Close()

	enc1 := decodeEncode(t, fd, mp4.OptimizeTrun)
	buf1 := bytes.NewBuffer(enc1)
	enc2 := decodeEncode(t, buf1, mp4.OptimizeTrun)
	diff := deep.Equal(enc2, enc1)
	if diff != nil {
		t.Errorf("Second write gives diff %s", diff)
	}
}

func TestDecodeEncodeNoOptimize(t *testing.T) {
	inFile := "testdata/1.m4s"

	data, err := os.ReadFile(inFile)
	if err != nil {
		t.Fatal(err)
	}
	buf0 := bytes.NewBuffer(data)
	enc := decodeEncode(t, buf0, mp4.OptimizeNone)
	diff := deep.Equal(enc, data)
	if diff != nil {
		t.Errorf("First encode gives diff %s", diff)
	}
}

func decodeEncode(t *testing.T, r io.Reader, optimize mp4.EncOptimize) []byte {
	f, err := mp4.DecodeFile(r)
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.Buffer{}
	f.EncOptimize = optimize
	err = f.Encode(&buf)
	if err != nil {
		t.Error(err)
	}
	return buf.Bytes()
}

func TestMoofEncrypted(t *testing.T) {
	inFile := "testdata/moof_enc.m4s"
	inFileGoldenDumpPath := "testdata/golden_moof_enc_m4s_dump.txt"
	fd, err := os.Open(inFile)
	if err != nil {
		t.Fatal(err)
	}
	defer fd.Close()

	f, err := mp4.DecodeFile(fd)
	if err != io.EOF && err != nil {
		t.Error(err)
	}

	var bufOut bytes.Buffer
	f.FragEncMode = mp4.EncModeBoxTree
	err = f.Encode(&bufOut)
	if err != nil {
		t.Error(err)
	}

	inSeg, err := os.ReadFile(inFile)
	if err != nil {
		t.Fatal(err)
	}

	diff := deep.Equal(inSeg, bufOut.Bytes())
	if diff != nil {
		tmpOutput := "testdata/moof_enc_tmp.mp4"
		err := writeGolden(t, tmpOutput, bufOut.Bytes())
		if err == nil {
			t.Errorf("Encoded output not same as input for %s. Wrote %s", inFile, tmpOutput)
		} else {
			t.Errorf("Encoded output not same as input for %s, but error %s when writing  %s", inFile, err, tmpOutput)
		}
	}

	err = compareOrUpdateInfo(t, f, inFileGoldenDumpPath)
	if err != nil {
		t.Error(err)
	}
}

func TestDecodeEncodeCencFragmentedFile(t *testing.T) {
	inData, err := os.ReadFile("testdata/prog_8s_enc_dashinit.mp4")
	if err != nil {
		t.Fatal(err)
	}
	inBuf := bytes.NewBuffer(inData)
	decFile, err := mp4.DecodeFile(inBuf)
	if err != nil {
		t.Error(err)
	}
	outSlice := make([]byte, 0, len(inData))
	outBuf := bytes.NewBuffer(outSlice)
	decFile.FragEncMode = mp4.EncFragFileMode(mp4.EncModeBoxTree)
	err = decFile.Encode(outBuf)
	if err != nil {
		t.Error(err)
	}
	outData := outBuf.Bytes()
	if !bytes.Equal(inData, outData) {
		t.Errorf("generated bytes differ from input")
	}
}

func TestCommonSampleDuration(t *testing.T) {
	cases := []struct {
		inFile          string
		trackID         uint32
		wantedCommonDur uint32
		wantedError     string
	}{
		{"testdata/1.m4s", 2, 3000, ""},
		{"testdata/golden_1_frag.m4s", 2, 3000, ""},
		{"testdata/1.m4s", 1, 0, "fragment.CommonSampleDuration: no track with trex trackID=1"},
	}
	for _, c := range cases {
		fd, err := os.Open(c.inFile)
		if err != nil {
			t.Fatal(err)
		}
		defer fd.Close()
		f, err := mp4.DecodeFile(fd)
		if err != nil {
			t.Fatal(err)
		}
		trex := &mp4.TrexBox{
			TrackID: c.trackID,
		}
		for _, s := range f.Segments {
			gotCommonDur, err := s.CommonSampleDuration(trex)
			if c.wantedError != "" {
				if err == nil {
					t.Errorf("case %s: wanted error %s, got nil", c.inFile, c.wantedError)
				}
				if err.Error() != c.wantedError {
					t.Errorf("case %q: wanted error %q, got %q", c.inFile, c.wantedError, err.Error())
				}
				continue
			}
			if err != nil {
				t.Fatal(err)
			}
			if gotCommonDur != c.wantedCommonDur {
				t.Errorf("case %s: got commonDur %d, wanted %d", c.inFile, gotCommonDur, c.wantedCommonDur)
			}
		}
	}
}
