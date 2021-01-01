package main

import (
	"time"
)

func main() {

	// Start up the MQTT listener
	mqttInit()

	// Sleep indefinitely
	for {
		time.Sleep(60 * time.Second)
	}

}
