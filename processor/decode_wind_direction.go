package processor

import (
	"log/slog"
	"math"

	"github.com/nathanmsmith/rtldavis/protocol"
)

// Decode the windspeed reading from a message.
// TODO: write tests when you're a little more sure about the accuracy of the measurements
func DecodeWindDirection(m protocol.Message) int16 {
	// From Dekay (https://github.com/dekay/im-me/blob/master/pocketwx/src/protocol.txt):
	// > Byte 2: Wind direction from 1 to 360 degrees.  Wind direction is updated every
	// > transmission.  The wind reading is contained in a single byte that limits the
	// > maximum value to 255.  It is converted to a range of 1 to 360 degrees by
	// > scaling the byte value by 360 / 255.  A wind speed reading of 0xd3 = 211
	// > (decimal) * 360 / 255 = 297.
	//
	// However, this is apparently totally out of date for Vantage Vues which apparently use a magnetic sensor instead of
	// a potentiometer for determining wind direction.

	// In his code (https://github.com/lheijst/weewx-rtldavis/blob/master/bin/user/rtldavis.py#L1049-L1059),
	// Luc mentions this as well. He calculates this as:
	// wind_dir_vue = wind_dir_raw * 1.40625 + 0.3
	// But that appears to be for the VP2, not the Vantage Vue! And he doesn't even use this value anywhere.
	// Red herring.
	// https://www.wxforum.net/index.php?topic=21967.msg213307#msg213307
	//
	// Dario (https://www.carluccio.de/davis-vue-hacking-part-2/) uses a modified version of Dekay's solution:
	// WindDir= (9 + Byte2 * 342 / 255)
	// I'm not sure where his came from.
	//
	// Kobuki and rdsman came to a different approach (https://www.wxforum.net/index.php?topic=22189.msg247945#msg247945):
	// Take the two bytes from packet 2 and the bottom byte from packet 4.
	// rdsman notes of the bottom byte:
	// > Bit 0 has no effect on the wind direction. Bit 1 always causes it to read 1 digit higher.
	//
	// Kobuki and rdsman then seem to differ on the constant to multiply this by.
	// https://github.com/kobuki/VPTools/blob/master/Examples/ISSRx/ISSRx.ino#L93
	//
	// 2026-01-01: Empirical results indicate that dekay's data is most accurate for me.

	luc_wind_dir := math.Round(float64(m.Data[2])*1.40625 + 0.3)
	dekay_wind_dir := math.Round(float64(m.Data[2]) * 360 / 255)
	rawDirection := (m.Data[2] << 1) | (m.Data[4]&2)>>1
	kabuki_wind_dir := math.Round(float64(rawDirection) * 360 / 512)
	rdsman_wind_dir := math.Round(float64(rawDirection) * 0.3515625)
	dario_wind_dir := math.Round(9 + float64(m.Data[2])*342/255)

	slog.Info("Parsed wind direction", "luc", luc_wind_dir, "kabuki", kabuki_wind_dir, "rdsman", rdsman_wind_dir, "dekay", dekay_wind_dir, "dario", dario_wind_dir)

	return int16(dekay_wind_dir)
}
