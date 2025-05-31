package processor

import "github.com/nathanmsmith/rtldavis/protocol"

// From Davis:
// > The solar radiation reading gives a measure of the amount of solar radiation hitting the solar radiation sensor at any given time, expressed in Watts/sq. meter (W/m^2).
// https://www.davisinstruments.com/pages/solar-radiation-an-explanation
func DecodeSolarRadiation(m protocol.Message) (float32, error) {
	rawTenBits := (doublebyte(m.Data[3])<<8 | doublebyte(m.Data[4])) >> 6
	solarRadiation := float32(rawTenBits) * 1.757936

	return solarRadiation, nil
}
