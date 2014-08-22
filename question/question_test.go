// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package question

import (
	"flag"
	"fmt"
	"io"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/kless/term"
	"github.com/kless/term/readline"
	"github.com/kless/yoda/valid"
)

var (
	NeedUser = flag.Bool("user", false, "need user interaction")
	Time     = flag.Uint("t", 2, "time in seconds to wait to write the answers in automatic mode")

	// Interactive mode:
	// It is needed a fifo for os.Stdin (os.Stderr is used by 'go test') since
	// it is connected to the TTY. So, you write go to the term.
	pr *io.PipeReader
	pw *io.PipeWriter
)

func init() {
	flag.Parse()
	term.InputFD = syscall.Stderr

	if *NeedUser {
		term.Input = os.Stderr
	} else {
		pr, pw = io.Pipe()
		term.Input = pr
	}
}

func TestQuestion(t *testing.T) {
	var err error

	q := NewCustom(nil, "", "", "oui", "non")
	defer func() {
		if err = q.Restore(); err != nil {
			t.Error(err)
		}
	}()

	fmt.Print("\n== Questions\n\n")

	if !*NeedUser {
		Int := "-1"
		Uint := "1"
		Float := "1.1"
		String := "foo"
		False := "false"
		True := "true"

		auto := map[int][]string{
			1:  []string{Float, "", "R. C."},
			2:  []string{""},
			3:  []string{String, Int},
			4:  []string{True, Int, ""},
			5:  []string{String, Float},
			6:  []string{String, "", "Oui"},
			7:  []string{String, ""},
			8:  []string{False},
			9:  []string{String, ""},
			10: []string{String, "blue"},
			11: []string{Uint},
			12: []string{String, ""},
			13: []string{String, ""},
			14: []string{"photo", "cryp", ""},
			15: []string{String},
		}

		go func() {
			for i := 1; i <= 15; i++ {
				for _, v := range auto[i] {
					time.Sleep(time.Duration(*Time) * time.Second)
					// Remember that the terminal is in raw mode.
					fmt.Fprintf(pw, "%s%s", v, readline.CRLF)
				}
			}
		}()
	}

	q.Prompt("1. What is your name?").
		Check(valid.C_Required | valid.C_StrictString).Min(4)
	PrintAnswer(q.ReadString())

	q.Prompt("2. What color is your hair?").Default("brown")
	PrintAnswer(q.ReadString())

	q.Prompt("3. What temperature is there?").Default("-2")
	PrintAnswer(q.ReadInt64())

	q.Prompt("4. How old are you?").Default("16")
	PrintAnswer(q.ReadUint64())

	q.Prompt("5. How tall are you?").Default("1.74")
	PrintAnswer(q.ReadFloat64())

	q.Prompt("6. Are you french?").Check(valid.C_Required)
	PrintAnswer(q.ReadBool())

	q.Prompt("7. Do you watch television?").Default("true")
	PrintAnswer(q.ReadBool())

	q.Prompt("8. Do you read books?").Default("false")
	PrintAnswer(q.ReadBool())

	color := []string{"red", "blue", "black"}

	q.Prompt("9. What is your favourite color?").Default("blue")
	PrintAnswer(q.ChoiceString(color))

	q.Prompt("10. Another favourite color?")
	PrintAnswer(q.ChoiceString(color))

	q.Prompt("11. Choose number").Default("3")
	PrintAnswer(q.ChoiceInt([]int{1, 3, 5}))

	q.Prompt("12. Email").Default("ja@contac.me").Check(valid.C_StrictString)
	PrintAnswer(q.ReadEmail())

	q.Prompt("13. Web").Default("https://foo.com").Check(valid.C_DNS)
	PrintAnswer(q.ReadURL())

	q.Prompt("14. Hobby").Check(valid.C_StrictString)
	PrintAnswer(q.ReadStringSlice())

	q.Prompt("15. A film").Default("Terminator")
	PrintAnswer(q.ReadString())
}
