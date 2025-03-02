package protocol

import (
	"fmt"
	"slices"
)

// Barometric pressure sensor - Tracks atmospheric pressure
// Rain gauge - Measures precipitation amounts
// Anemometer - Measures wind speed
// Wind vane - Determines wind direction
// Solar radiation sensor (via calculated values, not a direct sensor)

// Decode the type of a message.
// Returns an error if the message type is not recognized.
func DecodeMessageType(m Message) (messageType byte, err error) {
	// Get the top nibble of the first byte.
	// Types:
	// 0:
	// 1:
	// 2: supercap voltage
	// 3: unknown
	// 4: UV (not on Vantage Vue ISS?)
	// 5: Rain Rate
	// 6: Solar Radiation (not on Vantage Vue ISS?)
	// 7: Solar Cell Output / Solar Power (Vue only)
	// 8: Temperature
	// 9: 10-min average wind gust
	// 10: Humidity
	// 11:
	// 12:
	// 13:
	// 14: Rain
	// 15:
	knownMessageTypes := []byte{0x02, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0E}

	messageType = (m.Data[0] >> 4) & 0x0F
	if !slices.Contains(knownMessageTypes, messageType) {
		err = fmt.Errorf("Unknown message type: %02x", messageType)
	}

	return messageType, err
}
