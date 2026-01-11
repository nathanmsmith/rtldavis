package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeHumidity(t *testing.T) {
	message := createMessage([]byte{0xA0, 0x06, 0x52, 0x83, 0x38, 0x00, 0x5a, 0xC8})
	humidity, err := DecodeHumidity(message)
	assert.NoError(t, err)
	assert.Equal(t, 89.9, humidity)
}
