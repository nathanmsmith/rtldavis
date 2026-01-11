package processor

import (
	"errors"
	"log/slog"

	"github.com/nathanmsmith/rtldavis/protocol"
)

func DecodeSolarRadiation(m protocol.Message) (float32, error) {
	// From Dekay (https://github.com/dekay/im-me/blob/master/pocketwx/src/protocol.txt):
	// > Bytes 3 and 4 are solar radiation.  The first byte is MSB and the second LSB.
	// > The lower nibble of the 4th byte is again always 5, so they only use the first
	// > three nibbles.  A value of FF in the third byte indicates that no sensor is
	// > present.
	// > Solar radiation = (byte3 << 8 + byte4) >> 6) * 1.757936
	// >
	// > Reference:
	// > http: //www.wxforum.net/index.php?topic=18489.msg178506#msg178506
	// > Reference:
	// > http: //www.wxforum.net/index.php?topic=18489.msg190548#msg190548
	//
	// Dario says the bit is 0x07
	// https://github.com/dcbo/ISS-MQTT-Gateway/blob/master/src/main.cpp

	slog.Info("Solar reading received", "raw_byte_data", bytesToSpacedHex(m.Data))
	if GetMessageType(m) != 0x06 {
		return -1, errors.New("message does not have solar reading")
	}

	if m.Data[3] == 0xFF {
		return -1, errors.New("no sensor")
	}

	radiation := float32(((int16(m.Data[3])<<8)+int16(m.Data[4]))>>6) * 1.757936
	return radiation, nil
}
