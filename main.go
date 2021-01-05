// Copyright 2021 Safecast.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {

	// Get our external IP address
	rsp, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		fmt.Printf("Can't get our own IP address: %v\n", err)
		os.Exit(0)
	}
	defer rsp.Body.Close()
	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Printf("Error fetching IP addr: %v\n", err)
		os.Exit(0)
	}
	thisServerAddressIPv4 = string(bytes.TrimSpace(buf))

	// Get HTTP port
	if len(os.Args) > 1 {
		thisServerPort = ":" + os.Args[1]
	}

	// Start up the MQTT listener
	mqttInit()

	// Fire-up the HTTP handler
	go httpInboundHandler()

	// Sleep indefinitely
	for {
		time.Sleep(60 * time.Second)
	}

}
