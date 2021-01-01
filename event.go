// Copyright 2021 safecast.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"time"
)

// Event is our timeout-enabled event waiter abstraction
type Event struct {
	eq chan struct{}
}

// EventNew allocates a new semaphore, initially unsignalled
func eventNew() *Event {
	evt := Event{}
	evt.eq = make(chan struct{}, 1)
	return &evt
}

// IsSignalled will return true if the event has been signalled
func (evt *Event) IsSignalled() bool {
	return len(evt.eq) != 0
}

// Wait waits for the event until the specified timeout, and
// returns true if acquired, and false if timeout
func (evt *Event) Wait(timeout time.Duration) bool {
	select {
	case <-evt.eq:
		return true
	case <-time.After(timeout):
		return false
	}
}

// Signal ensures that the waiter becomes unblocked
func (evt *Event) Signal() bool {
	select {
	case evt.eq <- struct{}{}:
		return true
	default:
		return false
	}
}
