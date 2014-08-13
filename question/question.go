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
func NewCustom(s *valid.Schema, prefixPrompt, prefixError, trueStr, falseStr string) *Question {
	t, err := term.New()
	if err != nil {
		panic(err)
	}

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

// Prompt sets a new prompt.
func (q *Question) Prompt(str string) *Question {
	q.prompt = str
	q.schema.Bydefault = ""
	q.schema.SetModifier(0)
	return q
}

// Default sets a value by default.
func (q *Question) Default(str string) *Question {
	q.schema.Bydefault = str
	return q
}

// Mod sets the modifier flags.
func (q *Question) Mod(flag valid.Modifier) *Question {
	q.schema.SetModifier(flag)
	return q
}

// Min sets the checking for the minimum length of a string,
// or the minimum value of a numeric type.
func (q *Question) Min(n interface{}) *Question {
	q.schema.SetMin(n)
	return q
}

// Max sets the checking for the maximum length of a string,
// or the maximum value of a numeric type.
func (q *Question) Max(n interface{}) *Question {
	q.schema.SetMax(n)
	return q
}

// Range sets the checking for the minimum and maximum lengths of a string,
// or the minimum and maximum values of a numeric type.
func (q *Question) Range(min, max interface{}) *Question {
	q.schema.SetRange(min, max)
	return q
}

// == Read

// read is the base to read and validate input.
func (q *Question) read() (iface interface{}, err error) {
	var hadError bool
	line, err := q.newLine()
	if err != nil {
		return nil, err
	}

	for {
		input, err := line.Read()
		if err != nil {
			return nil, err
		}

		switch q.schema.DataType {
		case yoda.T_bool:
			iface, err = valid.Bool(q.schema, input)
		case yoda.T_int64:
			iface, err = valid.Int64(q.schema, input)
		case yoda.T_uint64:
			iface, err = valid.Uint64(q.schema, input)
		case yoda.T_float64:
			iface, err = valid.Float64(q.schema, input)
		case yoda.T_string:
			iface, err = valid.String(q.schema, input)
		default:
			panic("unimplemented")
		}

		// Validation
		if err == valid.ErrRequired {
			os.Stderr.Write(readline.DelLine_CR)
			fmt.Fprintf(os.Stderr, "%s%s", q.prefixError, err)
			term.Output.Write(readline.CursorUp)
			hadError = true
			continue
		}
		// Error of type.
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
		return iface, nil
	}
}

// ReadBool prints the prompt waiting to get a string that represents a boolean.
func (q *Question) ReadBool() (bool, error) {
	q.schema.DataType = yoda.T_bool

	iface, err := q.read()
	if err != nil {
		return false, err
	}
	return iface.(bool), nil
}

// ReadInt64 prints the prompt waiting to get an integer number.
func (q *Question) ReadInt64() (int64, error) {
	q.schema.DataType = yoda.T_int64

	iface, err := q.read()
	if err != nil {
		return 0, err
	}
	return iface.(int64), nil
}

// ReadUint64 prints the prompt waiting to get an unsigned integer number.
func (q *Question) ReadUint64() (uint64, error) {
	q.schema.DataType = yoda.T_uint64

	iface, err := q.read()
	if err != nil {
		return 0, err
	}
	return iface.(uint64), nil
}

// ReadFloat64 prints the prompt waiting to get a floating-point number.
func (q *Question) ReadFloat64() (float64, error) {
	q.schema.DataType = yoda.T_float64

	iface, err := q.read()
	if err != nil {
		return 0, err
	}
	return iface.(float64), nil
}

// ReadString prints the prompt waiting to get a string.
func (q *Question) ReadString() (string, error) {
	q.schema.DataType = yoda.T_string

	iface, err := q.read()
	if err != nil {
		return "", err
	}
	return iface.(string), nil
}

// == Slices

// ReadInt64Slice reads multiple int64.
// You have to press Enter twice to finish.
func (q *Question) ReadInt64Slice() (values []int64, err error) {
	q.schema.DataType = yoda.T_Int64Slice

	if _, err = q.newLine(); err != nil {
		return nil, err
	}

	q.schema.IsSlice = true
	term.Output.Write(readline.CRLF)

	for {
		v, err := q.ReadInt64()
		if err != nil {
			return nil, err
		}
		if v == 0 {
			break
		}
		values = append(values, v)
	}

	q.schema.IsSlice = false
	return
}

// ReadUint64Slice reads multiple uint64.
// You have to press Enter twice to finish.
func (q *Question) ReadUint64Slice() (values []uint64, err error) {
	q.schema.DataType = yoda.T_Uint64Slice

	if _, err = q.newLine(); err != nil {
		return nil, err
	}

	q.schema.IsSlice = true
	term.Output.Write(readline.CRLF)

	for {
		v, err := q.ReadUint64()
		if err != nil {
			return nil, err
		}
		if v == 0 {
			break
		}
		values = append(values, v)
	}

	q.schema.IsSlice = false
	return
}

// ReadFloat64Slice reads multiple float64.
// You have to press Enter twice to finish.
func (q *Question) ReadFloat64Slice() (values []float64, err error) {
	q.schema.DataType = yoda.T_Float64Slice

	if _, err = q.newLine(); err != nil {
		return nil, err
	}

	q.schema.IsSlice = true
	term.Output.Write(readline.CRLF)

	for {
		v, err := q.ReadFloat64()
		if err != nil {
			return nil, err
		}
		if v == 0 {
			break
		}
		values = append(values, v)
	}

	q.schema.IsSlice = false
	return
}

// ReadSliceString reads multiple strings.
// You have to press Enter twice to finish.
func (q *Question) ReadStringSlice() (values []string, err error) {
	q.schema.DataType = yoda.T_StringSlice

	if _, err = q.newLine(); err != nil {
		return nil, err
	}

	q.schema.IsSlice = true
	term.Output.Write(readline.CRLF)

	for {
		v, err := q.ReadString()
		if err != nil {
			return nil, err
		}
		if v == "" {
			break
		}
		values = append(values, v)
	}

	q.schema.IsSlice = false
	return
}

// == Choices

var (
	down2 = []byte{13, 10, 13, 10}
	up2   = []byte(fmt.Sprintf("%s%s", readline.CursorUp, readline.CursorUp))
)

// ChoiceInt prints the prompt waiting to get an int that is in the slice.
func (q *Question) ChoiceInt(choices []int) (int, error) {
	term.Output.Write(down2)
	fmt.Fprintf(term.Output, "   >>> %v", choices)
	term.Output.Write(up2)

	for {
		value, err := q.ReadInt64()
		if err != nil {
			return 0, err
		}
		for _, v := range choices {
			if v == int(value) {
				return int(value), nil
			}
		}

		os.Stderr.Write(readline.DelLine_CR)
		fmt.Fprintf(os.Stderr, "%s%s", q.prefixError, "invalid choice")
		term.Output.Write(readline.CursorUp)
	}
}

// ChoiceFloat64 prints the prompt waiting to get a float64 that is in the slice.
func (q *Question) ChoiceFloat64(choices []float64) (float64, error) {
	term.Output.Write(down2)
	fmt.Fprintf(term.Output, "   >>> %v", choices)
	term.Output.Write(up2)

	for {
		value, err := q.ReadFloat64()
		if err != nil {
			return 0, err
		}
		for _, v := range choices {
			if v == float64(value) {
				return float64(value), nil
			}
		}

		os.Stderr.Write(readline.DelLine_CR)
		fmt.Fprintf(os.Stderr, "%s%s", q.prefixError, "invalid choice")
		term.Output.Write(readline.CursorUp)
	}
}

// ChoiceString prints the prompt waiting to get a string that is in the slice.
func (q *Question) ChoiceString(choices []string) (string, error) {
	term.Output.Write(down2)
	fmt.Fprintf(term.Output, "   >>> %v", choices)
	term.Output.Write(up2)

	for {
		value, err := q.ReadString()
		if err != nil {
			return "", err
		}
		for _, v := range choices {
			if v == string(value) {
				return string(value), nil
			}
		}

		os.Stderr.Write(readline.DelLine_CR)
		fmt.Fprintf(os.Stderr, "%s%s", q.prefixError, "invalid choice")
		term.Output.Write(readline.CursorUp)
	}
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

// SetRFCEmail sets an email based in RFC like regular expression.
func (q *Question) SetRFCEmail() *Question {
	q.val = validate.NewRFCEmail(q.flag)
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

			if q.schema.DataType != yoda.T_bool {
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

// == Terminal handler

// ExitAtCtrlD exits when it is pressed Ctrl-D, with value n.
func (q *Question) ExitAtCtrlD(n int) {
	go func() {
		select {
		case <-readline.ChanCtrlD:
			term.Output.Write(readline.DelLine_CR)
			os.Exit(n)
		}
	}()
}

// Restore restores terminal settings.
func (q *Question) Restore() error { return q.term.Restore() }

// == Utility

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
