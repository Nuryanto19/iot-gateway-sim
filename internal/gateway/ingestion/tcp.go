package ingestion

import (
	"context"
	"io"
	"iot-gateway-sim/internal/model"
	"log"
	"net"
)

func StartTCPServer(ctx context.Context, address string, dataChan chan<- model.SensorData) {

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("[TCP-Listener] Failed to listen from %s: %v", address, err)
	}

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	log.Printf("[TCP-Listener] Listening on %s", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				log.Println("[TCP-Listener] Listener have closed, shutdown process.")
				return
			default:
				log.Printf("[TCP-Listener] Failed to accept connection from %s : %v", address, err)
				continue
			}
		}

		go connHandler(ctx, conn, dataChan)
	}
}

func connHandler(ctx context.Context, conn net.Conn, dataChan chan<- model.SensorData) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	log.Printf("[TCP-Listener] New micro controller connected: %s", remoteAddr)

	buf := make([]byte, 8)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, err := io.ReadFull(conn, buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("[TCP-Listener] Error reading data from %s: %v", remoteAddr, err)
				}
				return
			}

			data := model.Unpack(buf)
			data.Protocol = "tcp"

			select {
			case dataChan <- data:
			default:
				log.Printf("[TCP-Listener] Dropping data from %s due to backpressure", remoteAddr)
			}
		}
	}
}
