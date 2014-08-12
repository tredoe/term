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

	isSlice bool // If accept multiple answers.
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
func New() *Question {
	return NewCustom(valid.NewSchema(0), _PREFIX, _PREFIX_ERR, _STR_TRUE, _STR_FALSE)
}

// TODO: remove
// clean removes data related to the prompt.
func (q *Question) clean() {
	q.prompt = ""
	//q.suffixPrompt = ""
	//q.defaultVal = nil
	//q.flag = 0
}

// Prompt sets a new prompt.
func (q *Question) Prompt(str string) *Question {
	q.prompt = str
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
	q.clean()

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

// readSlice is the base to read slices.
func (q *Question) readSlice() (iface interface{}, err error) {
	_, err = q.newLine()
	if err != nil {
		return nil, err
	}

	term.Output.Write(readline.CRLF)
	q.isSlice = true

	switch q.schema.DataType {
	case yoda.T_SliceInt64:
		iface = make([]int64, 0)
		for {
			v, err := q.ReadInt64()
			if err != nil {
				return nil, err
			}
			if v == 0 {
				break
			}
			iface = append(iface.([]int64), v)
		}

	case yoda.T_SliceUint64:
		iface = make([]uint64, 0)
		for {
			v, err := q.ReadUint64()
			if err != nil {
				return nil, err
			}
			if v == 0 {
				break
			}
			iface = append(iface.([]uint64), v)
		}

	case yoda.T_SliceFloat64:
		iface = make([]float64, 0)
		for {
			v, err := q.ReadFloat64()
			if err != nil {
				return nil, err
			}
			if v == 0 {
				break
			}
			iface = append(iface.([]float64), v)
		}

	case yoda.T_SliceString:
		iface = make([]string, 0)
		for {
			v, err := q.ReadString()
			if err != nil {
				return nil, err
			}
			if v == "" {
				break
			}
			iface = append(iface.([]string), v)
		}

	default:
		panic("unimplemented")
	}

	q.isSlice = false
	return
}

// ReadSliceString reads multiple strings.
// You have to press Enter to finish.
func (q *Question) ReadSliceString() ([]string, error) {
	q.schema.DataType = yoda.T_SliceString

	iface, err := q.readSlice()
	if err != nil {
		return nil, err
	}
	return iface.([]string), nil
}

// == Choices
/*
// readChoice is the base to read choices.
func (q *Question) readChoice(choices interface{}) (value interface{}, err error) {
	valida := validate.NewSlice(choices, q.flag)

	// Default value
	if q.defaultVal != nil {
		valida.SetDefault(q.defaultVal)
		q.suffixPrompt = msgForDefault(q.defaultVal)
	}

	fmt.Fprintf(term.Output, "   >>> %s\r\n", valida.JoinChoices())

	value, err = q.read(q.newLine())
	q.clean()
	return
}

// ChoiceInt prints the prompt waiting to get an integer number that is in the slice.
func (q *Question) ChoiceInt(choices []int) (int, error) {
	value, err := q.readChoice(choices)
	return value.(int), err
}

// ChoiceUint prints the prompt waiting to get an unsigned number that is in the slice.
func (q *Question) ChoiceUint(choices []uint) (uint, error) {
	value, err := q.readChoice(choices)
	return value.(uint), err
}

// ChoiceFloat prints the prompt waiting to get a float number that is in the slice.
func (q *Question) ChoiceFloat(choices []float32) (float32, error) {
	value, err := q.readChoice(choices)
	return value.(float32), err
}

// ChoiceString prints the prompt waiting to get a string that is in the slice.
func (q *Question) ChoiceString(choices []string) (string, error) {
	value, err := q.readChoice(choices)
	return value.(string), err
}
*/
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

// ANSI codes to set graphic mode
const (
	setBold = "\033[1m" // Bold on
	setOff  = "\033[0m" // All attributes off
)

// The values by default are set to bold.
const lenAnsi = len(setBold) + len(setOff)

// newLine gets a line type ready to show questions.
func (q *Question) newLine() (*readline.Line, error) {
	fullPrompt := ""
	extraChars := 0

	if !q.isSlice {
		fullPrompt = q.prefixPrompt + q.prompt

		// Add spaces
		if fullPrompt[len(fullPrompt)-1] == '?' {
			fullPrompt += " "
		} else if !q.isSlice {
			fullPrompt += ": "
		}

		// Default value
		if q.schema.Bydefault != "" {
			extraChars = lenAnsi // The value by default uses ANSI characters.

			if q.schema.DataType != yoda.T_bool {
				q.suffixPrompt = fmt.Sprintf("[%s%s%s] ", setBold, q.schema.Bydefault, setOff)
			} else {
				b, err := valid.Bool(q.schema, q.schema.Bydefault)
				if err != nil {
					return nil, err
				}
				if b {
					q.suffixPrompt = fmt.Sprintf("[%s%s%s/%s] ",
						setBold, q.trueStr, setOff, q.falseStr)
				} else {
					q.suffixPrompt += fmt.Sprintf("[%s/%s%s%s] ",
						q.trueStr, setBold, q.falseStr, setOff)
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

// ExitAtCtrlC exits when it is pressed Ctrl-C, with value n.
func (q *Question) ExitAtCtrlC(n int) {
	go func() {
		select {
		case <-readline.ChanCtrlC:
			q.Restore()
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
	if err == nil {
		fmt.Printf("  answer: ")
		if _, ok := i.(string); ok {
			fmt.Printf("%q\r\n", i)
		} else {
			fmt.Printf("%v\r\n", i)
		}
	} else if err != readline.ErrCtrlD {
		fmt.Printf("%s\r\n", err)
	}
}
