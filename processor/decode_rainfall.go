package processor

import (
	"errors"
	"log/slog"

	"github.com/nathanmsmith/rtldavis/protocol"
)

func DecodeRainfall(m protocol.Message) (int16, error) {
	// From Dekay (https://github.com/dekay/DavisRFM69/wiki/Message-Protocol):
	// > It is a running total of bucket tips that wraps back around to 0 eventually from the ISS. It is up to the console to keep track of changes in this byte. Only bits 0 through 6 of byte 3 are used, so the counter will overflow after 0x7F (127).
	//
	// Dario concurs: https://www.carluccio.de/davis-vue-hacking-part-2/
	//
	// Kobuki uses 0x80 instead of 0x7F:
	// https://github.com/kobuki/VPTools/blob/61e39ac9c561d439939bd8bbe1b9e77b72b7be27/Examples/ISSRx/ISSRx.ino#L126-L132
	//
	// From Luc:
	// > We have seen rain counters wrap around at 127 and others wrap around at 255.  When we filter the highest bit, both counter types will wrap at 127.

	slog.Info("Rainfall reading received", "raw_byte_data", bytesToSpacedHex(m.Data))

	if GetMessageType(m) != 0x0E {
		return -1, errors.New("message does not have rainfall")
	}

	rainClicks := int16(m.Data[3] & 0x7F)
	return rainClicks, nil
}
