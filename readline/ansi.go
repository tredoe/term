// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// The references about ANSI Escape sequences have been got from
// http://ascii-table.com/ansi-escape-sequences.php and
// http://www.termsys.demon.co.uk/vtansi.htm

package readline

// Characters
var (
	_CR   = []byte{13}     // Carriage return -- \r
	CRLF  = []byte{13, 10} // CR+LF is used for a new line in raw mode -- \r\n
	ctrlC = []rune("^C")
	ctrlD = []rune("^D")
)

// ANSI terminal escape controls
var (
	// == Cursor control
	CursorUp       = []byte("\033[A") // Up
	cursorDown     = []byte("\033[B") // Down
	cursorForward  = []byte("\033[C") // Forward
	cursorBackward = []byte("\033[D") // Backward

	toNextLine     = []byte("\033[E") // To next line
	toPreviousLine = []byte("\033[F") // To previous line

	// == Erase Text
	delScreenToUpper = []byte("\033[2J\033[0;0H") // Erase the screen; move upper

	delToRight       = []byte("\033[0K")       // Erase to right
	DelLine_CR       = []byte("\033[2K\r")     // Erase line; carriage return
	delLine_cursorUp = []byte("\033[2K\033[A") // Erase line; cursor up

	//delChar      = []byte("\033[1X") // Erase character
	delChar      = []byte("\033[P") // Delete character, from current position
	delBackspace = []byte("\033[D\033[P")

	// == Misc.
	//insertChar  = []byte("\033[@")   // Insert CHaracter
	//setLineWrap = []byte("\033[?7h") // Enable Line Wrap
)
