// Copyright 2021 Safecast.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Safecast/ttdefs"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var eventQ chan ttdefs.SafecastData

// mqttInit initializes the listener
func mqttInit() {

	// Create the queue
	eventQ = make(chan ttdefs.SafecastData, 10000)
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

	var e ttdefs.SafecastData
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
		mqttQOS := 0
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
func mqttEventQHandler(ch <-chan ttdefs.SafecastData) {

	for {

		// Pull the event from the channel
		data := <-ch

		// Process the event
		watcherPut(data)

	}

}
