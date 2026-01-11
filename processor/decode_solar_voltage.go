package processor

import (
	"errors"
	"log/slog"

	"github.com/nathanmsmith/rtldavis/protocol"
)

//     # message example:
//     # 70 01 F5 CE 43 86 58 E2

func DecodeSolarVoltage(m protocol.Message) (float32, error) {
	// The Davis Vantage Pro has a UV Index (0x04) and Solar Voltage (0x06) sensor.
	// Dekay has documented both here: https://github.com/dekay/im-me/blob/master/pocketwx/src/protocol.txt
	// The Vantage Vue has neither, but has a solar panel voltage sensor (0x07).
	//
	// Dario documented the reading, but couldn't find the units/value.
	// https://github.com/dcbo/ISS-MQTT-Gateway/blob/master/src/main.cpp
	//
	// Luc, Kobuki, and Matthew Wall all divide by 300 instead of 100.
	// https://github.com/lheijst/weewx-rtldavis/blob/master/bin/user/rtldavis.py#L1098
	// https://github.com/matthewwall/weewx-meteostick/blob/master/bin/user/meteostick.py#L774-L776
	// https://github.com/kobuki/VPTools/blob/61e39ac9c561d439939bd8bbe1b9e77b72b7be27/Examples/ISSRx/ISSRx.ino#L174

	slog.Info("Solar voltage reading received", "raw_byte_data", bytesToSpacedHex(m.Data))
	if GetMessageType(m) != 0x07 {
		return -1, errors.New("message does not have solar voltage reading")
	}

	if m.Data[3] == 0xFF {
		return -1, errors.New("no sensor")
	}

	solarVoltage := float32((m.Data[3] << 2) | ((m.Data[4] & 0xC0) >> 6))
	return solarVoltage, nil
}
