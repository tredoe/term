// Copyright 2014 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package question

import (
	"fmt"
	"os"

	"github.com/kless/term"
	"github.com/kless/term/readline"
	"github.com/kless/yoda/valid"
)

type extraType int

const (
	_ extraType = iota
	t_email
	t_url
)

// readExtra is the base to read and validate extra types.
func (q *Question) readExtra(t extraType) (value string, err error) {
	var hadError bool
	line, err := q.newLine()
	if err != nil {
		return "", err
	}

	for {
		input, err := line.Read()
		if err != nil {
			return "", err
		}

		switch t {
		case t_email:
			value, err = valid.Email(q.schema, input)
		case t_url:
			value, err = valid.URL(q.schema, input)
		default:
			panic("unimplemented")
		}

		if err != nil {
			os.Stderr.Write(readline.DelLine_CR)
			fmt.Fprintf(os.Stderr, "%s%s", q.prefixError, err)
			term.Output.Write(readline.CursorUp)
			hadError = true
			continue
		}

		if hadError {
			os.Stderr.Write(readline.DelLine_CR)
		}
		return value, nil
	}
}

// ReadEmail prints the prompt waiting to get an email.
func (q *Question) ReadEmail() (value string, err error) {
	return q.readExtra(t_email)
}

// ReadURL prints the prompt waiting to get an URL.
func (q *Question) ReadURL() (value string, err error) {
	return q.readExtra(t_url)
}
