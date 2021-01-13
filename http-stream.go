// Copyright 2021 Safecast.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
)

// WS support
var upgrader = websocket.Upgrader{} // use default options

// The main stream template
var streamTemplate *template.Template

// Initialize for stream processing
func streamInit() (err error) {
	var contents []byte
	contents, err = resourceRead("main.html")
	if err != nil {
		return
	}
	streamTemplate = template.Must(template.New("").Parse(string(contents)))
	return
}

// Handle http requests to the root url
func streamLaunch(rsp http.ResponseWriter, req *http.Request) {
	isTLS := req.Header.Get("X-Forwarded-Proto") == "https"
	scheme := "ws:"
	if isTLS {
		scheme = "wss:"
	}
	url := scheme + "//" + req.Host + httpTopicStream1 + req.URL.String()
	streamTemplate.Execute(rsp, url)
}

// Handle ws template
func httpStreamHandler(rsp http.ResponseWriter, req *http.Request) {

	// Get the args
	target, args, err := httpArgs(req, "")
	if err != nil {
		return
	}
	if strings.HasSuffix(target, "/") {
		target = strings.TrimSuffix(target, "/")
	}

	// Remove stream prefix from target, in all variations
	target = "/" + target
	target = strings.TrimPrefix(target, httpTopicStream2)
	target = strings.TrimPrefix(target, httpTopicStream1)
	target = strings.TrimPrefix(target, "/")

	// Upgrade the endpoint to use websockets
	c, err := upgrader.Upgrade(rsp, req, nil)
	if err != nil {
		fmt.Printf("upgrade: %s\n", err)
		return
	}
	defer c.Close()

	// Launch the reader
	inboundExited := false
	go processInbound(c, rsp, req, &inboundExited)

	// Generate a unique watcher ID
	requestorIPV4, _ := getRequestorIPv4(req)
	watcherID := watcherCreate(requestorIPV4, target, args)

	// Data watching loop
	for !inboundExited {

		// Get more data from the watcher
		data, ipinfo, err := watcherGet(watcherID, 10*time.Second)
		if err != nil {
			break
		}

		// Write either the accumulated notification text, or the idle message,
		// counting on the fact that one or the other will eventually fail when
		// the HTTP client goes away
		if len(data) == 0 {

			s := "waiting for events\n"
			err = c.WriteMessage(websocket.TextMessage, []byte(s))
			if err != nil {
				break
			}

		} else {

			err = nil
			for _, sd := range data {
				events := filterClassify(sd, ipinfo)
				for _, e := range events {
					s := fmt.Sprintf("%s%02d %0.02f%% %s %.0fkm %s %s\n",
						e.class, int(e.percent*10),
						e.percent*100, e.summary,
						e.distance/1000, e.city, e.country)
					err = c.WriteMessage(websocket.TextMessage, []byte(s))
					if err != nil {
						break
					}
					time.Sleep(250 * time.Millisecond)
				}
			}
			if err != nil {
				break
			}

		}

	}

	// Done
	watcherDelete(watcherID)
	return

}

// Function that proceses messages coming in from the JS client
func processInbound(c *websocket.Conn, w http.ResponseWriter, r *http.Request, exited *bool) {
	fmt.Printf("inbound: enter\n")
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
	fmt.Printf("inbound: exit\n")
	*exited = true
}
