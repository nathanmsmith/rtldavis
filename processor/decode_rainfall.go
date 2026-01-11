package processor

import (
	"errors"
	"log/slog"

	"github.com/nathanmsmith/rtldavis/protocol"
)

func DecodeRainfall(m protocol.Message) (int16, error) {
	slog.Info("Rainfall reading received", "raw_byte_data", bytesToSpacedHex(m.Data))

	if GetMessageType(m) != 0x0E {
		return -1, errors.New("message does not have rainfall")
	}

	rainClicks := int16(m.Data[3] & 0x7F)
	return rainClicks, nil
}
