// Copyright (c) 2012, SoundCloud Ltd.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/soundcloud/doozer-journal

package journal

// Operation values.
const (
	OpInvalid           = -1
	OpSet     Operation = 0
	OpDel               = 1
)

// Operation represents a doozerd operation.
type Operation int

// NewOperation returns an Operation based of it's string representation.
func NewOperation(opStr string) (op Operation) {
	switch opStr {
	case "set":
		op = OpSet
	case "del":
		op = OpDel
	default:
		op = OpInvalid
	}

	return
}

// String returns the string representation of an Operation.
func (op Operation) String() (str string) {
	switch op {
	case OpSet:
		str = "set"
	case OpDel:
		str = "del"
	default:
		str = "invalid"
	}

	return
}
