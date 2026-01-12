package main

import (
	"context"
	"iot-gateway-sim/internal/simulator"
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

	tcpTargetAddr := getEnv("SIM_TARGET_TCP_ADDR", "localhost:5000")
	udpTargetAddr := getEnv("SIM_TARGET_UDP_ADDR", "localhost:5001")

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Go(func() {
		simulator.SimulateTCPSever(ctx, tcpTargetAddr, 101)
	})
	wg.Go(func() {
		simulator.SimulateUDPServer(ctx, udpTargetAddr, 202)
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	wg.Wait()
	log.Println("All simutor have stoped. Exiting...")
}
