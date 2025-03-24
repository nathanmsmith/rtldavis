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
	Humidity   *float32  `json:"humidity"`
	ReceivedAt time.Time `json:"received_at"`
}

// type RainDatum struct {
// 	Temperature *float32  `json:"temperature"`
// 	ReceivedAt  time.Time `json:"received_at"`
// }

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
			}

			// TODO: other measurements
			// if UV, set UV
			// if humidity, set humidity
			// rain rate
			// etc

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
	defer resp.Body.Close()

	slog.Info("Successfully POSTed weather data", "payload", payload)

	if resp.StatusCode != http.StatusOK {
		slog.Error("Server returned non-OK status", "status", resp.Status)
		return
	}

	wp.clearData()
}

// Clear out any existing weather data.
func (wp *WeatherProcessor) clearData() {
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
