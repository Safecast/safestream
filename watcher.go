// Copyright 2021 Safecast.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/Safecast/ttdefs"
	"github.com/google/uuid"
)

// The active watcher data structure
type activeWatcher struct {
	watcherID string
	ipinfo    IPInfoData
	target    string
	args      map[string]string
	event     *Event
	buf       []ttdefs.SafecastData
}

var watchers = []activeWatcher{}
var watcherLock sync.RWMutex

// Create a new watcher
func watcherCreate(requestorIP string, target string, args map[string]string) (watcherID string) {

	// Create the watcher
	watcherID = uuid.New().String()
	watcher := activeWatcher{}
	watcher.watcherID = watcherID
	watcher.target = target
	watcher.args = args
	watcher.event = eventNew()

	// For localhost debugging
	if requestorIP == "127.0.0.1" {
		requestorIP = "65.96.197.34"
	}

	// Look up info about the requestor
	url := "http://ip-api.com/json/" + requestorIP
	rsp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s: error getting location info: %s\n", requestorIP, err)
	} else {
		defer rsp.Body.Close()
		buf, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			fmt.Printf("%s: error reading location info: %s\n", requestorIP, err)
		} else {
			err := json.Unmarshal(buf, &watcher.ipinfo)
			if err != nil {
				fmt.Printf("%s: error unmarshaling: %s\n", requestorIP, err)
			}
		}
	}

	// Add to queue of watchers
	watcherLock.Lock()
	watchers = append(watchers, watcher)
	fmt.Printf("watchers: added (now %d)\n", len(watchers))
	watcherLock.Unlock()

	return
}

// Delete a watcher
func watcherDelete(watcherID string) {

	watcherLock.Lock()
	numWatchers := len(watchers)
	for i, watcher := range watchers {
		if watcher.watcherID == watcherID {
			if i == numWatchers-1 {
				watchers = watchers[0:i]
			} else {
				watchers = append(watchers[0:i], watchers[i+1:]...)
			}
			fmt.Printf("watchers: removed (now %d)\n", len(watchers))
			break
		}
	}
	watcherLock.Unlock()

	return

}

// Get data from a watcher
func watcherGet(watcherID string, timeout time.Duration) (data []ttdefs.SafecastData, ipinfo IPInfoData, err error) {
	var watcher activeWatcher

	// Find the watcher
	watcherLock.Lock()
	for _, watcher = range watchers {
		if watcher.watcherID == watcherID {
			break
		}
	}
	watcherLock.Unlock()

	// If not found, we're done
	if watcher.watcherID != watcherID {
		err = fmt.Errorf("watcher not found")
		return
	}

	// Wait with timeout
	watcher.event.Wait(timeout)

	// Get the buffer
	watcherLock.Lock()
	for i := range watchers {
		if watchers[i].watcherID == watcherID {
			data = watchers[i].buf
			ipinfo = watchers[i].ipinfo
			watchers[i].buf = []ttdefs.SafecastData{}
			break
		}
	}
	watcherLock.Unlock()

	return

}

// Append data from a watcher
func watcherPut(data ttdefs.SafecastData) {

	// Scan all watchers
	watcherLock.Lock()
	for i := range watchers {
		if filterMatches(data, watchers[i].target, watchers[i].args) {
			watchers[i].buf = append(watchers[i].buf, data)
			watchers[i].event.Signal()
		}
	}
	watcherLock.Unlock()

	return

}
