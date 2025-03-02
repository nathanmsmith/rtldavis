package protocol

import (
	"errors"
	"log"
)

// Decode the temperature reading from a message.
// Returns an error if there is no temperature reading.
func DecodeTemperature(m Message) (float32, error) {
	// From https://github.com/dekay/im-me/blob/master/pocketwx/src/protocol.txt:
	// > Byte 3 and 4 are temperature.  The first byte is MSB and the second LSB.  The
	// > value is signed with 0x0000 representing 0F.  This reading in the old version
	// > of the ISS was taked from an analog sensor and measured by an A/D.  The newer
	// > ISS uses a digital sensor but still represents the data in the same way.  160
	// > counts (0xa0) represents 1 degree F.  A message of 80 04 70 0f 99 00 91 11
	// > represents temperature as 0x0f99, or 3993 decimal.  Divide 3993 by 160 to get
	// > the console reading of 25.0F
	//
	// Luc Heijst's approach appears to be slightly different but return the same values.
	// He shifts the packet bit into a 12-bit value. Temperature is reported to a tenth
	// of a degree, so we divide by 10 to get the value in Fahrenheit.
	//
	// Examples
	// # 80 00 00 33 8D 00 25 11 (digital temp)
	// # 81 00 00 59 45 00 A3 E6 (analog temp)
	// # 81 00 DB FF C3 00 AB F8 (no sensor)

	log.Printf("Temperature reading received, raw byte data: %x", m.Data)
	if GetMessageType(m) != 0x08 {
		return -1, errors.New("Message does not have temperature")
	}

	if m.Data[4]&0x08 == 0 {
		return -1, errors.New("Temperature reading is not from digital sensor. Analog sensor not supported")
	}

	raw := (int16(m.Data[3]) << 4) + (int16(m.Data[4]) >> 4)
	if raw == 0x0FFC {
		return -1, errors.New("No sensor")
	}

	temperature := float32(raw) / 10.0
	return temperature, nil
}
