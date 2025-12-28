package main

import (
	"context"
	"iot-gateway-sim/internal/simulation"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Go(func() {
		simulation.SimulateTCPSever(ctx, "localhost:5000", 101)
	})
	wg.Go(func() {
		simulation.SimulateUDPServer(ctx, "localhost:5001", 202)
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
