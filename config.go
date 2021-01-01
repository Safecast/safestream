// Copyright 2017 Inca Roads LLC.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

// Global configuration Parameters
package main

// Safecast public MQTT service params
const mqttServer = "tcp://m14.cloudmqtt.com:17506"
const mqttTopic = "device/#"
const mqttUsername = "zxskatjq"
const mqttPassword = "C0IVCvumIOSg"

// httpHTTPPort (here for golint)
const httpPort string = ":80"

// httpPortAlternate (here for golint)
const httpPortAlternate string = ":8888"

// httpTopicMain0 (here for golint)
const httpTopicMain0 string = "/"

// httpTopicMain1 (here for golint)
const httpTopicMain1 string = "/index.html"

// httpTopicMain2 (here for golint)
const httpTopicMain2 string = "/index.htm"

// httpTopicPing (here for golint)
const httpTopicPing string = "/ping"
