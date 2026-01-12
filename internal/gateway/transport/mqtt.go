package transport

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client mqtt.Client
}

func newTLSConfig(caFile, certFile, keyFile string) *tls.Config {

	caCert, err := os.ReadFile(caFile)
	if err != nil {
		log.Fatalf("Failed to read CA file: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	clientCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load key pair client: %v", err)
	}

	return &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{clientCert},
	}
}

func NewMQTTClient(broker, clientID, caFile, certFile, keyFile string) (*MQTTClient, error) {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)

	tlsConfig := newTLSConfig(caFile, certFile, keyFile)
	opts.SetTLSConfig(tlsConfig)

	opts.SetConnectRetry(true)
	opts.SetAutoReconnect(true)
	opts.SetCleanSession(true)

	opts.OnConnect = func(c mqtt.Client) {
		log.Printf("MQTT Client connected to %s", broker)
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("MQTT connection lost: %v", err)
	}

	client := mqtt.NewClient(opts)

	token := client.Connect()
	if ok := token.WaitTimeout(5 * time.Second); !ok {
		return nil, fmt.Errorf("mqtt client timout")
	}

	if token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("MQTT connect failed: %v", token.Error())
	}

	return &MQTTClient{client: client}, nil
}

func (m *MQTTClient) Publish(topic string, payload []byte) error {

	token := m.client.Publish(topic, 1, false, payload)

	if !token.WaitTimeout(5 * time.Second) {
		log.Printf("MQTT publish to topic %s timed out", topic)
		return fmt.Errorf("mqtt publish timed out")
	}

	return nil
}

func (m *MQTTClient) Disconnect() {
	log.Println("Disconnecting from MQTT broker...")
	m.client.Disconnect(250)
	log.Println("MQTT client disconnected.")
}
