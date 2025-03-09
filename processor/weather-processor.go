package processor

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/nathanmsmith/rtldavis/protocol"
)

type WindDatum struct {
	WindSpeed     float32   `json:"wind_speed"`
	WindDirection float64   `json:"wind_direction"`
	ReceivedAt    time.Time `json:"received_at"`
}

type TemperatureDatum struct {
	Temperature *float32  `json:"temperature"`
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
	Temperature *float32  `json:"temperature"`
	ReceivedAt  time.Time `json:"received_at"`
}

type UVDatum struct {
	UVIndex    *float32  `json:"uv_index"`
	ReceivedAt time.Time `json:"received_at"`
}

type WeatherDatum struct {
	SentAt time.Time `json:"sent_at"`
}

// POSTs weather data to a server every N seconds
// or when all data is collected.
type WeatherProcessor struct {
	// The data
	temperature *TemperatureDatum
	wind        *WindDatum
	humidity    *HumidityDatum

	mutex       sync.Mutex
	batchSize   int
	interval    time.Duration
	serverURL   string
	messageChan chan protocol.Message
	done        chan struct{}
}

func NewBatchProcessor(serverURL string, interval time.Duration, batchSize int) *WeatherProcessor {
	bp := &WeatherProcessor{
		batchSize:   batchSize,
		interval:    interval,
		serverURL:   serverURL,
		messageChan: make(chan protocol.Message, batchSize),
		done:        make(chan struct{}),
	}

	// Start the background processing
	go bp.processMessages()
	go bp.sendBatchPeriodically()

	return bp
}

func (wp *WeatherProcessor) hasSomeDataFields() bool {
	return wp.temperature != nil
}

func (wp *WeatherProcessor) hasAllDataFields() bool {
	return wp.temperature != nil
}

func (bp *WeatherProcessor) processMessages() {
	for {
		select {
		case message := <-bp.messageChan:
			bp.mutex.Lock()

			// process wind speed, wind direction
			// set wind speed, wind direction

			// get message type
			// switch
			// if UV, set UV
			// if temperature, set temperature
			// if humidity, set humidity
			// etc

			if bp.hasAllDataFields() {
				bp.sendData()
			}
			bp.mutex.Unlock()

		case <-bp.done:
			return
		}
	}
}

func (bp *WeatherProcessor) sendBatchPeriodically() {
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

func (bp *WeatherProcessor) sendData() {
	if !bp.hasNoDataFields() {
		return
	}

	// Prepare the payload

	payload, err := json.Marshal(map[string]interface{}{
		"packets":         bp.messages,
		"count":           len(bp.messages),
		"payload_sent_at": time.Now(),
	})
	if err != nil {
		log.Printf("Error marshaling messages: %v", err)
		return
	}

	// Send the HTTP POST request
	resp, err := http.Post(bp.serverURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Error sending batch: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Successfully sent messages: %s", payload)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Server returned non-OK status: %d", resp.StatusCode)
		return
	}

	// Clear the messages slice after successful send
	// TODO: bp.clearData()
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
