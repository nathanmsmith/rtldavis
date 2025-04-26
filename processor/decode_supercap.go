package processor

import (
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
	// Luc's calculation is different, but he also doesn't have a Vantage Vue.
	// https://github.com/lheijst/weewx-rtldavis/blob/master/bin/user/rtldavis.py#L1098
	//
	// Kobuki uses 300 instead of 100. He also doesn't have a Vantage Vue
	// https://github.com/kobuki/VPTools/blob/61e39ac9c561d439939bd8bbe1b9e77b72b7be27/Examples/ISSRx/ISSRx.ino#L174
	voltage := float32((m.Data[3]<<2)+((m.Data[4]&0xC0)>>6)) / 100

	slog.Info("Parsed supercap voltage", "voltage", voltage)

	return float32(voltage), nil
}
