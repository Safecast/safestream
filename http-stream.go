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

// Handle http requests to the root url
func streamLaunch(rsp http.ResponseWriter, req *http.Request) {
	streamTemplate.Execute(rsp, "ws://"+req.Host+httpTopicStream1+req.URL.String())
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

		// Get more data from the watcher, using a timeout computed by trial and
		// error as a reasonable amount of time to catch an error on the Write
		// when the client has gone away.  Longer than that, sometimes the response
		// time in picking up an error becomes quite unpredictable and long.
		data, ipinfo, err := watcherGet(watcherID, 16*time.Second)
		if err != nil {
			break
		}

		// Write either the accumulated notification text, or the idle message,
		// counting on the fact that one or the other will eventually fail when
		// the HTTP client goes away
		if len(data) == 0 {

			s := time.Now().UTC().Format("2006-01-02T15:04:05Z") + " waiting for events\n"
			err = c.WriteMessage(websocket.TextMessage, []byte(s))
			if err != nil {
				break
			}

		} else {

			err = nil
			for _, sd := range data {
				events := filterClassify(sd, ipinfo)
				for _, e := range events {
					s := fmt.Sprintf("%s %0.02f%% %s %.0fkm %s %s\n",
						e.class, e.percent*100, e.summary,
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

// The HTML template
var streamTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>

window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
	var playButton = document.getElementById("play");
	var p1 = document.getElementById("p1")

    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
		d.scrollIntoView()
    };

	(function () {
	  'use strict';

	  const URL = 'https://s3-us-west-2.amazonaws.com/s.cdpn.io/123941/Yodel_Sound_Effect.mp3';

	  const context = new AudioContext();

	  let yodelBuffer;

	  window.fetch(URL)
	    .then(response => response.arrayBuffer())
	    .then(arrayBuffer => context.decodeAudioData(arrayBuffer))
	    .then(audioBuffer => {
	      playButton.disabled = false;
	      yodelBuffer = audioBuffer;
	    });

		playButton.onclick = function(evt) {
			playButton.style.display = "none";
			play(yodelBuffer)
	        return false;
		};

		p1.onclick = function(evt) {
			play(yodelBuffer)
	        return false;
		};

	  function play(audioBuffer) {
	    const source = context.createBufferSource();
	    source.buffer = audioBuffer;
	    source.connect(context.destination);
	    source.start();
	  }
	}());

    var enterfn = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("now listening for Safecast events");
        }
        ws.onclose = function(evt) {
            print("connection closed");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print(evt.data);
			p1.click()
        }
        ws.onerror = function(evt) {
            print("error: " + evt.data);
        }
        return false;
    };

	var exitfn = function(evt) {
        if (!ws) {
            return false;
        }
        ws.send(input.value);
        return false;
    };

	var e = document.getElementById("open");
    if (e) { e.onclick = enterfn };
	var e = document.getElementById("send")
    if (e) { e.onclick = exitfn };
	var e = document.getElementById("close")
    if (e) {
		e.onclick = function(evt) {
	        if (!ws) {
	            return false;
	        }
	        ws.close();
	        return false;
		};
    };

	enterfn();

});
</script>
</head>
<body>
<!--
<p>Click "Open" to create a connection to the server,
"Send" to send a message to the server and "Close" to close the connection.
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
-->
<div id="output"></div>
<div id="footer">
<p><p>
<form>
<!--
<p><input id="input" type="text" value="">
<button id="send">request</button>
<p>
-->
<button id="play" disabled>Listen to Event Stream</button>
<p hidden>
<button id="p1">
</button>
</form>
</div>
</body>
</html>
`))
