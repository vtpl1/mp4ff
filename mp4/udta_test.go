package mp4

import "testing"

func TestUdta(t *testing.T) {
	udta := &UdtaBox{}
	unknown := &UnknownBox{
		Name:       "\xa9enc",
		SizeN:      12,
		NotDecoded: []byte{0, 0, 0, 0},
	}

	udta.AddChild(unknown) // Any arbitrary box
	boxDiffAfterEncodeAndDecode(t, udta)
}
