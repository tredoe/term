// Copyright 2010 Jonas mg
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

var (
	down2 = []byte{13, 10, 13, 10}
	up2   = []byte(fmt.Sprintf("\r%s%s", readline.CursorUp, readline.CursorUp))
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
				term.Output.Write(readline.CursorDown)
				term.Output.Write(readline.DelLine_cursorUp)
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
				term.Output.Write(readline.CursorDown)
				term.Output.Write(readline.DelLine_cursorUp)
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
				term.Output.Write(readline.CursorDown)
				term.Output.Write(readline.DelLine_cursorUp)
				return string(value), nil
			}
		}

		os.Stderr.Write(readline.DelLine_CR)
		fmt.Fprintf(os.Stderr, "%s%s", q.prefixError, "invalid choice")
		term.Output.Write(readline.CursorUp)
	}
}
