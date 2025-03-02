package protocol

import (
	"log"
)

// Decode the windspeed from the
// Returns an error if there is no temperature reading.
func DecodeWindspeed(m Message) (float32, error) {
	// From https://github.com/dekay/im-me/blob/master/pocketwx/src/protocol.txt:
	// > Byte 1: Wind speed in mph. Wind speed is updated every transmission. Simple.
	//
	// From Luc Heijst:
	// > Each data packet of iss or anemometer contains wind info,
	// > but it is only valid when received from the channel with
	// > the anemometer connected
	// > message examples:
	// > 51 06 B2 FF 73 00 76 61
	// > E0 00 00 4E 05 00 72 61 (no sensor)

	log.Printf("Wind reading received, raw byte data: %x", m.Data)
	windspeed := float32(m.Data[1])

	log.Printf("windspeed: %.2f", windspeed)

	return windspeed, nil
}
