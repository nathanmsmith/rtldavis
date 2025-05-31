// TODO
// a link: https://www.wxforum.net/index.php?topic=24981.msg244036#msg244036
// https://www.wxforum.net/index.php?topic=18489.msg286136#msg286136

//
// https://www.manula.com/manuals/pws/davis-kb/1/en/topic/uv-sensor

package processor

import "github.com/nathanmsmith/rtldavis/protocol"

func DecodeUVIndex(m protocol.Message) (float32, error) {
	// From Dekay (https://github.com/dekay/DavisRFM69/wiki/Message-Protocol):
	// > Bytes 3 and 4 are for UV Index. The first byte is MSB and the second
	// > LSB. The lower nibble of the 4th byte is always 5, so they only use the
	// > first three nibbles. A value of FF in the third byte indicates that no
	//  sensor is present.

	// https://github.com/cmatteri/CC1101-Weather-Receiver/blob/master/weewx/ccwxrxvp2.py#L121

	// index = (((m.Data[3] << 8) + m.Data[4]) >> 6) / 50.0

	// Take the top 10 bits of the double byte
	rawTenBits := (doublebyte(m.Data[3])<<8 | doublebyte(m.Data[4])) >> 6
	uvIndex := float32(rawTenBits) / 50.0

	// https://www.wxforum.net/index.php?topic=18489.msg190548;topicseen#msg190548
	return uvIndex, nil
}
