package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Example from Dekay
// https://github.com/dekay/im-me/blob/master/pocketwx/src/protocol.txt
// a0 06 52 83 38 00 5a c8
// ((0x38 >> 4) << 8) + 0x83 = 131 + 768 = 899 = 89.9% Relative Humidity
func TestDecodeHumidity(t *testing.T) {
	message := createMessage([]byte{0xA0, 0x06, 0x52, 0x83, 0x38, 0x00, 0x5a, 0xC8})
	humidity, err := DecodeHumidity(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(89.9), humidity)
}
