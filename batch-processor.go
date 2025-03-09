package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/nathanmsmith/rtldavis/protocol"
)

type BatchProcessor struct {
	packets    []protocol.DecodedPacket
	mutex      sync.Mutex
	batchSize  int
	interval   time.Duration
	serverURL  string
	packetChan chan protocol.DecodedPacket
	done       chan struct{}
}

func NewBatchProcessor(serverURL string, interval time.Duration, batchSize int) *BatchProcessor {
	bp := &BatchProcessor{
		packets:    make([]protocol.DecodedPacket, 0),
		batchSize:  batchSize,
		interval:   interval,
		serverURL:  serverURL,
		packetChan: make(chan protocol.DecodedPacket, batchSize),
		done:       make(chan struct{}),
	}

	// Start the background processing
	go bp.processMessages()
	go bp.sendBatchPeriodically()

	return bp
}

func (bp *BatchProcessor) processMessages() {
	for {
		select {
		case packet := <-bp.packetChan:
			bp.mutex.Lock()
			bp.packets = append(bp.packets, packet)

			// If we've reached batch size, send immediately
			if len(bp.packets) >= bp.batchSize {
				bp.sendBatch()
			}
			bp.mutex.Unlock()

		case <-bp.done:
			return
		}
	}
}

func (bp *BatchProcessor) sendBatchPeriodically() {
	ticker := time.NewTicker(bp.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			bp.mutex.Lock()
			if len(bp.packets) > 0 {
				bp.sendBatch()
			}
			bp.mutex.Unlock()
		case <-bp.done:
			return
		}
	}
}

func (bp *BatchProcessor) sendBatch() {
	if len(bp.packets) == 0 {
		return
	}

	// Prepare the payload
	payload, err := json.Marshal(map[string]interface{}{
		"packets":         bp.packets,
		"count":           len(bp.packets),
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
	bp.packets = make([]protocol.DecodedPacket, 0)
}

func (bp *BatchProcessor) AddPacket(packet protocol.DecodedPacket) {
	bp.packetChan <- packet
}

func (bp *BatchProcessor) Stop() {
	close(bp.done)

	// Send any remaining messages
	bp.mutex.Lock()
	defer bp.mutex.Unlock()
	if len(bp.packets) > 0 {
		bp.sendBatch()
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
