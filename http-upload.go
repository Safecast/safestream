// Copyright 2021 Safecast.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
)

// Handle http requests to the root url
func httpUploadHandler(rsp http.ResponseWriter, req *http.Request) {

	// Get the args
	_, args, err := httpArgs(req, "")
	if err != nil {
		return
	}
	filename := args["file"]
	if filename == "" {
		io.WriteString(rsp, "file not specified")
		return
	}
	collection := args["collection"]
	if collection == "" {
		io.WriteString(rsp, "collection not specified")
		return
	}

	// Read the file
	var contents []byte
	contents, err = ioutil.ReadAll(req.Body)
	if err != nil {
		io.WriteString(rsp, fmt.Sprintf("can't read file: %s", err))
		return
	}
	if len(contents) == 0 {
		io.WriteString(rsp, "file not uploaded, or file is empty")
		return
	}

	// Make sure that the collection directory is created
	usr, err := user.Current()
	if err != nil {
		io.WriteString(rsp, fmt.Sprintf("can't get home directory: %s\n", err))
		return
	}
	directory := usr.HomeDir + filePathCollections + collection
	os.MkdirAll(directory, 0777)

	// Write the file
	err = ioutil.WriteFile(directory+"/"+filename, contents, 0644)
	if err != nil {
		io.WriteString(rsp, fmt.Sprintf("can't write file: %s\n", err))
		return
	}

	// Done
	return

}
