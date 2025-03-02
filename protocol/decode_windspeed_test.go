package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeWindspeed(t *testing.T) {
	message := createMessage([]byte{0xA0, 0x00, 0x00, 0xC9, 0x3D, 0x00, 0x2A, 0x87})

	temp, err := DecodeWindspeed(message)
	assert.ErrorContains(t, err, "Message does not have temperature")
	assert.Equal(t, float32(-1.0), temp)
}
