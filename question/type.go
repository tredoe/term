// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package question

import (
	"fmt"
	"os"
	"strings"

	"github.com/kless/term"
	"github.com/kless/term/readline"
	"github.com/kless/yoda"
	"github.com/kless/yoda/valid"
)

// read is the base to read and validate input.
func (q *Question) read(typ yoda.Type) (iface interface{}, err error) {
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

		switch typ {
		case yoda.String:
			iface, err = valid.String(q.schema, input)
		case yoda.Bool:
			iface, err = valid.Bool(q.schema, input)
		case yoda.Int64:
			iface, err = valid.Int64(q.schema, input)
		case yoda.Uint64:
			iface, err = valid.Uint64(q.schema, input)
		case yoda.Float64:
			iface, err = valid.Float64(q.schema, input)
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
		return iface, nil
	}
}

// ReadBool prints the prompt waiting to get a string that represents a boolean.
func (q *Question) ReadBool() (bool, error) {
	q.isBool = true
	iface, err := q.read(yoda.Bool)
	if err != nil {
		return false, err
	}
	return iface.(bool), nil
}

// ReadInt64 prints the prompt waiting to get an integer number.
func (q *Question) ReadInt64() (int64, error) {
	iface, err := q.read(yoda.Int64)
	if err != nil {
		return 0, err
	}
	return iface.(int64), nil
}

// ReadUint64 prints the prompt waiting to get an unsigned integer number.
func (q *Question) ReadUint64() (uint64, error) {
	iface, err := q.read(yoda.Uint64)
	if err != nil {
		return 0, err
	}
	return iface.(uint64), nil
}

// ReadFloat64 prints the prompt waiting to get a floating-point number.
func (q *Question) ReadFloat64() (float64, error) {
	iface, err := q.read(yoda.Float64)
	if err != nil {
		return 0, err
	}
	return iface.(float64), nil
}

// ReadString prints the prompt waiting to get a string.
func (q *Question) ReadString() (string, error) {
	iface, err := q.read(yoda.String)
	if err != nil {
		return "", err
	}

	return iface.(string), nil
}

// == Slices

// ReadInt64Slice reads multiple int64.
// You have to press Enter twice to finish.
func (q *Question) ReadInt64Slice() (values []int64, err error) {
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

// sprintSlice returns a slice with the option by default in bold.
func (q *Question) sprintSlice(s interface{}) string {
	src := fmt.Sprintf("%v", s)
	dst := ""

	if q.schema.Bydefault != "" {
		if idx1 := strings.Index(src, q.schema.Bydefault); idx1 != -1 {
			dst = src[:idx1] + readline.ANSI_SET_BOLD

			if idx2 := strings.Index(src[idx1:], " "); idx2 != -1 {
				dst += src[idx1:idx1+idx2] + readline.ANSI_SET_OFF
				dst += src[idx1+idx2:]
			} else {
				dst += src[idx1:len(src)-1] + readline.ANSI_SET_OFF + "]"
			}
			return dst
		}
	}
	return src
}

/// readChoice is the base to read and validate input from a set of choices.
func (q *Question) readChoice(typ yoda.Type, choices interface{}) (iface interface{}, err error) {
	fmt.Fprintf(term.Output, "%s%s\r\n%s%s\r\n",
		q.prefixPrompt, q.prompt, _PREFIX_PS2, q.sprintSlice(choices),
	)

	line, err := readline.NewLine(q.term, _PREFIX_PS2, q.prefixError, 0, nil)
	if err != nil {
		return nil, err
	}

	for {
		input, err := line.Read()
		if err != nil {
			return nil, err
		}

		switch typ {
		case yoda.String:
			iface, _ = valid.String(q.schema, input)
			for _, v := range choices.([]string) {
				if v == iface.(string) {
					return iface, nil
				}
			}
		case yoda.Int:
			iface, err = valid.Int(q.schema, input)
			for _, v := range choices.([]int) {
				if v == iface.(int) {
					return iface, nil
				}
			}
		case yoda.Float64:
			iface, err = valid.Float64(q.schema, input)
			for _, v := range choices.([]float64) {
				if v == iface.(float64) {
					return iface, nil
				}
			}
		default:
			panic("unimplemented")
		}

		os.Stderr.Write(readline.DelLine_CR)
		fmt.Fprintf(os.Stderr, "%s%s", q.prefixError, "invalid choice")
		term.Output.Write(readline.CursorUp)
	}
}

// ChoiceInt prints the prompt waiting to get an int that is in the slice.
func (q *Question) ChoiceInt(choices []int) (int, error) {
	value, err := q.readChoice(yoda.Int, choices)
	if err != nil {
		return 0, err
	}
	return value.(int), nil
}

// ChoiceFloat64 prints the prompt waiting to get a float64 that is in the slice.
func (q *Question) ChoiceFloat64(choices []float64) (float64, error) {
	value, err := q.readChoice(yoda.Float64, choices)
	if err != nil {
		return 0, err
	}
	return value.(float64), nil
}

// ChoiceString prints the prompt waiting to get a string that is in the slice.
func (q *Question) ChoiceString(choices []string) (string, error) {
	value, err := q.readChoice(yoda.String, choices)
	if err != nil {
		return "", err
	}
	return value.(string), nil
}
