package sei_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/vtpl1/mp4ff/sei"
)

func TestSei136Clock(t *testing.T) {
	cl := sei.CreateClockTS()
	cl.TimeOffsetValue = 1300
	cl.NFrames = 5
	cl.Hours = 12
	cl.Minutes = 30
	cl.Seconds = 10
	cl.ClockTimeStampFlag = true
	cl.UnitsFieldBasedFlag = true
	cl.FullTimeStampFlag = false
	cl.SecondsFlag = true
	cl.MinutesFlag = true
	cl.HoursFlag = true
	cl.DiscontinuityFlag = false
	cl.CntDroppedFlag = false
	cl.CountingType = 1
	cl.TimeOffsetLength = 11

	tc := sei.TimeCodeSEI{}
	tc.Clocks = append(tc.Clocks, cl)
	tc.Clocks = append(tc.Clocks, cl)
	pl := tc.Payload()
	str := tc.String()
	if str == "" {
		t.Error("String() failed")
	}
	if len(pl) == 0 {
		t.Error("Payload() failed")
	}

	seiData := sei.SEIData{
		PayloadType: sei.SEITimeCodeType,
		PayloadData: pl,
	}
	sei1, err := sei.DecodeTimeCodeSEI(&seiData)
	if err != nil {
		t.Error(err)
	}
	tcDec := sei1.(*sei.TimeCodeSEI)
	decCl := tcDec.Clocks[0]
	diff := deep.Equal(cl, decCl)
	if diff != nil {
		t.Error(diff)
	}
}
