// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package quest

// ANSI codes to set graphic mode
const (
	setOff  = "\033[0m" // All attributes off
	setBold = "\033[1m" // Bold on
)

// The values by default are set to bold.
const lenAnsi = len(setBold) + len(setOff)
