// Copyright (c) 2012, SoundCloud Ltd.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/soundcloud/doozer-journal

package main

import (
	"fmt"
	"github.com/soundcloud/doozer-journal/coordinator"
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

	entries := make(chan coordinator.Entry, 1024)
	errChan := make(chan error)

	go coordinator.Watch(cmd.Conn, cmd.Rev, entries, errChan)

	entryHandler(j, entries, errChan)
}

func entryHandler(j *journal.Journal, entries chan coordinator.Entry, errChan chan error) {
	for {
		select {
		case e, ok := <-entries:
			if !ok {
				return
			}

			var entry *journal.Entry
			if e.IsSet {
				entry = journal.NewEntry(e.Rev, journal.OpSet, e.Path, e.Value)
			} else {
				entry = journal.NewEntry(e.Rev, journal.OpDel, e.Path, []byte{})
			}

			err := j.Append(entry)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}
			b, err := journal.Marshal(entry)
			if err != nil {
				return
			}

			fmt.Fprintf(os.Stdout, "%s\n", string(b))
		case err := <-errChan:
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}
		}
	}
}
