package processor

import (
	"errors"

	"log/slog"

	"github.com/nathanmsmith/rtldavis/protocol"
)

func DecodeRainRate(m protocol.Message) (float32, error) {
	// From Dekay (https://github.com/dekay/DavisRFM69/wiki/Message-Protocol):

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

	// > Byte 3 and 4 are temperature.  The first byte is MSB and the second LSB.  The
	// > value is signed with 0x0000 representing 0F.  This reading in the old version
	// > of the ISS was taked from an analog sensor and measured by an A/D.  The newer
	// > ISS uses a digital sensor but still represents the data in the same way.  160
	// > counts (0xa0) represents 1 degree F.  A message of 80 04 70 0f 99 00 91 11
	// > represents temperature as 0x0f99, or 3993 decimal.  Divide 3993 by 160 to get
	// > the console reading of 25.0F
	//
	// Luc Heijst's approach appears to be slightly different but return the same values.
	// He shifts the packet bit into a 12-bit value. Temperature is reported to a tenth
	// of a degree, so we divide by 10 to get the value in Fahrenheit.
	//
	// Examples

	slog.Info("Rain rate reading received", "raw_byte_data", bytesToSpacedHex(m.Data))
	if GetMessageType(m) != 0x05 {
		return -1, errors.New("Message does not have rain rate")
	}

	if m.Data[3] == 0xFF {
		slog.Info("No rain detected", "raw_byte_data", bytesToSpacedHex(m.Data))
	}

	return 0, nil
}
