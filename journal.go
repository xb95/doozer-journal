// Copyright (c) 2012, SoundCloud Ltd.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/soundcloud/doozer-journal

package main

import (
	"github.com/soundcloud/doozer-journal/journal"
	"os"
	"path"
	"strings"
)

var cmdJournal = &Command{
	Name:      "journal",
	Desc:      "takes an initial snapshot & journals mutations",
	UsageLine: "journal",
}

func init() {
	cmdJournal.Run = runJournal
}

func runJournal(cmd *Command, args []string) {
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

	watchroot := path.Join(root, "**")

	rev := cmd.Rev
	for {
		ev, err := cmd.Conn.Wait(watchroot, rev+1)
		if err != nil {
			exitWithError("Unable to wait on event: %s\n", err)
		}

		// Advance in time!
		rev = ev.Rev

		var entry *journal.Entry
		if ev.IsSet() {
			entry = journal.NewEntry(ev.Rev, journal.OpSet, ev.Path, ev.Body)
		} else if ev.IsDel() {
			entry = journal.NewEntry(ev.Rev, journal.OpDel, ev.Path, []byte{})
		} else {
			continue
		}

		err = j.Append(entry)
		if err != nil {
			exitWithError("%s\n", err)
		}

		b, err := journal.Marshal(entry)
		if err != nil {
			return
		}

		if !strings.HasPrefix(ev.Path, journal.INTERNAL_PREFIX) {
			logInfo("%s\n", string(b))
		}

	}
}
