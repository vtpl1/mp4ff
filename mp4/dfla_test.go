package mp4_test

import (
	"testing"

	"github.com/Eyevinn/mp4ff/mp4"
)

func TestEncodeDecodeDfla(t *testing.T) {

	b := &mp4.Dfla{}
	if b.Type() != "dfLa" {
		t.Error("error")
	}
}
