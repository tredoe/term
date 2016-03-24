// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build !plan9,!windows

package term

import (
	"io"
	"os"

	"github.com/kless/term/sys"
)

// A Terminal represents a general terminal interface.
type Terminal struct {
	mode modeType

	// Contain the state of a terminal, enabling to restore the original settings
	oldState, lastState sys.Termios

	// Window size
	size sys.Winsize

	fd, outFd int // File descriptor
	in        io.Reader
	out       io.Writer
}

// New creates a new terminal interface.
func New() (*Terminal, error) {
	return NewWith(os.Stdin, os.Stdout)
}

// NewWith creates a new terminal with in/out files
func NewWith(in, out *os.File) (*Terminal, error) {
	return NewWithAll(in, out, int(in.Fd()), int(out.Fd()))
}

// NewWithAll creates a new terminal with in/out and explicitly specified file descriptors
func NewWithAll(in io.Reader, out io.Writer, inFd, outFd int) (*Terminal, error) {
	t := &Terminal{in: in, out: out, fd: inFd, outFd: outFd}

	// Get the actual state
	if err := sys.Getattr(t.fd, &t.lastState); err != nil {
		return nil, os.NewSyscallError("sys.Getattr", err)
	}

	t.oldState = t.lastState // the actual state is copied to another one
	return t, nil
}

// == Restore
//

type State struct {
	wrap sys.Termios
}

// OriginalState returns the terminal's original state.
func (t *Terminal) OriginalState() State {
	return State{t.oldState}
}

// Restore restores the original settings for the term.
func (t *Terminal) Restore() error {
	if t.mode != 0 {
		if err := sys.Setattr(t.fd, sys.TCSANOW, &t.oldState); err != nil {
			return os.NewSyscallError("sys.Setattr", err)
		}
		t.lastState = t.oldState
		t.mode = 0
	}
	return nil
}

// Restore restores the settings from State.
func Restore(fd int, st State) error {
	if err := sys.Setattr(fd, sys.TCSANOW, &st.wrap); err != nil {
		return os.NewSyscallError("sys.Setattr", err)
	}
	return nil
}

// == Modes
//

// RawMode sets the terminal to something like the "raw" mode. Input is available
// character by character, echoing is disabled, and all special processing of
// terminal input and output characters is disabled.
//
// NOTE: in tty "raw mode", CR+LF is used for output and CR is used for input.
func (t *Terminal) RawMode() error {
	// Input modes - no break, no CR to NL, no NL to CR, no carriage return,
	// no strip char, no start/stop output control, no parity check.
	t.lastState.Iflag &^= (sys.BRKINT | sys.IGNBRK | sys.ICRNL | sys.INLCR |
		sys.IGNCR | sys.ISTRIP | sys.IXON | sys.PARMRK)

	// Output modes - disable post processing.
	t.lastState.Oflag &^= sys.OPOST

	// Local modes - echoing off, canonical off, no extended functions,
	// no signal chars (^Z,^C).
	t.lastState.Lflag &^= (sys.ECHO | sys.ECHONL | sys.ICANON | sys.IEXTEN | sys.ISIG)

	// Control modes - set 8 bit chars.
	t.lastState.Cflag &^= (sys.CSIZE | sys.PARENB)
	t.lastState.Cflag |= sys.CS8

	// Control chars - set return condition: min number of bytes and timer.
	// We want read to return every single byte, without timeout.
	t.lastState.Cc[sys.VMIN] = 1 // Read returns when one char is available.
	t.lastState.Cc[sys.VTIME] = 0

	// Put the terminal in raw mode after flushing
	if err := sys.Setattr(t.fd, sys.TCSAFLUSH, &t.lastState); err != nil {
		return os.NewSyscallError("sys.Setattr", err)
	}
	t.mode |= RawMode
	return nil
}

// EchoMode turns the echo mode.
func (t *Terminal) EchoMode(echo bool) error {
	if !echo {
		//t.lastState.Lflag &^= (sys.ECHO | sys.ECHOE | sys.ECHOK | sys.ECHONL)
		t.lastState.Lflag &^= sys.ECHO
	} else {
		//t.lastState.Lflag |= (sys.ECHO | sys.ECHOE | sys.ECHOK | sys.ECHONL)
		t.lastState.Lflag |= sys.ECHO
	}

	if err := sys.Setattr(t.fd, sys.TCSANOW, &t.lastState); err != nil {
		return os.NewSyscallError("sys.Setattr", err)
	}

	if echo {
		t.mode |= EchoMode
	} else {
		t.mode &^= EchoMode
	}
	return nil
}

// CharMode sets the terminal to single-character mode.
func (t *Terminal) CharMode() error {
	// Disable canonical mode, and set buffer size to 1 byte.
	t.lastState.Lflag &^= sys.ICANON
	t.lastState.Cc[sys.VTIME] = 0
	t.lastState.Cc[sys.VMIN] = 1

	if err := sys.Setattr(t.fd, sys.TCSANOW, &t.lastState); err != nil {
		return os.NewSyscallError("sys.Setattr", err)
	}
	t.mode |= CharMode
	return nil
}

// SetMode sets the terminal attributes given by state.
// Warning: The use of this function is not cross-system.
func (t *Terminal) SetMode(state sys.Termios) error {
	if err := sys.Setattr(t.fd, sys.TCSANOW, &state); err != nil {
		return os.NewSyscallError("sys.Setattr", err)
	}

	t.lastState = state
	t.mode |= OtherMode
	return nil
}

// == Utility
//

// Fd returns the file descriptor referencing the term.
func (t *Terminal) Fd() int {
	return t.fd
}

func (t *Terminal) Input() io.Reader {
	return t.in
}

func (t *Terminal) OutFd() int {
	return t.outFd
}

func (t *Terminal) Output() io.Writer {
	return t.out
}

// GetSize returns the size of the term.
func (t *Terminal) GetSize() (row, column int, err error) {
	if err = sys.GetWinsize(t.OutFd(), &t.size); err != nil {
		return
	}
	return int(t.size.Row), int(t.size.Col), nil
}
