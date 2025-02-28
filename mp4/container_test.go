package mp4_test

import (
	"bytes"
	"testing"

	"github.com/vtpl1/mp4ff/mp4"
)

func TestContainerBox(t *testing.T) {
	var b mp4.Box = mp4.NewGenericContainerBox("test")
	_ = b
}

func TestGenericContainer(t *testing.T) {
	// Just check that it doesn't crash
	c := mp4.NewGenericContainerBox("test")
	// c.AddChild(&VsidBox{SourceID: 42})
	w := bytes.Buffer{}
	err := c.Encode(&w)
	if err != nil {
		t.Error(err)
	}
	err = c.Info(&w, "", "", "  ")
	if err != nil {
		t.Error(err)
	}
}
