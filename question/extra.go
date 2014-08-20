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
	//"github.com/kless/yoda"
	"github.com/kless/yoda/valid"
)

// ValidFunc is the type of the function called to validate extra types.
type ValidFunc func(s *valid.Schema, str string) (string, error)

// readExtra is the base to read and validate extra types.
func (q *Question) readExtra(validFn ValidFunc) (string, error) {
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

		value, err := validFn(q.schema, input)
		if err != nil {
			return "", err
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

// ReadEmail
func (q *Question) ReadEmail() (value string, err error) {
	for {
		if value, err = q.ReadString(); err != nil {
			return
		}

		if value, err = valid.Email(q.schema, value); err != nil {
			os.Stderr.Write(readline.DelLine_CR)
			fmt.Fprintf(os.Stderr, "%s%s", q.prefixError, err)

			term.Output.Write(readline.CursorUp)
			//hadError = true
			continue
		}

		return //value, err
	}

//	return q.readExtra(valid.Email(q.schema, ""))
}

// == Regexp

/*
// NewRegexp sets a regular expression "re" with a name to be identified.
func (q *Question) NewRegexp(name, re string) *Question {
	q.val = validate.NewRegexp(q.flag, name, re)
	return q
}

// SetBasicEmail sets an email based into a basic format like regular expression.
func (q *Question) SetBasicEmail() *Question {
	q.val = validate.NewBasicEmail(q.flag)
	return q
}

// AddRegexp adds another regular expression.
// Whether the name to identify the reg.exp. has been set, it is changed.
func (q *Question) AddRegexp(name, re string) *Question {
	q.val.AddRegexp(name, re)
	return q
}

// Read prints the prompt waiting to get a string that matches with a set of
// regular expressions.
func (q *Question) Read() (string, error) {
	if q.defaultVal != nil {
		q.val.SetDefault(q.defaultVal)
		q.suffixPrompt = msgForDefault(q.defaultVal)
	}

	value, err := q.read(q.newLine(), q.val)
	q.clean()
	return value.(string), err
}
*/
