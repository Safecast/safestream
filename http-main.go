// Copyright 2021 Safecast.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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

	// Process the stream
	if method == "GET" && target == "" {
		streamLaunch(rsp, req)
		return
	}

	// If the resource has a slash in it, it's in a collection
	if strings.Contains(target, "/") {
		s := strings.Split(target, "/")
		collection := s[0]
		filename := s[1]
		contents, err := uploadRead(collection, filename)
		if err != nil {
			fmt.Printf("uploadRead: %s\n", err)
		} else {
			rsp.Write(contents)
		}
		return
	}

	// Attempt to load the resource
	contents, _ := resourceRead(target)
	rsp.Write(contents)

	// Done
	return

}

// Get a resource path
func resourcePath(filename string) (path string) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	path = dir + filePathResources + filename
	return
}

// Get a resource URL
func resourceURL(filename string) (url string) {
	url = "http://" + thisServerAddressIPv4 + thisServerPort + "/" + filename
	return
}

// Get a resource's contents
func resourceRead(filename string) (contents []byte, err error) {
	contents, err = ioutil.ReadFile(resourcePath(filename))
	return
}
