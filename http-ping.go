// Copyright 2021 Safecast.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Handle http requests for health checks
func httpPingHandler(rsp http.ResponseWriter, req *http.Request) {
	io.WriteString(rsp, fmt.Sprintf("%s\n", time.Now().UTC().Format("2006-01-02T15:04:05Z")))
}
