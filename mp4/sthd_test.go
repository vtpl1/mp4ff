package mp4_test

import (
	"bytes"
	"testing"

	"github.com/go-test/deep"
	"github.com/vtpl1/mp4ff/mp4"
)

func TestSthd(t *testing.T) {
	encBox := &mp4.SthdBox{}

	buf := bytes.Buffer{}
	err := encBox.Encode(&buf)
	if err != nil {
		t.Error(err)
	}

	decBox, err := mp4.DecodeBox(0, &buf)
	if err != nil {
		t.Error(err)
	}

	if diff := deep.Equal(encBox, decBox); diff != nil {
		t.Error(diff)
	}
}
