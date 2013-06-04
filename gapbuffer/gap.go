// Copyright 2013 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package gapbuffer implements the gap buffer.
package gapbuffer

import (
	"fmt"
	//"unicode/utf8"

	//"bitbucket.org/ares/term"
)

// Buffer size by default.
const (
	_BUFFER_CAP = 4096
	_BUFFER_LEN = 64 // Initial length
)

type WriteMode int

const (
	Insert WriteMode = iota
	Overwrite
)

// A GapBuffer represents a gap buffer.
type GapBuffer struct {
	size     int // Number of characters added
	gapStart int
	gapEnd   int

	bufEnd int
	cursor int
	mode   WriteMode

	buf     []rune

/*	columns   int // Number of columns for actual window
	promptLen int
	pos       int    // Pointer position into buffer
	size      int    // Amount of characters added
*/
}

// New creates and initializes a new GapBuffer using values by default.
func New() *GapBuffer {
	return NewGapBuffer(make([]rune, _BUFFER_LEN, _BUFFER_CAP))
}

// NewGapBuffer creates and initializes a new GapBuffer using buf as its initial contents.
func NewGapBuffer(buf []rune) *GapBuffer {
	lastIndex := len(buf) - 1

	return &GapBuffer{
		buf:      buf,
		bufEnd:   lastIndex,
		gapEnd:   lastIndex,
		gapStart: lastIndex / 2,
	}
}

// NextChar moves the cursor to next character.
func (b *GapBuffer) NextChar() bool {
	if b.cursor < len(b.buf) {
		b.cursor++
		b.gapStart++
		b.gapEnd++

		b.buf[b.cursor] = b.buf[b.gapEnd]
		b.buf[b.gapEnd] = 0
		return true
	}
	return false
}

// PrevChar moves the cursor to previous character.
func (b *GapBuffer) PrevChar() bool {
	if b.cursor > 0 {
		b.cursor--
		b.buf[b.gapEnd] = b.buf[b.cursor]
		b.buf[b.cursor] = 0

		b.gapStart--
		b.gapEnd--
		return true
	}
	return false
}

// NextWord moves the cursor to next word.
func (b *GapBuffer) NextWord() {
	for ok := false; ; {
		ok = b.NextChar()
		if !ok || b.buf[b.cursor] == 32 {
			return
		}
	}
}

func (b *GapBuffer) Show() {
	
}

// PrevWord moves the cursor to previous word.
func (b *GapBuffer) PrevWord() {
	for ok := false; ; {
		ok = b.PrevChar()
		if !ok || b.buf[b.cursor-1] == 32 {
			return
		}
	}
}

// InsertChar inserts a character in the cursor position.
func (b *GapBuffer) InsertChar(r rune) error {
	b.buf[b.cursor] = r
	b.cursor++

	return nil
}

// InsertChars inserts several characters.
func (b *GapBuffer) InsertChars(runes []rune) error {
	for _, r := range runes {
		if err := b.InsertChar(r); err != nil {
			return err
		}
	}
	return nil
}

// Print prints the buffer.
func (b *GapBuffer) Print() {
	fmt.Printf(" Cursor:%*d · Gap start:%*d · Gap end:%*d  [",
		3, b.cursor, 3, b.gapStart, 3, b.gapEnd)

	for i := 0; i < len(b.buf); i++ {
		if i > b.gapStart && i < b.gapEnd {
			fmt.Print("_")
		} else if i == b.gapStart || i == b.gapEnd {
			fmt.Print("|")
		} else {
			if b.buf[i] == 0 {
				fmt.Print("*")
			} else {
				fmt.Printf("%c", b.buf[i])
			}
		}
	}
	fmt.Println("]")
}

