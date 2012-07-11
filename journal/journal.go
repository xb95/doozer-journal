// Copyright (c) 2012, SoundCloud Ltd.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/soundcloud/doozer-journal

package journal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// INTERNAL_PREFIX is the path fragment which identifies paths used by doozerd
// internally.
const INTERNAL_PREFIX = "/ctl"

// ENTRY_END marks the end of the logline representation of an entry.
const ENTRY_END = "\n"

// FIELD_SEPARATOR is the delimiter used to separate fields inside of an entry.
const FIELD_SEPARATOR string = "|"

// JournalEntry represents a journal log entry.
type JournalEntry struct {
	Rev   int64
	Op    Operation
	Path  string
	Value []byte
}

// NewEntry returns a JournalEntry instance.
func NewEntry(r int64, op Operation, p string, v []byte) (entry *JournalEntry) {
	return &JournalEntry{Rev: r, Op: op, Path: p, Value: v}
}

// Journal represents an append-only log which accepts an either writable or readable
// file.
type Journal struct {
	File         *os.File
	SyncInterval time.Duration
	SyncOps      int64
	opCounter    int64
}

// New returns a Journal instance with sane sync defaults.
func New(logfile *os.File) (j *Journal) {
	j = &Journal{logfile, 10 * time.Second, 100, 0}

	go j.syncLoop()

	return
}

// Append writes a JournalEntry to the end of the journal log.
func (j *Journal) Append(entry *JournalEntry) (err error) {
	if !strings.HasPrefix(entry.Path, INTERNAL_PREFIX) {
		var payload []byte

		payload, err = Marshal(entry)
		if err != nil {
			return
		}

		length := fmt.Sprintf("%08d", len(payload))
		line := length + " " + string(payload) + ENTRY_END

		_, err := j.File.Write([]byte(line))
		if err != nil {
			return fmt.Errorf("Unable to append '%s' to journal: %s", string(payload), err.Error())
		}

		if j.opCounter >= j.SyncOps {
			j.Sync()
			j.opCounter = 0
		}
	}

	return
}

// Sync forces an fsync() on the journal log file.
func (j *Journal) Sync() (err error) {
	err = j.File.Sync()
	if err != nil {
		return fmt.Errorf("Unable to sync journal: %s", err.Error())
	}

	return
}

// syncLoop schedules Sync calls based on the treshholds defined for SyncInterval &
// SyncOps.
func (j *Journal) syncLoop() {
	tick := time.Tick(j.SyncInterval)

	for _ = range tick {
		if j.opCounter > 0 {
			j.Sync()
			j.opCounter = 0
		}
	}
}

type Reader struct {
	file   *os.File
	offset int64
}

func NewReader(j *Journal) (r *Reader) {
	r = &Reader{file: j.File}

	return
}

func (r *Reader) ReadEntry() (entry *JournalEntry, err error) {
	l := make([]byte, 8, 8)
	n, err := r.file.ReadAt(l, r.offset)
	if err != nil {
		return
	}

	r.offset += int64(n)

	space := make([]byte, 1, 1)
	n, err = r.file.ReadAt(space, r.offset)
	if err != nil {
		return
	}

	if string(space) != " " {
		return nil, fmt.Errorf("corrupted journal file: missing space")
	}

	r.offset += int64(n)

	length, err := strconv.Atoi(string(l))
	if err != nil {
		return
	}

	payload := make([]byte, length, length)
	n, err = r.file.ReadAt(payload, r.offset)
	if err != nil {
		return
	}

	r.offset += int64(n)

	entry, err = Unmarshal(payload)
	if err != nil {
		return
	}

	end := make([]byte, len(ENTRY_END), len(ENTRY_END))
	n, err = r.file.ReadAt(end, r.offset)
	if err != nil {
		return
	}

	r.offset += int64(n)

	if string(end) != ENTRY_END {
		return nil, fmt.Errorf("corrupted journal file: frame doesn't end with %s", ENTRY_END)
	}

	return
}

func Marshal(entry *JournalEntry) (payload []byte, err error) {
	rev := strconv.FormatInt(entry.Rev, 10)
	val := string(entry.Value)
	str := strings.Join([]string{rev, entry.Op.String(), entry.Path, val}, FIELD_SEPARATOR)

	payload = []byte(str)

	return
}

func Unmarshal(payload []byte) (entry *JournalEntry, err error) {
	l := strings.SplitN(string(payload), FIELD_SEPARATOR, 4)

	rev, err := strconv.ParseInt(l[0], 10, 64)
	if err != nil {
		return
	}

	entry = NewEntry(rev, NewOperation(l[1]), l[2], []byte(l[3]))

	return
}
