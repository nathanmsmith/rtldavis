package protocol

import (
	"fmt"
	"log"
)

func DecodeMsg(m Message) (obs []string) {
	/* sensor messages arrive with 'channel' numbers, which has no
	   relation to a go chan or an RF frequency. we only understand the
	   ISS ('integrated sensor suite') channel 0, used by Vantage Vue. */
	if m.ID != 0 {
		log.Printf("received message for unsupported channel %d", m.ID)
		return
	}

	windspeed_raw := m.Data[1]
	/* TODO: weewx rtldavis.py has a giant complicated error correction
	   matrix for windspeed.  But from comments it's not clear it
	   applies to Vantage Vue? */
	obs = append(obs, fmt.Sprintf("windspeed_raw %d", windspeed_raw))

	winddir_vue := float64(m.Data[2])*1.40625 + 0.3
	obs = append(obs, fmt.Sprintf("winddir %.2f", winddir_vue))

	msg_type := (m.Data[0] >> 4) & 0x0F
	/* most of the time we will use the 10-bit number in this weird place */
	raw := ((int16(m.Data[3]) << 2) + int16(m.Data[4])>>6) & 0x03FF
	switch msg_type {
	case 0x02:
		/* supercap voltage */
		obs = append(obs, fmt.Sprintf("supercap_v_raw %d", raw))
		if raw != 0x03FF {
			obs = append(obs, fmt.Sprintf("supercap_v %0.2f", float32(raw)/300.0))
		}
	case 0x04:
		/* UV radiation */
		obs = append(obs, fmt.Sprintf("uv_raw %d", raw))
		if raw != 0x03FF {
			obs = append(obs, fmt.Sprintf("uv %0.2f", raw/50.0))
		}
	case 0x05:
		/* rain rate */
		time_between_tips_raw := ((int16(m.Data[4]) & 0x30) << 4) + int16(m.Data[3])
		rain_rate := 0.0
		if time_between_tips_raw == 0x03FF {
			rain_rate = 0.0 /* time between tips is infinity => no rain */
		} else if m.Data[4]&0x40 == 0 {
			/* "heavy rain", time-between-tips is scaled up by 4 bits */
			/* TODO: we are assuming the 0.1 inch (=2.54mm) size bucket */
			rain_rate = 3600.0 / (float64(time_between_tips_raw) / 16.0) * 2.54
		} else {
			/* "light rain" formula */
			rain_rate = 3600.0 / float64(time_between_tips_raw) * 2.54
		}
		obs = append(obs, fmt.Sprintf("rain_rate_mmh %.2f", rain_rate))
	case 0x06:
		/* solar radiation */
		if raw < 0x03fE {
			sr := float64(raw) * 1.757936
			obs = append(obs, fmt.Sprintf("solar_radiation %.2f", sr))
		}
	case 0x07:
		/* solar panel output */
		if raw != 0x03FF {
			obs = append(obs, fmt.Sprintf("solar_panel_v %.2f", float32(raw)/300.0))
		}
	case 0x08:
		/* temperature */
		raw = (int16(m.Data[3]) << 4) + (int16(m.Data[4]) >> 4)
		if raw != 0x0FFC {
			if m.Data[4]&0x08 != 0 {
				/* digital sensor */
				obs = append(obs, fmt.Sprintf("temp %.2f", float32(raw)/10.0))
			} else {
				/* TODO: thermistor */
				log.Printf("can't interpret analog temperature sensor reading %d", raw)
				obs = append(obs, fmt.Sprintf("temp_raw %d", raw))
			}
		}
	case 0x09:
		/* 10-min average wind gust */
		gust_raw := m.Data[3]
		gust_index := m.Data[5] >> 4 /* ??? */
		obs = append(obs, fmt.Sprintf("gust_mph %d", gust_raw))
		obs = append(obs, fmt.Sprintf("gust_index %d", gust_index))
	case 0x0a:
		/* humidity */
		raw = (int16(m.Data[4]>>4) << 8) + int16(m.Data[3])
		if raw != 0 {
			if m.Data[4]&0x08 != 0 {
				/* digital sensor */
				obs = append(obs, fmt.Sprintf("humidity %.2f", float32(raw)/10.0))
			} else {
				/* TODO: two other types of digital sensor and one analog */
				log.Printf("can't interpret humidity sensor reading %d", raw)
				obs = append(obs, fmt.Sprintf("humidity_raw %d", raw))
			}
		}
	case 0x0e:
		/* rain */
		raw = int16(m.Data[3]) /* "rain count raw"?? */
		/* ignore the high bit because apparently some counters wrap at
		   128 and others at 256, and it doesn't matter */
		if raw != 0x80 {
			raw &= 0x7f
			obs = append(obs, fmt.Sprintf("rain_count %d", raw))
		}
	default:
		log.Printf("Unknown data packet type 0x%02x: %02x", m.ID, m.Data)
	}
	return
}
