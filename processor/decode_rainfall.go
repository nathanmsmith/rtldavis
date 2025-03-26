package processor

import (
	"github.com/nathanmsmith/rtldavis/protocol"
)

func DecodeRainfall(m protocol.Message)  {
	rainClicks := (m.Data[3] & 0x7F)
	slog.Info("Rain click reading received", , "raw_byte_data", bytesToSpacedHex(m.Data)))
}
