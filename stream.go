// Copyright 2021 Safecast.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"time"
)

// Handle http requests to the root url
func stream(rsp http.ResponseWriter, req *http.Request, target string, args map[string]string) {

	// Browser clients buffer output before display UNLESS this is the content type
	rsp.Header().Set("Content-Type", "application/json")

	// Begin
	rsp.Write([]byte(time.Now().UTC().Format("2006-01-02T15:04:05Z") + " watching " + target + "\n"))

	// Generate a unique watcher ID
	watcherID := watcherCreate(target, args)

	// Data watching loop
	for {

		// This is an obscure but critical function that flushes partial results
		// back to the client, so that it may display these partial results
		// immediately rather than wait until the end of the transaction.
		f, ok := rsp.(http.Flusher)
		if ok {
			f.Flush()
		} else {
			break
		}

		// Get more data from the watcher, using a timeout computed by trial and
		// error as a reasonable amount of time to catch an error on the Write
		// when the client has gone away.  Longer than that, sometimes the response
		// time in picking up an error becomes quite unpredictable and long.
		data, err := watcherGet(watcherID, 16*time.Second)
		if err != nil {
			break
		}

		// Write either the accumulated notification text, or the idle message,
		// counting on the fact that one or the other will eventually fail when
		// the HTTP client goes away
		if len(data) == 0 {

			s := time.Now().UTC().Format("2006-01-02T15:04:05Z") + " waiting for events\n"
			_, err = rsp.Write([]byte(s))
			if err != nil {
				break
			}

		} else {

			for _, sd := range data {

				class, percent, summary := filterClassify(sd)
				if class == filterClassUnknown {
					continue
				}

				s := fmt.Sprintf("%s %0.02f%% %s\n", class, percent*100, summary)
				_, err = rsp.Write([]byte(s))
				if err != nil {
					break
				}

			}

		}

	}

	// Done
	watcherDelete(watcherID)
	return

}
