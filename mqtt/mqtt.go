package mqtt

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	httpServer "iotTester/http"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var Client mqtt.Client

const (
	broker              = "broker.hivemq.com"
	port                = 8883
	mqtt_username       = ""
	mqtt_password       = ""
	initialSubscription = "/api/"
)

// var token = os.Getenv("ENV_TOKEN")

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	httpServer.ReceiveDataFromMqtt(msg.Topic(), msg.Payload())
}

func sub(client mqtt.Client) {
	topic := initialSubscription + "#"
	token := client.Subscribe(topic, 2, nil)
	token.Wait()
	log.Println("Mqtt listening...")
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")

	sub(client)
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v \n", err)
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
	ca, err := GetCertificatesPEM(broker + ":8883")
	if err != nil {
		log.Fatalln(err.Error())
	}
	certpool.AppendCertsFromPEM([]byte(ca))
	return &tls.Config{
		RootCAs: certpool,
	}
}

func Connect() {
	tlsConfig := NewTlsConfig()
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%s:%d", broker, port))
	opts.SetTLSConfig(tlsConfig)
	opts.SetClientID("iotDeviceTester")
	opts.SetUsername(mqtt_username)
	opts.SetPassword(mqtt_password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	Client = mqtt.NewClient(opts)

	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}
