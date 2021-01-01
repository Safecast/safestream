// Copyright 2021 Safecast.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/safecast/ttdata"
)

// The active watcher data structure
type activeWatcher struct {
	watcherID string
	target    string
	args      map[string]string
	event     *Event
	buf       []ttdata.SafecastData
}

var watchers = []activeWatcher{}
var watcherLock sync.RWMutex

// Create a new watcher
func watcherCreate(target string, args map[string]string) (watcherID string) {

	watcherID = uuid.New().String()

	watcher := activeWatcher{}
	watcher.watcherID = watcherID
	watcher.target = target
	watcher.args = args
	watcher.event = eventNew()

	watcherLock.Lock()
	watchers = append(watchers, watcher)
	fmt.Printf("watchers: %s added (now %d)\n", watcher.target, len(watchers))
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
			fmt.Printf("watchers: %s removed (now %d)\n", watcher.target, len(watchers))
			break
		}
	}
	watcherLock.Unlock()

	return

}

// Get data from a watcher
func watcherGet(watcherID string, timeout time.Duration) (data []ttdata.SafecastData, err error) {
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
			watchers[i].buf = []ttdata.SafecastData{}
			break
		}
	}
	watcherLock.Unlock()

	return

}

// Append data from a watcher
func watcherPut(data ttdata.SafecastData) {

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
