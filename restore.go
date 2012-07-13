// Copyright (c) 2012, SoundCloud Ltd.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/soundcloud/doozer-journal

package main

import (
	"github.com/soundcloud/doozer-journal/journal"
	"io"
	"os"
)

var cmdRestore = &Command{
	Name:      "restore",
	Desc:      "replays journal",
	UsageLine: "restore",
}

func init() {
	cmdRestore.Run = runRestore
}

func runRestore(cmd *Command, args []string) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0666)
	if err != nil {
		exitWithError("Unable to open file: %s\n", err)
	}

	j := journal.New(f)
	r := journal.NewReader(j)

	for {
		entry, err := r.ReadEntry()
		if err == io.EOF {
			break
		}
		if err != nil {
			exitWithError("%s\n", err)
		}

		switch entry.Op {
		case journal.OpSet:
			_, err = cmd.Conn.Set(entry.Path, -1, entry.Value)
			if err != nil {
				exitWithError("Error setting '%s' to '%s'\n", entry.Path, string(entry.Value))
			}
		case journal.OpDel:
			err = cmd.Conn.Del(entry.Path, -1)
			if err != nil {
				exitWithError("Error deleting '%s'\n", entry.Path)
			}
		default:
			exitWithError("Unknown operation %s\n", entry.Op)
		}
	}
}
