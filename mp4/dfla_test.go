package mp4_test

import (
	"testing"

	"github.com/Eyevinn/mp4ff/mp4"
)

func TestEncodeDecodeDfla(t *testing.T) {
	t.Run("dfLa", func(t *testing.T) {
		dfla := mp4.CreateDfla(8000)
		boxDiffAfterEncodeAndDecode(t, dfla)
	})
}
