// Copyright (c) 2012, SoundCloud Ltd.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/soundcloud/doozer-journal

package journal

import (
	"strconv"
	"strings"
)

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
func NewEntry(r int64, op Operation, p string, v []byte) (entry JournalEntry) {
	return JournalEntry{Rev: r, Op: op, Path: p, Value: v}
}

// NewEntryFromLog takes a log line in form of a string and constructs a JournalEntry
// instance from it.
func NewEntryFromLog(log string) (entry JournalEntry, err error) {
	l := strings.SplitN(log, FIELD_SEPARATOR, 4)

	rev, err := strconv.ParseInt(l[0], 10, 64)
	if err != nil {
		return
	}

	val, err := strconv.Unquote(l[3])
	if err != nil {
		return
	}

	entry = NewEntry(rev, NewOperation(l[1]), l[2], []byte(val))

	return
}

// ToLog returns the string representation of a JournalEntry in the journal log file.
func (e JournalEntry) ToLog() string {
	rev := strconv.FormatInt(e.Rev, 10)
	val := strconv.Quote(string(e.Value))

	return strings.Join([]string{rev, e.Op.String(), e.Path, val}, FIELD_SEPARATOR)
}
