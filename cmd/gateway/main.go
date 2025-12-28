package main

import (
	"context"
	"iot-gateway-sim/internal/ingestion"
	"iot-gateway-sim/internal/processing"
	"iot-gateway-sim/internal/transport"
	"iot-gateway-sim/pkg/model"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	dataCh := make(chan model.SensorData, 100)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}

	mqttClient, err := transport.NewMQTTClient("ssl://localhost:8883", "iot-gateway-1",
		"gateway-certs/ca.crt",
		"gateway-certs/client.crt",
		"gateway-certs/client.key",
	)
	if err != nil {
		log.Fatalf("FATAL: MQTT auth failed: %v", err)
	}
	defer mqttClient.Disconnect()
	startPipeline(ctx, wg, dataCh, mqttClient)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	cancel()
	wg.Wait()
}

func startPipeline(
	ctx context.Context,
	wg *sync.WaitGroup,
	dataCh chan model.SensorData,
	mqttClient *transport.MQTTClient,
) {
	wg.Add(3)

	wg.Go(func() {
		ingestion.StartTCPServer(ctx, ":5000", dataCh)
	})

	wg.Go(func() {
		ingestion.StartUDPServer(ctx, ":5001", dataCh)
	})

	wg.Go(func() {
		processing.RunBufferPipeline(ctx, dataCh, mqttClient)
	})

}
