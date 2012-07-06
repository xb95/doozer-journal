// Copyright (c) 2012, SoundCloud Ltd.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/soundcloud/doozer-journal

package coordinator

import (
	"fmt"
	"github.com/soundcloud/doozer"
)

// Entry represents a doozerd entry.
type Entry struct {
	Rev   int64
	Path  string
	Value []byte
	IsSet bool
}

// Walk reads the complete doozerd state at the passed revision and emits every
// entry(includes only files) over the entries channel. All errors are send to
// errChan.
func Walk(conn *doozer.Conn, rev int64, entries chan Entry, errChan chan error) {
	doozer.Walk(conn, rev, "/", func(path string, info *doozer.FileInfo, e error) (err error) {
		if e != nil {
			errChan <- fmt.Errorf("Error walking tree: %s\n", e.Error())
		}

		if !info.IsDir {
			val, _, err := conn.Get(path, &rev)
			if err != nil {
				errChan <- fmt.Errorf("Error getting value for '%s': %s\n", path, err.Error())
			}

      entries <- Entry{Rev: info.Rev, Path: path, Value: val, IsSet: true}
		}

		return
	})

	close(entries)
}
