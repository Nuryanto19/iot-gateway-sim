package ingestion

import (
	"context"
	"iot-gateway-sim/pkg/model"
	"log"
	"net"
)

func StartUDPServer(ctx context.Context, address string, dataChan chan<- model.SensorData) {

	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		log.Fatalf("[UDP-Listener] Listen Error from %s: %v", address, err)
	}

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	buf := make([]byte, 8)
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				log.Printf("[UDP-Listener] ReadFrom error: %v", err)
				continue
			}
		}

		if n != 8 {
			log.Printf("[UDP-Listener] invalid packet size from %v, expect %d byte", addr, n)
			continue
		}

		sensorData := model.Unpack(buf)
		sensorData.Protocol = "udp"

		select {
		case dataChan <- sensorData:
		default:
			log.Printf("[UDP-Listener] Dropping data from %s due to backpressure", addr)

		}

	}
}
