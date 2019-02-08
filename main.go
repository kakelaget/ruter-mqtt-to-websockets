package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	uri := os.Getenv("PS_BROKER_URI")
	user := os.Getenv("PS_MQTT_USERNAME")
	pass := os.Getenv("PS_MQTT_PASSWORD")

	MQTT.DEBUG = log.New(os.Stdout, "", 0)
	MQTT.ERROR = log.New(os.Stdout, "", 0)

	tlsconfig := NewTLSConfig()

	opts := MQTT.NewClientOptions()
	opts.AddBroker(uri)
	opts.SetUsername(user)
	opts.SetPassword(pass)
	opts.SetClientID("miljohack")
	opts.SetTLSConfig(tlsconfig)
	opts.SetDefaultPublishHandler(f)
	opts.SetKeepAlive(30)

	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	client.Subscribe("#", 0, nil)

	select {}
}

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s - ", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func NewTLSConfig() *tls.Config {
	// Import trusted certificates
	certpool := x509.NewCertPool()

	filepath := os.Getenv("PS_MQTT_CERT")
	pemCerts, err := ioutil.ReadFile(filepath)
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: false,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{},
	}
}

