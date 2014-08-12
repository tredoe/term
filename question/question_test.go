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

func TestQuest(t *testing.T) {
	//var ans interface{}
	var err error

	q := NewCustom(nil, "", "", "oui", "non")
	defer func() {
		if err = q.Restore(); err != nil {
			t.Error(err)
		}
	}()

	fmt.Print("\n== Questions\n\n")

	if *NeedUser {
		q.ExitAtCtrlC(1)
	} else {
		Int := "-1"
		Uint := "1"
		Float := "1.1"
		String := "foo"
		False := "false"
		True := "true"

		auto := map[int][]string{
			0:  []string{String, "", "Oui"},
			1:  []string{Float, "", "R. C."},
			2:  []string{""},
			3:  []string{String, Int},
			4:  []string{True, Int, ""},
			5:  []string{String, Float},
			6:  []string{String, ""},
			7:  []string{False},

			8:  []string{String, ""},
			9:  []string{String, Uint, "10"},
			10: []string{Uint, "10"},
			11: []string{String, ""},
			12: []string{String, ""},
			13: []string{"photo", "cryp", ""},
		}

		go func() {
			for i := 0; i <= 13; i++ {
				for _, v := range auto[i] {
					time.Sleep(time.Duration(*Time) * time.Second)
					// Remember that the terminal is in raw mode.
					fmt.Fprintf(pw, "%s%s", v, readline.CRLF)
				}
			}
		}()
	}

	q.Prompt("0. Are you french?").Mod(valid.M_Required)
	PrintAnswer(q.ReadBool())

	q.Prompt("1. What is your name?").Mod(valid.M_Required)//.Min(4)
	PrintAnswer(q.ReadString())

	q.Prompt("2. What color is your hair?").Default("brown")
	PrintAnswer(q.ReadString())

	q.Prompt("3. What temperature is there?").Default("-2")
	PrintAnswer(q.ReadInt64())

	q.Prompt("4. How old are you?").Default("16")
	PrintAnswer(q.ReadUint64())

	q.Prompt("5. How tall are you?").Default("1.74")
	PrintAnswer(q.ReadFloat64())

	q.Prompt("6. Do you watch television?").Default("true")
	PrintAnswer(q.ReadBool())

	q.Prompt("7. Do you read books?").Default("false")
	PrintAnswer(q.ReadBool())
/*
	color := []string{"red", "blue", "black"}
	q.Prompt("8. What is your favourite color?").Default("blue")
	ans, err = q.ChoiceString(color)
	PrintAnswer(ans, err)

	q.Prompt("9. Another favourite color?")
	ans, err = q.ChoiceString(color)
	PrintAnswer(ans, err)

	q.Prompt("10. Choose number").Default(uint(3))
	ans, err = q.ChoiceUint([]uint{1, 3, 5})
	PrintAnswer(ans, err)

	q.Prompt("11. Email").Default("ja@contac.me")
	ans, err = q.SetBasicEmail().Read()
	PrintAnswer(ans, err)

	q.Prompt("12. Contact").Default("https://foo.com")
	ans, err = q.SetBasicEmail().AddRegexp("contact address", `^http`).Read()
	PrintAnswer(ans, err)

	q.Prompt("13. Hobby")
	ans, err = q.ReadSliceString()
	PrintAnswer(ans, err)*/
}
