package processing

import (
	"context"
	"encoding/json"
	"iot-gateway-sim/internal/transport"
	"iot-gateway-sim/pkg/model"
	"log"
	"time"
)

const (
	bufferSize    = 10
	flushInterval = 5 * time.Second
)

func RunBufferPipeline(ctx context.Context, dataChan <-chan model.SensorData, mqttClient *transport.MQTTClient) {
	log.Println("[Processor] Buffer pipeline started")

	buffer := make([]model.SensorData, 0, bufferSize)
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	flush := func(reason string) {
		if len(buffer) == 0 {
			return
		}
		log.Printf("[Processor] %s. Flushing %d data.", reason, len(buffer))
		flushBuffer(buffer, mqttClient)
		buffer = buffer[:0]
	}

	for {
		select {
		case <-ctx.Done():
			flush("Shutdown signal received")
			log.Println("[Processor] Pipeline stopped.")
			return

		case data, ok := <-dataChan:
			if !ok {
				flush("Data channel closed")
				log.Println("[Processor] Pipeline stopped.")
				return
			}

			buffer = append(buffer, data)

			if len(buffer) >= bufferSize {
				flush("Buffer full")
			}

		case <-ticker.C:
			flush("Flust interval reached")
		}
	}
}

func flushBuffer(buffer []model.SensorData, mqttClient *transport.MQTTClient) {
	if len(buffer) == 0 {
		return
	}

	payload, err := json.Marshal(buffer)
	if err != nil {
		log.Printf("[Processor-Flush] Failed to marshaling JSON: %v", err)
		return
	}

	topic := "sensor/batch"
	if err := mqttClient.Publish(topic, payload); err != nil {
		log.Printf("[Processor-Flush] Failed to publish batch to MQTT: %v", err)
	} else {
		log.Printf("[Processor-Flush] Publish %d data to topic %s", len(buffer), topic)
	}
}
