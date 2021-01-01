package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/safecast/ttdata"
)

const mqttServer = "tcp://m14.cloudmqtt.com:17506"
const mqttTopic = "device/#"
const mqttQOS = 0
const mqttUsername = "zxskatjq"
const mqttPassword = "C0IVCvumIOSg"

var eventQ chan ttdata.SafecastData

// mqttInit initializes the listener
func mqttInit() {

	// Create the queue
	eventQ = make(chan ttdata.SafecastData, 10000)
	go mqttEventQHandler(eventQ)

	// Get connect parameters
	hostname, _ := os.Hostname()
	clientID := hostname + strconv.Itoa(time.Now().Second()) + "x"
	connOpts := MQTT.NewClientOptions().AddBroker(mqttServer).SetClientID(clientID).SetCleanSession(true)
	connOpts.SetUsername(mqttUsername)
	connOpts.SetPassword(mqttPassword)
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)
	connOpts.OnConnect = mqttConnectHandler

	// Connect
	client := MQTT.NewClient(connOpts)
	for {
		token := client.Connect()
		token.Wait()
		if token.Error() != nil {
			fmt.Printf("MQTT connect error: %s\n", token.Error())
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

}

// mqttMessageReceived handles received messages
func mqttMessageReceived(client MQTT.Client, message MQTT.Message) {

	var e ttdata.SafecastData
	err := json.Unmarshal(message.Payload(), &e)
	if err != nil {
		fmt.Printf("mqtt message error: %s\n", err)
		return
	}
	eventQ <- e

}

// mqttConnectHandler handles connections
func mqttConnectHandler(c MQTT.Client) {
	for {
		token := c.Subscribe(mqttTopic, byte(mqttQOS), mqttMessageReceived)
		token.Wait()
		if token.Error() != nil {
			fmt.Printf("MQTT subscribe error: %s\n", token.Error())
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
}

// Event queue handler
func mqttEventQHandler(ch <-chan ttdata.SafecastData) {

	for {

		// Pull the event from the channel
		e := <-ch

		// Process the event
		fmt.Printf("Processing: %+v\n", e)

	}

}
