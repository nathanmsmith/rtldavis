package processor

import (
	"errors"
	"log/slog"

	"github.com/nathanmsmith/rtldavis/protocol"
)

func DecodeSolarRadiation(m protocol.Message) (float32, error) {
	// The Davis Vantage Pro has a UV Index (0x04) and Solar Radiation (0x06) sensor.
	// Dekay has documented both here: https://github.com/dekay/im-me/blob/master/pocketwx/src/protocol.txt
	// The Vantage Vue has neither, but has a solar panel voltage sensor (0x07).
	//
	// Dario documented the reading, but couldn't find the units/value.
	// https://github.com/dcbo/ISS-MQTT-Gateway/blob/master/src/main.cpp
	//
	// Kobuki divides the result by 300. The unit is still unclear
	// https://github.com/kobuki/VPTools/blame/master/Examples/ISSRx/ISSRx.ino#L178-L181
	//
	// Another repo divides the result by 100:
	// https://github.com/HydroSense/FeatherM0_Davis_ISS_rx/blob/4d0fa1d2adae59eabcc53f4300e2acf0090b87bd/FeatherM0_Davis_ISS_rx.ino#L235

	slog.Info("Solar voltage reading received", "raw_byte_data", bytesToSpacedHex(m.Data))
	if GetMessageType(m) != 0x07 {
		return -1, errors.New("message does not have solar voltage reading")
	}

	if m.Data[3] == 0xFF {
		return -1, errors.New("no sensor")
	}

	solarRadiation := float32((m.Data[3] << 2) | ((m.Data[4] & 0xC0) >> 6))
	return solarRadiation, nil
}
