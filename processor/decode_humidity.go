package processor

import (
	"errors"
	"log/slog"

	"github.com/nathanmsmith/rtldavis/protocol"
)

func DecodeHumidity(m protocol.Message) (float32, error) {
	// From Dekay (https://github.com/dekay/im-me/blob/master/pocketwx/src/protocol.txt):
	// >Humidity is represented as two bytes in Byte 3 and Byte 4 as a ten bit value.
	// >Bits 5 and 4 in Byte 4 are the two most significant bits.  Byte 3 is the
	// >low order byte. The ten bit value is then 10x the humidity value displayed on
	// >the console.  The function of the four low order bits in Byte 3 that cause the
	// >apparent jitter are not known.
	//
	// http://madscientistlabs.blogspot.com/2012/05/its-not-heat.html
	// https://www.carluccio.de/davis-vue-hacking-part-2/
	// https://github.com/kobuki/VPTools/blob/61e39ac9c561d439939bd8bbe1b9e77b72b7be27/Examples/ISSRx/ISSRx.ino#L156-L158
	// https://github.com/dcbo/ISS-MQTT-Gateway/blob/master/src/main.cpp

	slog.Info("Humidity reading received", "raw_byte_data", bytesToSpacedHex(m.Data))
	if GetMessageType(m) != 0x0A {
		return -1, errors.New("message does not have humidity")
	}

	humidity := float32((int16(m.Data[4]>>4))<<8|int16(m.Data[3])) / 10.0
	return humidity, nil
}
