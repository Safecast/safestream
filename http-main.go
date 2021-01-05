// Copyright 2021 Safecast.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Handle http requests to the root url
func httpMainHandler(rsp http.ResponseWriter, req *http.Request) {

	// Process the request URI, looking for things that will indicate "dev"
	method := req.Method
	if method == "" {
		method = "GET"
	}

	// Get the target
	target, _, err := httpArgs(req, "")
	if err != nil {
		return
	}
	if strings.HasSuffix(target, "/") {
		target = strings.TrimSuffix(target, "/")
	}

	// Exit if just the favicon
	if target == "favicon.ico" {
		return
	}

	// Process
	if method == "GET" {
		streamLaunch(rsp, req)
		return
	}

	// Done
	io.WriteString(rsp, fmt.Sprintf("unrecognized function\n"))
	return

}
