// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package question provides functions for printing questions and validating answers.
//
// The test enables to run in interactive mode (flag -user), or automatically (by
// default) where the answers are written directly every "n" seconds (flag -t).
//
// Note: the values for the input and output are got from the package base "term".
package question

import (
	"fmt"
	"os"

	"github.com/kless/term"
	"github.com/kless/term/readline"
	"github.com/kless/yoda"
	"github.com/kless/yoda/valid"
)

// A Question represents a question.
type Question struct {
	term   *term.Terminal
	schema *valid.Schema // Validation schema

	prefixError  string // String before of any error message
	prefixPrompt string // String before of the question
	suffixPrompt string // String after of the question, for values by default
	prompt       string // The question

	// Strings that represent booleans
	trueStr  string
	falseStr string
}

// NewCustom returns a Question with the given arguments; if any is empty,
// it is used the values by default.
// Panics if the terminal can not be set.
//
// prefixPrompt is the text placed before of the prompt.
// prefixError is placed before of show any error.
//
// trueStr and falseStr are the strings to be showed when the question
// needs a boolean like answer and it is used a value by default.
//
// The terminal is changed to raw mode so have to use (*Question) Restore()
// at finish.
//
// Handles interrupts CTRL-C to start a new line and CTRL-D to exit.
func NewCustom(s *valid.Schema, prefixPrompt, prefixError, trueStr, falseStr string) *Question {
	t, err := term.New()
	if err != nil {
		panic(err)
	}

	// Interrupt handling.
	go func() {
		for {
			select {
			case <- readline.ChanCtrlC:
			case <-readline.ChanCtrlD:
				term.Output.Write(readline.DelLine_CR)
				os.Exit(2)
			}
		}
	}()

	if s == nil {
		s = valid.NewSchema(0)
	}
	if prefixPrompt == "" {
		prefixPrompt = _PREFIX
	}
	if prefixError == "" {
		prefixError = _PREFIX_ERR
	}
	if trueStr == "" {
		trueStr = _STR_TRUE
	}
	if falseStr == "" {
		falseStr = _STR_FALSE
	}

	valid.SetBoolStrings(map[string]bool{trueStr: true, falseStr: false})

	return &Question{
		term:   t,
		schema: s,

		prefixError:  prefixError,
		prefixPrompt: prefixPrompt,
		trueStr:      trueStr,
		falseStr:     falseStr,
	}
}

// Values by default for a Question.
const (
	_PREFIX       = " + "
	_PREFIX_MULTI = "   * "
	_PREFIX_ERR   = "  [!] "
	_STR_TRUE     = "y"
	_STR_FALSE    = "n"
)

// New returns a Question using values by default.
// The terminal is changed to raw mode so have to use (*Question) Restore()
// at finish.
func New() *Question {
	return NewCustom(valid.NewSchema(0), _PREFIX, _PREFIX_ERR, _STR_TRUE, _STR_FALSE)
}

// Restore restores terminal settings.
func (q *Question) Restore() error { return q.term.Restore() }

// Prompt sets a new prompt.
func (q *Question) Prompt(str string) *Question {
	q.prompt = str
	q.schema.Bydefault = ""
	q.schema.SetChecker(0)
	return q
}

// Default sets a value by default.
func (q *Question) Default(str string) *Question {
	q.schema.Bydefault = str
	return q
}

// Check sets the checker flags.
func (q *Question) Check(flag valid.Checker) *Question {
	q.schema.SetChecker(flag)
	return q
}

// Min sets the checking for the minimum length of a string,
// or the minimum value of a numeric type.
// The valid types for the aregument are: int, float64.
func (q *Question) Min(n interface{}) *Question {
	q.schema.SetMin(n)
	return q
}

// Max sets the checking for the maximum length of a string,
// or the maximum value of a numeric type.
// The valid types for the aregument are: int, float64.
func (q *Question) Max(n interface{}) *Question {
	q.schema.SetMax(n)
	return q
}

// Range sets the checking for the minimum and maximum lengths of a string,
// or the minimum and maximum values of a numeric type.
// The valid types for the areguments are: int, float64.
func (q *Question) Range(min, max interface{}) *Question {
	q.schema.SetRange(min, max)
	return q
}

// == Printing

// The values by default are set to bold.
const lenAnsi = len(readline.ANSI_SET_BOLD) + len(readline.ANSI_SET_OFF)

// newLine gets a line type ready to show questions.
func (q *Question) newLine() (*readline.Line, error) {
	fullPrompt := ""
	extraChars := 0

	if !q.schema.IsSlice {
		fullPrompt = q.prefixPrompt + q.prompt

		// Add spaces
		if fullPrompt[len(fullPrompt)-1] == '?' {
			fullPrompt += " "
		} else if !q.schema.IsSlice {
			fullPrompt += ": "
		}

		// Default value
		if q.schema.Bydefault != "" {
			extraChars = lenAnsi // The value by default uses ANSI characters.

			if q.schema.DataType() != yoda.T_bool {
				q.suffixPrompt = fmt.Sprintf("[%s%s%s] ",
					readline.ANSI_SET_BOLD,
					q.schema.Bydefault,
					readline.ANSI_SET_OFF,
				)
			} else {
				b, err := valid.Bool(q.schema, q.schema.Bydefault)
				if err != nil {
					return nil, err
				}
				if b {
					q.suffixPrompt = fmt.Sprintf("[%s%s%s/%s] ",
						readline.ANSI_SET_BOLD,
						q.trueStr,
						readline.ANSI_SET_OFF,
						q.falseStr,
					)
				} else {
					q.suffixPrompt = fmt.Sprintf("[%s/%s%s%s] ",
						q.trueStr,
						readline.ANSI_SET_BOLD,
						q.falseStr,
						readline.ANSI_SET_OFF,
					)
				}
			}

			fullPrompt += q.suffixPrompt
		}
	} else {
		fullPrompt = _PREFIX_MULTI
	}

	// No history
	return readline.NewLine(q.term, fullPrompt, q.prefixError, extraChars, nil)
}

// PrintAnswer prints values returned by a Question.
func PrintAnswer(i interface{}, err error) {
	term.Output.Write(readline.DelLine_CR)

	if err == nil {
		msg := "  answer: "
		if _, ok := i.(string); ok {
			msg += "%q\r\n"
		} else {
			msg += "%v\r\n"
		}
		fmt.Fprintf(term.Output, msg, i)

	} else if err != readline.ErrCtrlD {
		fmt.Fprintf(term.Output, "%s\r\n", err)
	}
}
