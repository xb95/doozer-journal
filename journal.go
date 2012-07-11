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

	rev := cmd.Rev
	for {
		ev, err := cmd.Conn.Wait("/**", rev+1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to wait on event: %s\n", err.Error())
			os.Exit(1)
		}

		var entry *journal.Entry
		if ev.IsSet() {
			entry = journal.NewEntry(ev.Rev, journal.OpSet, ev.Path, ev.Body)
		} else {
			entry = journal.NewEntry(ev.Rev, journal.OpDel, ev.Path, []byte{})
		}

		err = j.Append(entry)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}

		b, err := journal.Marshal(entry)
		if err != nil {
			return
		}

		fmt.Fprintf(os.Stdout, "%s\n", string(b))

		rev = ev.Rev
	}

}
