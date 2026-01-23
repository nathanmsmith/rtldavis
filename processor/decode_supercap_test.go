package processor

import (
	"testing"
	"time"

	"github.com/nathanmsmith/rtldavis/dsp"
	"github.com/nathanmsmith/rtldavis/protocol"
	"github.com/stretchr/testify/assert"
)

// TestBatteryLowFlagSet verifies that when bit 3 of byte 0 is set,
// the BatteryLow flag is correctly extracted and populated in BatteryDatum
func TestBatteryLowFlagSet(t *testing.T) {
	// Message byte 0 structure: [MessageType(4 bits)][BatteryLow(1 bit)][TransmitterID(3 bits)]
	// For supercap (type 0x02) with battery low = true and transmitter ID = 0:
	// Binary: 0010 1000 = 0x28
	//         ^^^^        Message type 0x02
	//              ^      Battery low bit set
	//               ^^^   Transmitter ID 0

	// protocol.NewMessage skips first 2 bytes, so pkt.Data[2] becomes m.Data[0]
	pkt := dsp.Packet{
		Idx: 0,
		Data: []byte{
			0x00, 0x00,       // Skipped by NewMessage
			0x28,             // Message type 0x02, battery low=true, ID=0
			0x00,             // Wind speed
			0x00,             // Wind direction
			0xE1,             // Supercap voltage byte 3
			0x00,             // Supercap voltage byte 4
			0x00, 0x00,       // Padding
		},
	}

	message := protocol.NewMessage(pkt)

	// Verify the battery low flag was extracted correctly
	assert.True(t, message.BatteryLow, "BatteryLow flag should be true")

	// Verify message type is supercap (0x02)
	assert.Equal(t, byte(0x02), GetMessageType(message))

	// Decode supercap voltage (should succeed)
	voltage, err := DecodeSupercap(message)
	assert.NoError(t, err)
	assert.Greater(t, voltage, float32(0.0), "Voltage should be positive")
}

// TestBatteryLowFlagNotSet verifies that when bit 3 of byte 0 is clear,
// the BatteryLow flag is false
func TestBatteryLowFlagNotSet(t *testing.T) {
	// For supercap (type 0x02) with battery low = false and transmitter ID = 0:
	// Binary: 0010 0000 = 0x20
	//         ^^^^        Message type 0x02
	//              ^      Battery low bit clear
	//               ^^^   Transmitter ID 0

	pkt := dsp.Packet{
		Idx: 0,
		Data: []byte{
			0x00, 0x00,       // Skipped by NewMessage
			0x20,             // Message type 0x02, battery low=false, ID=0
			0x00,             // Wind speed
			0x00,             // Wind direction
			0xE1,             // Supercap voltage byte 3
			0x00,             // Supercap voltage byte 4
			0x00, 0x00,       // Padding
		},
	}

	message := protocol.NewMessage(pkt)

	// Verify the battery low flag was extracted correctly
	assert.False(t, message.BatteryLow, "BatteryLow flag should be false")

	// Verify message type is supercap (0x02)
	assert.Equal(t, byte(0x02), GetMessageType(message))

	// Decode supercap voltage (should succeed)
	voltage, err := DecodeSupercap(message)
	assert.NoError(t, err)
	assert.Greater(t, voltage, float32(0.0), "Voltage should be positive")
}

// TestBatteryDatumPopulatesIsLow verifies that the WeatherProcessor
// correctly populates the IsLow field in BatteryDatum
func TestBatteryDatumPopulatesIsLow(t *testing.T) {
	// Create a weather processor
	wp := NewWeatherProcessor("http://localhost:8080", "test-key", 5_000_000_000, 10)
	defer wp.Stop()

	// Test with battery low = true
	pktLow := dsp.Packet{
		Idx: 0,
		Data: []byte{
			0x00, 0x00,       // Skipped by NewMessage
			0x28,             // Message type 0x02, battery low=true, ID=0
			0x00,             // Wind speed
			0x00,             // Wind direction
			0xE1,             // Supercap voltage byte 3
			0x00,             // Supercap voltage byte 4
			0x00, 0x00,       // Padding
		},
	}
	messageLow := protocol.NewMessage(pktLow)
	wp.AddMessage(messageLow)

	// Give the goroutine time to process the message
	// In a real scenario, you might use a more sophisticated synchronization mechanism
	// but for tests, a small sleep is acceptable
	time.Sleep(10 * time.Millisecond)

	// Access the battery datum
	wp.mutex.Lock()
	batteryLow := wp.data.Battery
	wp.mutex.Unlock()

	assert.NotNil(t, batteryLow, "Battery datum should be populated")
	assert.True(t, batteryLow.IsLow, "IsLow should be true")

	// Test with battery low = false
	pktOk := dsp.Packet{
		Idx: 0,
		Data: []byte{
			0x00, 0x00,       // Skipped by NewMessage
			0x20,             // Message type 0x02, battery low=false, ID=0
			0x00,             // Wind speed
			0x00,             // Wind direction
			0xE1,             // Supercap voltage byte 3
			0x00,             // Supercap voltage byte 4
			0x00, 0x00,       // Padding
		},
	}
	messageOk := protocol.NewMessage(pktOk)
	wp.AddMessage(messageOk)

	// Give the goroutine time to process the message
	time.Sleep(10 * time.Millisecond)

	// Access the battery datum
	wp.mutex.Lock()
	batteryOk := wp.data.Battery
	wp.mutex.Unlock()

	assert.NotNil(t, batteryOk, "Battery datum should be populated")
	assert.False(t, batteryOk.IsLow, "IsLow should be false")
}
