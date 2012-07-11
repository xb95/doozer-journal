// Copyright (c) 2012, SoundCloud Ltd.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/soundcloud/doozer-journal

package main

import (
	"fmt"
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
			fmt.Fprintf(os.Stderr, "Journal already exists at %s\n", file)
		} else {
			fmt.Fprintf(os.Stderr, "Unable to open journal: %s\n", err.Error())
		}

		os.Exit(1)
	}

	j := journal.New(f)
	err = snapshot(cmd.Conn, cmd.Rev, j)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	err = j.Sync()
	if err != nil {
		fmt.Fprint(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
