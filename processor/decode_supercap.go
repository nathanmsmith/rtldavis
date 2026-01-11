package processor

import (
	"errors"
	"log/slog"

	"github.com/nathanmsmith/rtldavis/protocol"
)

// Decode the voltage of the super capacitor ("SuperCap").
// From Dekay:
// > The Supercap is a "super" capacitor that stores excess
// > energy from the ISS solar cell during the day, which is
// > then used to help power the console at night.
func DecodeSupercap(m protocol.Message) (float32, error) {
	// Dario's calculation. Source:
	// https://www.carluccio.de/davis-vue-hacking-part-2/
	// Goldcap [v]= ((Byte3 * 4) + ((Byte4 && 0xC0) / 64)) / 100
	//
	// Dekay defers to Dario here.
	// https://github.com/dekay/DavisRFM69/wiki/Message-Protocol
	//
	// Luc, Kobuki, and Matthew Wall all divide by 300 instead of 100.
	// https://github.com/lheijst/weewx-rtldavis/blob/master/bin/user/rtldavis.py#L1098
	// https://github.com/matthewwall/weewx-meteostick/blob/master/bin/user/meteostick.py#L774-L776
	// https://github.com/kobuki/VPTools/blob/61e39ac9c561d439939bd8bbe1b9e77b72b7be27/Examples/ISSRx/ISSRx.ino#L174
	// https://github.com/kobuki/VPTools/issues/13
	//
	// They are right, the max of the super cap is ~3V, not 8.
	// https://support.davisinstruments.com/article/0ics9tab6w-manual-vantage-vue-integrated-sensor-suite-manual-6250-6357
	//
	slog.Info("Supercap reading received", "raw_byte_data", bytesToSpacedHex(m.Data))
	if GetMessageType(m) != 0x02 {
		return -1, errors.New("message does not have supercap")
	}

	voltage := float32((m.Data[3]<<2)|((m.Data[4]&0xC0)>>6)) / 300

	return float32(voltage), nil
}
