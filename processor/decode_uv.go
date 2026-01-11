// TODO
// a link: https://www.wxforum.net/index.php?topic=24981.msg244036#msg244036
// https://www.wxforum.net/index.php?topic=18489.msg286136#msg286136

// https://www.davisinstruments.com/pages/solar-radiation-an-explanation
// https://www.manula.com/manuals/pws/davis-kb/1/en/topic/uv-sensor

package processor

import (
	"errors"
	"log/slog"

	"github.com/nathanmsmith/rtldavis/protocol"
)

// 40 00 00 FF C5 00 4C D2 -- Off (unplugged sensor)
// 40 00 00 00 45 00 98 29 -- 0
// 40 00 00 0B C5 00 73 40 -- 0.9
// 40 00 00 0B 05 00 65 14 -- 0.9
// 40 00 00 11 05 00 E1 B6 -- 1.3
// 40 00 00 12 85 00 A3 7E -- 1.5
// 40 00 00 12 45 00 B5 2A -- 1.4
// 40 00 00 11 45 00 EC 7A -- 1.4
// 40 00 00 11 C5 00 F7 E2 -- 1.4
// 40 00 00 0F 45 00 B4 18 -- 1.2
// 40 00 00 0D C5 00 C1 E0 -- 1.1
// 40 00 00 09 45 00 06 B8 -- 0.7
// 40 00 00 19 05 00 48 17 -- 2.0

func DecodeUVRadiation(m protocol.Message) (float32, error) {
	slog.Info("UV reading received", "raw_byte_data", bytesToSpacedHex(m.Data))
	if GetMessageType(m) != 0x04 {
		return -1, errors.New("message does not have uv")
	}

	uvIndex := float32((int16(m.Data[3])<<8)+int16(m.Data[4])) / 50.0
	return uvIndex, nil
}
