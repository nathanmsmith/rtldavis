package processor

import (
	"log/slog"
	"math"

	"github.com/nathanmsmith/rtldavis/protocol"
)

// Decode the windspeed reading from a message.
func DecodeWindDirection(m protocol.Message) int {
	// From https://github.com/dekay/im-me/blob/master/pocketwx/src/protocol.txt:
	// > Byte 2: Wind direction from 1 to 360 degrees.  Wind direction is updated every
	// > transmission.  The wind reading is contained in a single byte that limits the
	// > maximum value to 255.  It is converted to a range of 1 to 360 degrees by
	// > scaling the byte value by 360 / 255.  A wind speed reading of 0xd3 = 211
	// > (decimal) * 360 / 255 = 297.

	// https://github.com/kobuki/VPTools/blob/master/Examples/ISSRx/ISSRx.ino#L93
	// https://www.wxforum.net/index.php?topic=22189.msg213467#msg213467

	luc_wind_dir := float32(m.Data[2])*1.40625 + 0.3
	rawVal := (m.Data[2] << 1) | (m.Data[4]&2)>>1
	kabuki_wind_dir := math.Round(float64(rawVal) * 360 / 512)
	rdsman_wind_dir := math.Round(float64(rawVal) * 0.3515625)
	dekay_wind_dir := float64(m.Data[2]) * 360 / 255

	slog.Info("Parsed wind direction", "luc", luc_wind_dir, "kabuki", kabuki_wind_dir, "rdsman", rdsman_wind_dir, "dekay", dekay_wind_dir)

	return 0
}
