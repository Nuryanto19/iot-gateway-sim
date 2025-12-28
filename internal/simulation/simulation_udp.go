package simulation

import (
	"context"
	"iot-gateway-sim/pkg/model"
	"log"
	"math/rand/v2"
	"net"
	"time"
)

func SimulateUDPServer(ctx context.Context, addr string, id uint32) {

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var conn net.Conn

	for {
		select {
		case <-ctx.Done():
			log.Printf("[UDP] Stopping micro controller %d.", id)
			if conn != nil {
				conn.Close()
			}
			return
		case <-ticker.C:
			if conn == nil {
				newConn, err := net.Dial("udp", addr)
				if err != nil {
					log.Printf("[UDP] Connection failed micro controller %d: %v", id, err)
					conn = nil
					continue
				}
				conn = newConn
				log.Printf("[UDP] Micro contoller %d connected.", id)
			}

			data := model.SensorData{
				DeviceID: id,
				Value:    25.0 + rand.Float32()*10.0, // simulate temperature
			}

			_, err := conn.Write(data.Pack())
			if err != nil {
				log.Printf("[UDP] Error sending data from micro controller %d: %v", id, err)
				conn.Close()
				conn = nil
				continue
			}

			log.Printf("[UDP] Micro controller %d sent: %.2f", id, data.Value)

		}

	}
}
