package processor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"log/slog"

	"github.com/nathanmsmith/rtldavis/protocol"
)

type WindDatum struct {
	Speed      int16     `json:"speed"`
	Direction  int16     `json:"direction"`
	ReceivedAt time.Time `json:"received_at"`
}

type TemperatureDatum struct {
	Temperature float32   `json:"temperature"`
	ReceivedAt  time.Time `json:"received_at"`
}

type HumidityDatum struct {
	Humidity   float32   `json:"humidity"`
	ReceivedAt time.Time `json:"received_at"`
}

type RainDatum struct {
	InchesPerHour float32   `json:"inches_per_hour"`
	ReceivedAt    time.Time `json:"received_at"`
}

type SolarDatum struct {
	// Temperature *float32  `json:"temperature"`
	ReceivedAt time.Time `json:"received_at"`
}

type UVDatum struct {
	// UVIndex    *float32  `json:"uv_index"`
	ReceivedAt time.Time `json:"received_at"`
}

type WeatherDatum struct {
	Temperature *TemperatureDatum `json:"temperature"`
	Wind        *WindDatum        `json:"wind"`
	Rain        *RainDatum        `json:"rain"`
	SentAt      time.Time         `json:"sent_at"`
}

// POSTs weather data to a server every N seconds
// or when all data is collected.
type WeatherProcessor struct {
	// The data
	data WeatherDatum
	// temperature *TemperatureDatum
	// wind        *WindDatum
	// humidity    *HumidityDatum

	mutex       sync.Mutex
	batchSize   int
	interval    time.Duration
	serverURL   string
	messageChan chan protocol.Message
	done        chan struct{}
}

func NewWeatherProcessor(serverURL string, interval time.Duration, batchSize int) *WeatherProcessor {
	bp := &WeatherProcessor{
		batchSize:   batchSize,
		interval:    interval,
		serverURL:   serverURL,
		messageChan: make(chan protocol.Message, batchSize),
		done:        make(chan struct{}),
	}

	// Start the background processing
	go bp.processMessages()
	go bp.sendDataPeriodically()

	return bp
}

func (wp *WeatherProcessor) hasSomeDataFields() bool {
	return wp.data.Temperature != nil
}

// func (wp *WeatherProcessor) hasAllDataFields() bool {
// 	return wp.data.Temperature != nil
// }

func (wp *WeatherProcessor) processMessages() {
	for {
		select {
		case message := <-wp.messageChan:
			wp.mutex.Lock()

			slog.Info("Processing message", "raw_message", bytesToSpacedHex(message.Data))

			// process wind speed, wind direction
			// set wind speed, wind direction

			windSpeed := DecodeWindSpeed(message)
			windDirection := DecodeWindDirection(message)
			wp.data.Wind = &WindDatum{
				Speed:      windSpeed,
				Direction:  windDirection,
				ReceivedAt: message.ReceivedAt,
			}
			slog.Info("Saved wind data, will send soon", "windspeed", windSpeed, "direction", windDirection)

			switch GetMessageType(message) {

			// Super capacitor voltage
			case 0x02:
				_, err := DecodeSupercap(message)
				if err == nil {
					// wp.data.Capacitor = &RainDatum{
					// 	InchesPerHour: inchesPerHour,
					// 	ReceivedAt:    message.ReceivedAt,
					// }
					// slog.Info("Saved rain rate data, will send soon", "inchesPerHour", inchesPerHour)
				} else {
					slog.Error("Could not decode temperature from packet", "error", err)
				}

			// UV Index
			// https://github.com/dekay/DavisRFM69/wiki/Message-Protocol#message-4-uv-index
			case 0x04:

			// Rain Rate
			case 0x05:
				inchesPerHour, err := DecodeRainRate(message)
				if err == nil {
					wp.data.Rain = &RainDatum{
						InchesPerHour: inchesPerHour,
						ReceivedAt:    message.ReceivedAt,
					}
					slog.Info("Saved rain rate data, will send soon", "inchesPerHour", inchesPerHour)
				} else {
					slog.Error("Could not decode temperature from packet", "error", err)
				}

			// todo: Solar radiation?
			// https://github.com/dekay/DavisRFM69/wiki/Message-Protocol#message-6-solar-radiation
			// Dario says it's 0x07, Dekay 0x06
			// https://www.carluccio.de/davis-vue-hacking-part-2/
			case 0x06, 0x07:
				//    elif message_type == 6:
				//     # solar radiation
				//     # message examples
				//     # 61 00 DB 00 43 00 F4 3B
				//     # 60 00 00 FF C5 00 79 DA (no sensor)
				//     sr_raw = ((pkt[3] << 2) + (pkt[4] >> 6)) & 0x3FF
				//     if sr_raw < 0x3FE:
				//         data['solar_radiation'] = sr_raw * 1.757936
				//         dbg_parse(2, "solar_radiation_raw=0x%04x value=%s"
				//                   % (sr_raw, data['solar_radiation']))
				// elif message_type == 7:
				//     # solar cell output / solar power (Vue only)
				//     # message example:
				//     # 70 01 F5 CE 43 86 58 E2
				//     """When the raw values are divided by 300 the voltage comes
				//     in the range of 2.8-3.3 V measured by the machine readable
				//     format
				//     """
				//     solar_power_raw = ((pkt[3] << 2) + (pkt[4] >> 6)) & 0x3FF
				//     if solar_power_raw != 0x3FF:
				//         data['solar_power'] = solar_power_raw / 300.0
				//         dbg_parse(2, "solar_power_raw=0x%03x solar_power=%s"
				//                   % (solar_power_raw, data['solar_power']))
				//

			// Rain clicks
			case 0x0E:
				DecodeRainfall(message)

			// Temperature
			case 0x08:
				temperature, err := DecodeTemperature(message)
				if err == nil {
					wp.data.Temperature = &TemperatureDatum{
						Temperature: temperature,
						ReceivedAt:  message.ReceivedAt,
					}
					slog.Info("Saved temperature data, will send soon", "temp", temperature)
				} else {
					slog.Error("Could not decode temperature from packet", "error", err)
				}

			// gust speed Msg-ID 0x9 (every 50 seconds):

			// outside humidity Msg-ID 0xA (every 50 seconds):

			default:
				slog.Info("Unknown message type", "raw_message", bytesToSpacedHex(message.Data), "message_type", GetMessageType(message))
			}

			wp.mutex.Unlock()
		case <-wp.done:
			return
		}
	}
}

func (bp *WeatherProcessor) sendDataPeriodically() {
	ticker := time.NewTicker(bp.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			bp.mutex.Lock()
			if bp.hasSomeDataFields() {
				bp.sendData()
			}
			bp.mutex.Unlock()
		case <-bp.done:
			return
		}
	}
}

func (wp *WeatherProcessor) sendData() {
	if !wp.hasSomeDataFields() {
		return
	}

	wp.data.SentAt = time.Now()
	payload, err := json.Marshal(wp.data)
	if err != nil {
		slog.Error("Error marshaling data to JSON", "error", err)
		return
	}

	// Send the HTTP POST request
	resp, err := http.Post(wp.serverURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		slog.Error("Error POSTing data", "error", err, "payload", payload)
		return
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			slog.Error("Error closing response body", "error", closeErr)
		}
	}()

	slog.Info("Successfully POSTed weather data", "payload", payload)

	if resp.StatusCode != http.StatusCreated {
		slog.Error("Server returned non-Created status", "status", resp.Status)
		return
	}

	wp.clearData()
}

// Clear out any existing weather data.
func (wp *WeatherProcessor) clearData() {
	slog.Info("Clearing data")
	wp.data = WeatherDatum{}
}

func (wp *WeatherProcessor) AddMessage(message protocol.Message) {
	wp.messageChan <- message
}

func (bp *WeatherProcessor) Stop() {
	close(bp.done)

	// Send any remaining messages
	bp.mutex.Lock()
	defer bp.mutex.Unlock()
	if bp.hasSomeDataFields() {
		bp.sendData()
	}
}

// func main() {
// 	// Example usage
// 	processor := NewBatchProcessor(
// 		"http://example.com/api/batch",
// 		5*time.Second, // Send every 5 seconds
// 		100,           // or when batch size reaches 100
// 	)
//
// 	// Example of adding messages
// 	go func() {
// 		for i := 0; i < 1000; i++ {
// 			processor.AddMessage("test message")
// 			time.Sleep(100 * time.Millisecond)
// 		}
// 	}()
//
// 	// Run for a minute then stop
// 	time.Sleep(1 * time.Minute)
// 	processor.Stop()
// }
