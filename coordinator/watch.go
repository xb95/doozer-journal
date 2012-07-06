// Copyright (c) 2012, SoundCloud Ltd.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/soundcloud/doozer-journal

package coordinator

import (
	"fmt"
	"github.com/soundcloud/doozer"
)

// Watch tracks doozerd mutations and emits instanciated Entries over the entries
// channel. All errors are send to errChan.
func Watch(conn *doozer.Conn, rev int64, entries chan Entry, errChan chan error) {
	for {
		ev, err := conn.Wait("/**", rev+1)
		if err != nil {
			errChan <- fmt.Errorf("Unable to wait on event: %s\n", err.Error())
		}

		entries <- Entry{Rev: ev.Rev, Path: ev.Path, Value: ev.Body, IsSet: ev.IsSet()}

		rev = ev.Rev
	}
}
