package processor

import (
	"testing"

	"github.com/nathanmsmith/rtldavis/dsp"
	"github.com/nathanmsmith/rtldavis/protocol"
	"github.com/stretchr/testify/assert"
)

func createMessage(data []byte) protocol.Message {
	return protocol.Message{
		Packet: dsp.Packet{
			Idx:  0,
			Data: data,
		},
		ID: 0,
	}
}

func TestDecodeTemperatureInvalidMessage(t *testing.T) {
	message := createMessage([]byte{0xA0, 0x00, 0x00, 0xC9, 0x3D, 0x00, 0x2A, 0x87})

	temp, err := DecodeTemperature(message)
	assert.ErrorContains(t, err, "Message does not have temperature")
	assert.Equal(t, float32(-1.0), temp)
}

func TestDecodeTemperatureAnalogSensor(t *testing.T) {
	message := createMessage([]byte{0x81, 0x00, 0x00, 0x59, 0x45, 0x00, 0xA3, 0xE6})

	temp, err := DecodeTemperature(message)
	assert.ErrorContains(t, err, "Temperature reading is not from digital sensor. Analog sensor not supported")
	assert.Equal(t, float32(-1.0), temp)
}

func TestDecodeTemperatureNoSensor(t *testing.T) {
	message := createMessage([]byte{0x81, 0x00, 0xDB, 0xFF, 0xC8, 0x00, 0xAB, 0xF8})

	temp, err := DecodeTemperature(message)
	assert.ErrorContains(t, err, "No sensor")
	assert.Equal(t, float32(-1.0), temp)
}

func TestDecodeTemperature40Degrees(t *testing.T) {
	message := createMessage([]byte{0x80, 0x01, 0xa2, 0x19, 0x89, 0x04, 0x45, 0x19})

	temp, err := DecodeTemperature(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(40.8), temp)
}

func TestDecodeTemperature82Degrees(t *testing.T) {
	message := createMessage([]byte{0x80, 0x00, 0x00, 0x33, 0x8D, 0x00, 0x25, 0x11})

	temp, err := DecodeTemperature(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(82.4), temp)
}
