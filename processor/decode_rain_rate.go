package processor

import (
	"errors"

	"log/slog"

	"github.com/nathanmsmith/rtldavis/protocol"
)

func DecodeRainRate(m protocol.Message) (float32, error) {
	// From Dekay (https://github.com/dekay/DavisRFM69/wiki/Message-Protocol):
	// > Bytes 3 and 4 contain the rain rate information. The rate is actually the time in seconds between rain bucket tips in the ISS.
	// > The rain rate is calculated from the bucket tip rate and the size of the bucket (0.01" of rain for units sold in North America).
	// > More information here: https://www.wxforum.net/index.php?topic=23652.msg230631#msg230631

	// if (s[0]== 5):
	//   if (s[3]==255):
	//     rainrate = 0
	//   else:
	//     raw = (((s[4] & 0x30 ) / 16 * 250) + s[3])
	//     if (s[4] & 0x40) == 0x40:
	//       rainrate = 914.4 / raw
	//     elif (s[4] & 0x40) == 0x0:
	//       rainrate = (914.4*16) / raw
	//     else:
	//       print("rainrate fail")

	// > Rain is in Byte 3.  It is a running total of bucket tips that wraps back
	// > around to 0 eventually from the ISS.  It is up to the console to keep track of
	// > changes in this byte.  The example below is bound to confuse: the leading
	// > value is the elapsed time since data collection started (in seconds), all
	// > bytes have been converted to decimal, and the last two CRC bytes have been
	// > stripped off.  A tip of the rain bucket causes the value the ISS is sending
	// > from a steady value of 40 to a new value of 41.
	//
	// > 2426.3,224,16,33,40,1,0
	// > 2436.6,224,11,36,40,1,0
	// > 2446.8,224,9,29,41,2,0
	// > 2457.1,224,10,29,41,3,0

	//
	// Luc Heijst's approach appears to be slightly different but return the same values.
	// He shifts the packet bit into a 12-bit value. Temperature is reported to a tenth
	// of a degree, so we divide by 10 to get the value in Fahrenheit.
	//
	// Examples

	// https://github.com/dcbo/ISS-MQTT-Gateway
	// https://github.com/dcbo/ISS-MQTT-Gateway/blob/1ea7bab1e7c05f49519e7f18509698e05dc9ef04/src/main.cpp#L650

	slog.Info("Rain rate reading received", "raw_byte_data", bytesToSpacedHex(m.Data))
	if GetMessageType(m) != 0x05 {
		return -1, errors.New("Message does not have rain rate")
	}

	if m.Data[3] == 0xFF {
		slog.Info("No rain detected")

		// If the Most Significant Nibble of byte 4 is < 4 (i.e., 0100)
		// Light rain
	} else if m.Data[4]&0x40 == 0 {
		slog.Info("Light rain detected")
		// If the Most Significant Nibble of byte 4 is >= 4 (i.e., 0100)
		// Heavy rain
	} else if m.Data[4]&0x40 == 0x40 {
		slog.Info("Heavy rain detected")
	} else {
		// error
	}

	return 0, nil
}
