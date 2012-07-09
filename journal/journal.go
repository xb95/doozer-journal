// Copyright (c) 2012, SoundCloud Ltd.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/soundcloud/doozer-journal

package journal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// INTERNAL_PREFIX is the path fragment which identifies paths used by doozerd
// internally.
const INTERNAL_PREFIX string = "/ctl"

// ENTRY_SEPARATOR is the delimiter used to separate journal entries in the file.
const ENTRY_SEPARATOR rune = '\n'

// Journal represents an append-only log which accepts an either writable or readable
// file.
type Journal struct {
	File         *os.File
	SyncInterval time.Duration
	SyncOps      int64
	opChan       chan JournalEntry
}

// New returns a Journal instance with sane sync defaults.
func New(logfile *os.File) Journal {
	return Journal{logfile, 10 * time.Second, 100, make(chan JournalEntry, 1024)}
}

// Append writes a JournalEntry to the end of the journal log.
func (j Journal) Append(entry JournalEntry) (err error) {
	if !strings.HasPrefix(entry.Path, INTERNAL_PREFIX) {
		_, err := j.File.Write([]byte(entry.ToLog() + string(ENTRY_SEPARATOR)))
		if err != nil {
			return fmt.Errorf("Unable to append '%s' to journal: %s", entry.ToLog(), err.Error())
		}

		j.opChan <- entry
	}

	return
}

// Reads the Journal log file and emits every entry over the entries chan. All errors
// are send to the errChan.
func (j Journal) Read(entries chan JournalEntry, errChan chan error) {
	reader := bufio.NewReader(j.File)

	for {
		line, err := reader.ReadString(byte(ENTRY_SEPARATOR))
		if err == io.EOF {
			errChan <- nil

			close(entries)
			break
		}
		if err != nil && err != io.EOF {
			errChan <- err

			close(entries)
			break
		}

		cleanLine := strings.Trim(line, string(ENTRY_SEPARATOR))
		entry, err := NewEntryFromLog(cleanLine)
		if err != nil {
			errChan <- err

			close(entries)
			break
		}

		entries <- entry
	}
}

// Sync forces an fsync() on the journal log file.
func (j Journal) Sync() (err error) {
	err = j.File.Sync()
	if err != nil {
		return fmt.Errorf("Unable to sync journal: %s", err.Error())
	}

	return
}

// SyncLoop schedules Sync calls based on the treshholds defined for SyncInterval &
// SyncOps.
func (j Journal) SyncLoop() {
	tick := time.Tick(j.SyncInterval)
	var opCounter int64 = 0

	for {
		select {
		case _, ok := <-j.opChan:
			if !ok {
				return
			}

			opCounter += 1

			if opCounter >= j.SyncOps {
				j.Sync()
				opCounter = 0
			}
		case <-tick:
			println("sync after time!")
			j.Sync()
			opCounter = 0
		}
	}
}
