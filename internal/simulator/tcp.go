package simulator

import (
	"context"
	"iot-gateway-sim/internal/model"
	"log"
	"math/rand/v2"
	"net"
	"time"
)

func SimulateTCPSever(ctx context.Context, addr string, id uint32) {

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	var conn net.Conn

	for {
		select {
		case <-ctx.Done():
			if conn != nil {
				log.Printf("[TCP] Stopping micro controller %d.", id)
				conn.Close()
			}
			return
		case <-ticker.C:
			if conn == nil {
				log.Printf("[TCP] Trying to connect micro controller %d to %s", id, addr)

				newConn, err := net.Dial("tcp", addr)
				if err != nil {
					log.Printf("[TCP] Connection failded to micro controller %d: %v", id, err)
					conn = nil
					continue
				}
				conn = newConn
				log.Printf("[TCP] Micro contoller %d connected.", id)
			}

			data := model.SensorData{
				DeviceID: id,
				Value:    220.0 + rand.Float32()*5.0, // simulate electric voltage
			}

			_, err := conn.Write(data.Pack())
			if err != nil {
				log.Printf("[TCP] Error sending data from micro controller %d: %v", id, err)
				conn.Close()
				conn = nil
				continue
			}

			log.Printf("[TCP] Micro contoller %d sent: %.2f", id, data.Value)

		}

	}
}
