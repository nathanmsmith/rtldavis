package protocol

import (
	"testing"

	"github.com/nathanmsmith/rtldavis/dsp"
	"github.com/stretchr/testify/assert"
)

func createMessage(data []byte) Message {
	return Message{
		Packet: dsp.Packet{
			Idx:  0,
			Data: data,
		},
		ID: 0,
	}
}

func TestDecodeTemperatureInvalidMessage(t *testing.T) {
	message := createMessage([]byte{0x80, 0x14, 0x35, 0x19, 0x28, 0x00})

	result, _ := DecodeTemperature(testData.message)
	assert.Equal(t, a, b, "The two words should be the same.")

}
