// Copyright 2021 Safecast.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"time"
)

func main() {

	// Start up the MQTT listener
	mqttInit()

	// Fire-up the HTTP handler
	go httpInboundHandler()

	// Sleep indefinitely
	for {
		time.Sleep(60 * time.Second)
	}

}
