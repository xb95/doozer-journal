// Copyright (c) 2012, SoundCloud Ltd.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/soundcloud/doozer-journal

package main

import (
	"github.com/soundcloud/doozer-journal/journal"
	"os"
)

var cmdSnapshot = &Command{
	Name:      "snapshot",
	Desc:      "makes a snapshot and exits",
	UsageLine: "snapshot",
}

func init() {
	cmdSnapshot.Run = runSnapshot
}

func runSnapshot(cmd *Command, args []string) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_EXCL|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		if os.IsExist(err) {
			exitWithError("Journal already exists at %s\n", file)
		} else {
			exitWithError("Unable to open journal: %s\n", err)
		}
	}

	j := journal.New(f)
	err = snapshot(cmd.Conn, cmd.Rev, j)
	if err != nil {
		exitWithError("%s\n", err)
	}

	err = j.Sync()
	if err != nil {
		exitWithError("%s\n", err)
	}
}
