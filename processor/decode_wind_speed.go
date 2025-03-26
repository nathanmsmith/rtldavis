package processor

import (
	"github.com/nathanmsmith/rtldavis/protocol"
)

// Decode the windspeed reading from a message.
// TODO: write tests when you're a little more sure about the accuracy of the measurements
func DecodeWindSpeed(m protocol.Message) int16 {
	// From Dekay (https://github.com/dekay/im-me/blob/master/pocketwx/src/protocol.txt):
	// > Byte 1: Wind speed in mph.  Wind speed is updated every transmission.  Simple.
	//
	// Luc's approach uses error correction, but I'm not sure if this applies to the Vantage Vue.
	// https://github.com/lheijst/weewx-rtldavis/blob/master/bin/user/rtldavis.py#L1075
	//
	// Luc also attributes his code to Kobuki, but Kobuki doesn't use it in their code!
	// https://github.com/kobuki/VPTools/blob/master/Examples/ISSRx/ISSRx.ino#L102C24-L102C33
	//
	// Dario also uses the raw byte.
	// https://www.carluccio.de/davis-vue-hacking-part-2/
	return int16(m.Data[1])
}
