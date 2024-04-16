package mqtt

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	httpServer "iotTester/http"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var Client mqtt.Client

func Connect() {
	tlsConfig := NewTlsConfig()
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%s:%s", os.Getenv("BROKER"), os.Getenv("MQTT_PORT")))

	if os.Getenv("USE_MQTT") == "TRUE" {
		opts.SetTLSConfig(tlsConfig)
	}

	opts.SetClientID(os.Getenv("CLIENT_ID"))
	opts.SetUsername(os.Getenv("USERNAME"))
	opts.SetPassword(os.Getenv("PASSWORD"))
	opts.SetDefaultPublishHandler(
		func(client mqtt.Client, msg mqtt.Message) {
			httpServer.ReceiveDataFromMqtt(msg.Topic(), msg.Payload())
		},
	)

	opts.OnConnect = func(client mqtt.Client) {
		fmt.Println("Connected")

		sub(client)
	}

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Printf("Connect lost: %v \n", err)
	}

	Client = mqtt.NewClient(opts)

	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func sub(client mqtt.Client) {
	topic := os.Getenv("TOPIC") + "#"
	token := client.Subscribe(topic, 2, nil)
	token.Wait()
	log.Println("Mqtt listening...")
}

func GetCertificatesPEM(address string) (string, error) {
	conn, err := tls.Dial("tcp", address, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return "", err
	}
	defer conn.Close()
	var b bytes.Buffer
	for _, cert := range conn.ConnectionState().PeerCertificates {
		err := pem.Encode(&b, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})
		if err != nil {
			return "", err
		}
	}
	return b.String(), nil
}

func NewTlsConfig() *tls.Config {
	certpool := x509.NewCertPool()
	ca, err := GetCertificatesPEM(os.Getenv("BROKER") + ":" + os.Getenv("MQTT_PORT"))
	if err != nil {
		log.Fatalln("Error with Certificate", err.Error())
	}
	certpool.AppendCertsFromPEM([]byte(ca))
	return &tls.Config{
		RootCAs: certpool,
	}
}
