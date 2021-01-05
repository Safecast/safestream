// Copyright Safecast.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

// Github webhook that enables server auto-restart on commit
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Github webhook
func httpGithubHandler(rw http.ResponseWriter, req *http.Request) {

	// Unpack the request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Github webhook: error reading body: %s", err)
		return
	}
	var p PushPayload
	err = json.Unmarshal(body, &p)
	if err != nil {
		fmt.Printf("Github webhook: error unmarshaling body: %s", err)
		return
	}

	// Handle 'git commit -mm' and 'git commit -amm', used in dev builds, in a more aesthetically pleasing manner.
	if p.HeadCommit.Commit.Message == "m" {
		fmt.Printf(fmt.Sprintf("*** RESTARTING because %s pushed %s's commit to GitHub\n", p.Pusher.Name, p.HeadCommit.Commit.Committer.Name))
	} else {
		fmt.Printf(fmt.Sprintf("*** RESTARTING because %s pushed %s's commit to GitHub: %s\n",
			p.Pusher.Name, p.HeadCommit.Commit.Committer.Name, p.HeadCommit.Commit.Message))
	}

	// Exit
	os.Exit(0)

}
