package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeSolarRadiationInvalidMessage(t *testing.T) {
	message := createMessage([]byte{0x80, 0x00, 0x00, 0x45, 0x00, 0xAD, 0x21})

	radiation, err := DecodeSolarRadiation(message)
	assert.ErrorContains(t, err, "message does not have solar reading")
	assert.Equal(t, float32(-1.0), radiation)
}

func TestDecodeSolarRadiationNoSensor(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0xFF, 0xC5, 0x00, 0x11, 0x82})

	radiation, err := DecodeSolarRadiation(message)
	assert.ErrorContains(t, err, "no sensor")
	assert.Equal(t, float32(-1.0), radiation)
}

func TestDecodeSolarRadiation0(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x00, 0x45, 0x00, 0xAD, 0x21})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(1.757936), radiation)
}

func TestDecodeSolarRadiation2(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x00, 0x85, 0x00, 0xBB, 0x75})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(3.515872), radiation)
}

func TestDecodeSolarRadiation7(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x01, 0x05, 0x00, 0x97, 0xDD})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(7.031744), radiation)
}

func TestDecodeSolarRadiation11(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x01, 0x85, 0x00, 0x8C, 0x45})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(10.547616), radiation)
}

func TestDecodeSolarRadiation11Alt(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x01, 0xC5, 0x00, 0x81, 0x89})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(12.305552), radiation)
}

func TestDecodeSolarRadiation16(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x02, 0x45, 0x00, 0xC3, 0x41})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(15.8214245), radiation)
}

func TestDecodeSolarRadiation39(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x05, 0x85, 0x00, 0x50, 0x85})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(38.674591), radiation)
}

func TestDecodeSolarRadiation47(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x06, 0xC5, 0x00, 0x04, 0x19})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(47.464272), radiation)
}

func TestDecodeSolarRadiation58(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x08, 0x45, 0x00, 0x04, 0x80})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(58.011887), radiation)
}

func TestDecodeSolarRadiation74(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x08, 0xC5, 0x00, 0x1F, 0x18})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(61.527760), radiation)
}

func TestDecodeSolarRadiation153(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x15, 0xC5, 0x00, 0x1E, 0x2A})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(152.940430), radiation)
}

func TestDecodeSolarRadiation156(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x16, 0x45, 0x00, 0x5C, 0xE2})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(156.456299), radiation)
}

func TestDecodeSolarRadiation163(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x17, 0x45, 0x00, 0x6B, 0xD2})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(163.488052), radiation)
}

func TestDecodeSolarRadiation167(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x17, 0xC5, 0x00, 0x70, 0x4A})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(167.003922), radiation)
}

func TestDecodeSolarRadiation176(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x18, 0xC5, 0x00, 0x5C, 0x7B})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(174.035660), radiation)
}

func TestDecodeSolarRadiation202(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x1C, 0xC5, 0x00, 0x80, 0xBB})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(202.162643), radiation)
}

func TestDecodeSolarRadiation218(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x1E, 0x45, 0x00, 0xF5, 0x43})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(212.710251), radiation)
}

func TestDecodeSolarRadiation221(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x1D, 0x85, 0x00, 0xBA, 0x47})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(207.436447), radiation)
}

func TestDecodeSolarRadiation260(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x25, 0x05, 0x00, 0xCD, 0xDB})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(260.174530), radiation)
}

func TestDecodeSolarRadiation478(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x43, 0xC5, 0x00, 0xF2, 0x44})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(476.400665), radiation)
}

func TestDecodeSolarRadiation513(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x48, 0x45, 0x00, 0x19, 0x2D})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(508.043518), radiation)
}

func TestDecodeSolarRadiation513Alt(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x48, 0xC5, 0x00, 0x02, 0xB5})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(511.559387), radiation)
}

func TestDecodeSolarRadiation610(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x56, 0x85, 0x00, 0x57, 0x1B})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(608.245850), radiation)
}

func TestDecodeSolarRadiation782(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x6E, 0x85, 0x00, 0x3B, 0x1F})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(777.007690), radiation)
}

func TestDecodeSolarRadiation810(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x72, 0x45, 0x00, 0x1B, 0x49})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(803.376770), radiation)
}

func TestDecodeSolarRadiation816(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x74, 0x05, 0x00, 0xA4, 0x25})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(815.682312), radiation)
}

func TestDecodeSolarRadiation830(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x75, 0x45, 0x00, 0x9E, 0xD9})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(824.471985), radiation)
}

func TestDecodeSolarRadiation847(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x79, 0x45, 0x00, 0xEB, 0xB8})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(852.598938), radiation)
}

func TestDecodeSolarRadiation847Alt(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x78, 0x05, 0x00, 0xD1, 0x44})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(843.809265), radiation)
}

func TestDecodeSolarRadiation986(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0x8E, 0x05, 0x00, 0x80, 0xB6})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(-801.618835), radiation)
}

func TestDecodeSolarRadiation1761(t *testing.T) {
	message := createMessage([]byte{0x60, 0x00, 0x00, 0xFA, 0x85, 0x00, 0x9F, 0xE6})

	radiation, err := DecodeSolarRadiation(message)
	assert.NoError(t, err)
	assert.Equal(t, float32(-38.674591), radiation)
}
