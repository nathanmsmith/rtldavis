package processor

import (
	"log/slog"

	"github.com/nathanmsmith/rtldavis/protocol"
)

// I'm not totally sure how rainfall and rain rate differ yet.
func DecodeRainfall(m protocol.Message) {
	rainClicks := (m.Data[3] & 0x7F)
	slog.Info("Rain click reading received", "raw_byte_data", bytesToSpacedHex(m.Data), "rainClicks", rainClicks)
}
