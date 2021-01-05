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
	var header = document.getElementById("header");

    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
		d.scrollIntoView()
    };

	(function () {
	  'use strict';

	  const URL = "yodel.mp3"

	  const context = new AudioContext();

	  let yodelBuffer;
	  let resBuffer = new Map();

	  window.fetch(URL)
	    .then(response => response.arrayBuffer())
	    .then(arrayBuffer => context.decodeAudioData(arrayBuffer))
	    .then(audioBuffer => {
	      playButton.disabled = false;
	      yodelBuffer = audioBuffer;
	    });

	  function load(resURL) {
	  window.fetch(resURL)
	    .then(response => response.arrayBuffer())
	    .then(arrayBuffer => context.decodeAudioData(arrayBuffer))
	    .then(audioBuffer => {
		  resBuffer.set(resURL, audioBuffer)
	    });
      }

      load("air00.mp3")
      load("air01.mp3")
      load("air02.mp3")
      load("air03.mp3")
      load("air04.mp3")
      load("air05.mp3")
      load("air06.mp3")
      load("air07.mp3")
      load("air08.mp3")
      load("air09.mp3")
      load("air10.mp3")
      load("rad00.mp3")
      load("rad01.mp3")
      load("rad02.mp3")
      load("rad03.mp3")
      load("rad04.mp3")
      load("rad05.mp3")
      load("rad06.mp3")
      load("rad07.mp3")
      load("rad08.mp3")
      load("rad09.mp3")
      load("rad10.mp3")

		playButton.onclick = function(evt) {
			playButton.style.display = "none";
			play(yodelBuffer);
	        return false;
		};

		air00.onclick = function(evt) {
            play(resBuffer.get("air00.mp3"));
	        return false;
		};
		air01.onclick = function(evt) {
            play(resBuffer.get("air01.mp3"));
	        return false;
		};
		air02.onclick = function(evt) {
            play(resBuffer.get("air02.mp3"));
	        return false;
		};
		air03.onclick = function(evt) {
            play(resBuffer.get("air03.mp3"));
	        return false;
		};
		air04.onclick = function(evt) {
            play(resBuffer.get("air04.mp3"));
	        return false;
		};
		air05.onclick = function(evt) {
            play(resBuffer.get("air05.mp3"));
	        return false;
		};
		air06.onclick = function(evt) {
            play(resBuffer.get("air06.mp3"));
	        return false;
		};
		air07.onclick = function(evt) {
            play(resBuffer.get("air07.mp3"));
	        return false;
		};
		air08.onclick = function(evt) {
            play(resBuffer.get("air08.mp3"));
	        return false;
		};
		air09.onclick = function(evt) {
            play(resBuffer.get("air09.mp3"));
	        return false;
		};
		air10.onclick = function(evt) {
            play(resBuffer.get("air10.mp3"));
	        return false;
		};
		rad00.onclick = function(evt) {
            play(resBuffer.get("rad00.mp3"));
	        return false;
		};
		rad01.onclick = function(evt) {
            play(resBuffer.get("rad01.mp3"));
	        return false;
		};
		rad02.onclick = function(evt) {
            play(resBuffer.get("rad02.mp3"));
	        return false;
		};
		rad03.onclick = function(evt) {
            play(resBuffer.get("rad03.mp3"));
	        return false;
		};
		rad04.onclick = function(evt) {
            play(resBuffer.get("rad04.mp3"));
	        return false;
		};
		rad05.onclick = function(evt) {
            play(resBuffer.get("rad05.mp3"));
	        return false;
		};
		rad06.onclick = function(evt) {
            play(resBuffer.get("rad06.mp3"));
	        return false;
		};
		rad07.onclick = function(evt) {
            play(resBuffer.get("rad07.mp3"));
	        return false;
		};
		rad08.onclick = function(evt) {
            play(resBuffer.get("rad08.mp3"));
	        return false;
		};
		rad09.onclick = function(evt) {
            play(resBuffer.get("rad09.mp3"));
	        return false;
		};
		rad10.onclick = function(evt) {
            play(resBuffer.get("rad10.mp3"));
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
			header.style.display = "inline";
        }
        ws.onclose = function(evt) {
            print("connection closed");
            ws = null;
        }
        ws.onmessage = function(evt) {
			header.style.display = "none";
			var line = evt.data
            print(line);
			var firstWord = line.substr(0, line.indexOf(" "));
			document.getElementById(firstWord).click()
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
<div id="header" hidden>One moment please, as we wait for the next Safecast event...</div>
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
<button id="air00">
<button id="air01">
<button id="air02">
<button id="air03">
<button id="air04">
<button id="air05">
<button id="air06">
<button id="air07">
<button id="air08">
<button id="air09">
<button id="air10">
<button id="rad00">
<button id="rad01">
<button id="rad02">
<button id="rad03">
<button id="rad04">
<button id="rad05">
<button id="rad06">
<button id="rad07">
<button id="rad08">
<button id="rad09">
<button id="rad10">
</button>
</form>
</div>
</body>
</html>
`))
