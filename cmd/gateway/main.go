package main

import (
	"context"
	"iot-gateway-sim/internal/gateway/ingestion"
	"iot-gateway-sim/internal/gateway/processing"
	"iot-gateway-sim/internal/gateway/transport"
	"iot-gateway-sim/internal/model"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {

	brokerAddr := getEnv("MQTT_BROKER_ADDR", "ssl://localhost:8883")
	clientID := getEnv("MQTT_CLIENT_ID", "iot-gateway-1")
	tcpAddr := getEnv("GATEWAY_TCP_ADDR", ":5000")
	udpAddr := getEnv("GATEWAY_UDP_ADDR", ":5001")

	dataCh := make(chan model.SensorData, 100)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}

	mqttClient, err := transport.NewMQTTClient(
		brokerAddr,
		clientID,
		"gateway-certs/ca.crt",
		"gateway-certs/client.crt",
		"gateway-certs/client.key",
	)
	if err != nil {
		log.Fatalf("FATAL: MQTT auth failed: %v", err)
	}
	defer mqttClient.Disconnect()

	startPipeline(ctx, wg, dataCh, mqttClient, tcpAddr, udpAddr)

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
	tcpAddr, udpAddr string,
) {
	wg.Add(3)

	wg.Go(func() {
		ingestion.StartTCPServer(ctx, tcpAddr, dataCh)
	})

	wg.Go(func() {
		ingestion.StartUDPServer(ctx, udpAddr, dataCh)
	})

	wg.Go(func() {
		processing.RunBufferPipeline(ctx, dataCh, mqttClient)
	})

}
