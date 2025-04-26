package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeSupercap(t *testing.T) {
	message := createMessage([]byte{0x20, 0x04, 0xC3, 0xD4, 0xC1, 0x81, 0x89, 0xEE})
	voltage, err := DecodeSupercap(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(0.83), voltage)
}
