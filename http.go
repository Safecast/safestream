// Copyright 2021 Safecast.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

// Common support for all HTTP topic handlers
package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// httpInboundHandler kicks off inbound messages coming from all sources, then serve HTTP
func httpInboundHandler() {

	// Spin up misc handlers
	http.HandleFunc(httpTopicStream1, httpStreamHandler)
	http.HandleFunc(httpTopicStream2, httpStreamHandler)
	http.HandleFunc(httpTopicGithub, httpGithubHandler)
	http.HandleFunc(httpTopicPing, httpPingHandler)
	http.HandleFunc(httpTopicMain0, httpMainHandler)
	http.HandleFunc(httpTopicMain1, httpMainHandler)
	http.HandleFunc(httpTopicMain2, httpMainHandler)

	// Listen forcing IPV4 so we can do reverse lookup
	l, _ := net.Listen("tcp4", thisServerPort)
	server := &http.Server{}
	server.Serve(l)

}

// httpArgs parses the request URI and returns interesting things
func httpArgs(req *http.Request, topic string) (target string, args map[string]string, err error) {
	args = map[string]string{}

	// Trim the request URI
	target = req.RequestURI[len(topic):]

	// If nothing left, there were no args
	if len(target) == 0 {
		return
	}

	// Make sure that the prefix is "/", else the pattern matcher is matching something we don't want
	if strings.HasPrefix(target, "/") {
		target = strings.TrimPrefix(target, "/")
	}

	// See if there is a query, and if so process it
	str := strings.SplitN(target, "?", 2)
	if len(str) == 1 {
		return
	}
	target = str[0]
	remainder := str[1]

	// See if there is an anchor on the target, and special-case it
	str2 := strings.Split(target, "#")
	if len(str2) > 1 {
		target = str2[0]
		args["anchor"] = str2[1]
	}

	// Now that we know we have args, parse them
	values, err2 := url.ParseQuery(remainder)
	if err2 != nil {
		err = err2
		fmt.Printf("can't parse query: %s\n%s\n", err, str[1])
		return
	}

	// Generate the return arg in the format we expect
	for k, v := range values {
		if len(v) == 1 {
			str := v[0]

			// Safely unquote the value.  This fails if there are NO quotes, so
			// only replace the str value if no error occurs
			s, err2 := strconv.Unquote(str)
			if err2 == nil {
				str = s
			}

			args[k], err = url.PathUnescape(str)
			if err != nil {
				fmt.Printf("can't unescape: %s\n%s\n", err, str)
				return
			}
		}
	}

	// Done
	return

}
