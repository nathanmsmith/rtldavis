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
	Value      float32   `json:"value"`
	ReceivedAt time.Time `json:"received_at"`
}

type HumidityDatum struct {
	Value      float32   `json:"value"`
	ReceivedAt time.Time `json:"received_at"`
}

type RainDatum struct {
	InchesPerHour float32   `json:"inches_per_hour"`
	ReceivedAt    time.Time `json:"received_at"`
}

type BatteryDatum struct {
	Voltage    float32   `json:"voltage"`
	ReceivedAt time.Time `json:"received_at"`
}

type SolarDatum struct {
	Voltage    float32   `json:"voltage"`
	ReceivedAt time.Time `json:"received_at"`
}

type WeatherDatum struct {
	Temperature *TemperatureDatum `json:"temperature"`
	Wind        *WindDatum        `json:"wind"`
	Rain        *RainDatum        `json:"rain"`
	Humidity    *HumidityDatum    `json:"humidity"`

	Battery *BatteryDatum `json:"battery"`
	Solar   *SolarDatum   `json:"solar"`

	SentAt time.Time `json:"sent_at"`
}

// POSTs weather data to a server every N seconds
// or when all data is collected.
type WeatherProcessor struct {
	data        WeatherDatum
	mutex       sync.Mutex
	batchSize   int
	interval    time.Duration
	serverURL   string
	apiKey      string
	messageChan chan protocol.Message
	done        chan struct{}
}

func NewWeatherProcessor(serverURL string, apiKey string, interval time.Duration, batchSize int) *WeatherProcessor {
	bp := &WeatherProcessor{
		batchSize:   batchSize,
		interval:    interval,
		serverURL:   serverURL,
		apiKey:      apiKey,
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
				voltage, err := DecodeSupercap(message)
				if err == nil {
					wp.data.Battery = &BatteryDatum{
						Voltage:    voltage,
						ReceivedAt: message.ReceivedAt,
					}
					slog.Info("Saved super capacitor data, will send soon", "voltage", voltage)
				} else {
					slog.Error("Could not decode temperature from packet", "error", err)
				}

			// UV Index
			// https://github.com/dekay/DavisRFM69/wiki/Message-Protocol#message-4-uv-index
			case 0x04:
				slog.Error("Detected a UV Index reading. This is unexpected!!")

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

			case 0x06:
				slog.Error("Detected a solar radiation reading. This is unexpected!!")

			// todo: Solar radiation?
			// https://github.com/dekay/DavisRFM69/wiki/Message-Protocol#message-6-solar-radiation
			// Dario says it's 0x07, Dekay 0x06
			// https://www.carluccio.de/davis-vue-hacking-part-2/
			case 0x07:
				voltage, err := DecodeSolarVoltage(message)
				if err == nil {
					wp.data.Solar = &SolarDatum{
						Voltage:    voltage,
						ReceivedAt: message.ReceivedAt,
					}
					slog.Info("Saved solar voltage data, will send soon", "voltage", voltage)
				} else {
					slog.Error("Could not decode temperature from packet", "error", err)
				}

			// Temperature
			case 0x08:
				temperature, err := DecodeTemperature(message)
				if err == nil {
					wp.data.Temperature = &TemperatureDatum{
						Value:      temperature,
						ReceivedAt: message.ReceivedAt,
					}
					slog.Info("Saved temperature data, will send soon", "temp", temperature)
				} else {
					slog.Error("Could not decode temperature from packet", "error", err)
				}

			// gust speed Msg-ID 0x9 (every 50 seconds):

			// Humidity (every 50 seconds)
			case 0x0A:
				humidity, err := DecodeHumidity(message)
				if err == nil {
					wp.data.Humidity = &HumidityDatum{
						Value:      humidity,
						ReceivedAt: message.ReceivedAt,
					}
					slog.Info("Saved humidity data, will send soon", "temp", humidity)
				} else {
					slog.Error("Could not decode humidity from packet", "error", err)
				}

			// Rain clicks
			case 0x0E:
				DecodeRainfall(message)

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

	// Create the HTTP POST request
	req, err := http.NewRequest("POST", wp.serverURL, bytes.NewBuffer(payload))
	if err != nil {
		slog.Error("Error creating POST request", "error", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", wp.apiKey)

	// Send the HTTP POST request
	client := &http.Client{}
	resp, err := client.Do(req)
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
