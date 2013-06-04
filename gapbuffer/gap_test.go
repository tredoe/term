// Copyright 2013 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gapbuffer

import (
	"fmt"
	"testing"
)

func TestGap(t *testing.T) {
	buf := NewGapBuffer(make([]rune, 32, 128))
	buf.Print()

	buf.InsertChars([]rune("0123456789"))
	buf.Print()
	fmt.Println()

	for i := 0; i < 5; i++ {
		buf.PrevChar()
		buf.Print()
	}
	//buf.Print()
}
