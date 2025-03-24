package processor

import (
	"github.com/nathanmsmith/rtldavis/protocol"
)

// Decode the windspeed reading from a message.
func DecodeWindSpeed(m protocol.Message) int {
	// From https://github.com/dekay/im-me/blob/master/pocketwx/src/protocol.txt:
	// > Byte 1: Wind speed in mph.  Wind speed is updated every transmission.  Simple.
	return int(m.Data[1])
}
