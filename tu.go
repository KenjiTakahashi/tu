// tu
// Copyright (C) 2014 Karol 'Kenji Takahashi' Woźniak
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
// TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"unicode"

	"github.com/KenjiTakahashi/tu/titlecase"

	"github.com/mitchellh/cli"
)

func prepend(args []string, arg string) []string {
	args = append(args, "")
	copy(args[1:], args[0:])
	args[0] = arg
	return args
}

type PatternPiece struct {
	Sep  string
	Name string
}

type ParseCommand struct {
	ui      cli.Ui
	wg      sync.WaitGroup
	pattern []*PatternPiece
}

func (cmd *ParseCommand) ParsePattern(in string) (out []*PatternPiece) {
	current := &PatternPiece{}
	out = append(out, current)
	simple := true
	name := false

	isalnum := func(char rune) bool {
		return unicode.IsDigit(char) || unicode.IsLetter(char)
	}
	lookahead := func(data string) {
		if len(data) > 1 && data[1] == '%' {
			current = &PatternPiece{}
			out = append(out, current)
		}
	}

	for i, char := range in {
		switch {
		case char == '%':
			if name {
				if simple {
					current = &PatternPiece{}
					out = append(out, current)
				} else {
					current.Name += "%"
				}
			} else {
				simple = true
				name = true
			}
		case char == '{' && name:
			simple = false
		case char == '}' && name:
			name = false
			lookahead(in[i:])
		default:
			if simple && !isalnum(char) {
				name = false
			}
			if name && len(current.Sep) == 0 {
				current.Name += string(char)
			} else {
				current.Sep += string(char)
				lookahead(in[i:])
			}
		}
	}

	return
}

func (cmd *ParseCommand) Process(file string) {
	defer cmd.wg.Done()
	cmd.ui.Output(fmt.Sprintf("Processing `%s`", file))

	filename := path.Base(file)
	filename = strings.TrimSuffix(filename, path.Ext(filename))

	args := []string{}
	for _, pat := range cmd.pattern {
		if len(filename) == 0 {
			break
		}
		split := []string{filename}
		if len(pat.Sep) > 0 {
			split = strings.SplitN(filename, pat.Sep, 2)
		}
		args = append(args, fmt.Sprintf("set:%s=%s", pat.Name, split[0]))
		if len(split) > 1 {
			filename = split[1]
		}
	}
	args = append(args, file)

	tagutil := exec.Command("tagutil", args...)
	if err := tagutil.Run(); err != nil {
		cmd.ui.Error(err.Error())
	}
}

func (cmd *ParseCommand) Run(args []string) int {
	if len(args) < 2 {
		cmd.ui.Output(cmd.Help())
		return 1
	}

	cmd.pattern = cmd.ParsePattern(args[0])
	files := args[1:]

	for _, file := range files {
		cmd.wg.Add(1)
		go cmd.Process(file)
	}

	cmd.wg.Wait()
	return 0
}

func (cmd *ParseCommand) Help() string {
	return "tu w PATTERN FILES..."
}

func (cmd *ParseCommand) Synopsis() string {
	return "Writes tags based on file name(s) pattern"
}

type EditCommand struct {
	ui cli.Ui
}

func (cmd *EditCommand) Run(args []string) int {
	if len(args) < 1 {
		cmd.ui.Output(cmd.Help())
		return 1
	}

	args = prepend(args, "edit")
	tagutil := exec.Command("tagutil", args...)
	tagutil.Stdin = os.Stdin
	tagutil.Stdout = os.Stdout
	tagutil.Stderr = os.Stderr
	if err := tagutil.Run(); err != nil {
		cmd.ui.Error(err.Error())
		return 1
	}

	return 0
}

func (cmd *EditCommand) Help() string {
	return "tu i FILES..."
}

func (cmd *EditCommand) Synopsis() string {
	return "Interactively edits tags using $EDITOR"
}

type TitleCaseCommand struct {
	ui cli.Ui
	wg sync.WaitGroup
}

func (cmd *TitleCaseCommand) Process(file string) {
	defer cmd.wg.Done()
	cmd.ui.Output(fmt.Sprintf("Processing `%s`", file))

	tagutil := exec.Command("tagutil", "-F", "json", file)
	out, err := tagutil.Output()
	if err != nil {
		cmd.ui.Error(err.Error())
		return
	}

	var intags []map[string]string
	json.Unmarshal(out, &intags)

	var wg sync.WaitGroup
	ch := make(chan string)
	for _, tag := range intags {
		wg.Add(1)
		go func(tag map[string]string, ch chan string) {
			defer wg.Done()

			for k, v := range tag {
				ch <- fmt.Sprintf("set:%s=%s", k, titlecase.Convert(v))
			}
		}(tag, ch)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	outtags := make([]string, 0, len(intags) + 1)
	for tag := range ch {
		outtags = append(outtags, tag)
	}

	outtags = append(outtags, file)
	tagutil = exec.Command("tagutil", outtags...)
	if err := tagutil.Run(); err != nil {
		cmd.ui.Error(err.Error())
	}
}

func (cmd *TitleCaseCommand) Run(args []string) int {
	//TODO: [-t TAGS]
	if len(args) < 1 {
		cmd.ui.Output(cmd.Help())
		return 1
	}
	for _, file := range args {
		cmd.wg.Add(1)
		go cmd.Process(file)
	}

	cmd.wg.Wait()
	return 0
}

func (cmd *TitleCaseCommand) Help() string {
	return "tu t [-t TAGS] FILES..."
}

func (cmd *TitleCaseCommand) Synopsis() string {
	return "Title Cases the Tags"
}

type RenameCommand struct {
	ui cli.Ui
}

func (cmd *RenameCommand) Run(args []string) int {
	if len(args) < 2 {
		cmd.ui.Output(cmd.Help())
		return 1
	}

	if args[0] == "-Y" {
		args[1] = fmt.Sprintf("rename:%s", args[1])
	} else {
		args[0] = fmt.Sprintf("rename:%s", args[0])
	}
	args = prepend(args, "-p")

	tagutil := exec.Command("tagutil", args...)
	if err := tagutil.Run(); err != nil {
		cmd.ui.Error(err.Error())
		return 1
	}

	return 0
}

func (cmd *RenameCommand) Help() string {
	return "tu r [-Y] PATTERN FILES..."
}

func (cmd *RenameCommand) Synopsis() string {
	return "Renames the files from tags based on pattern"
}

func main() {
	ui := &cli.ConcurrentUi{Ui: &cli.BasicUi{Writer: os.Stdout}}
	commands := map[string]cli.CommandFactory{
		"w": func() (cli.Command, error) {
			return &ParseCommand{ui: ui}, nil
		},
		"i": func() (cli.Command, error) {
			return &EditCommand{ui: ui}, nil
		},
		"t": func() (cli.Command, error) {
			return &TitleCaseCommand{ui: ui}, nil
		},
		"r": func() (cli.Command, error) {
			return &RenameCommand{ui: ui}, nil
		},
	}

	cli := &cli.CLI{
		Args:     os.Args[1:],
		Commands: commands,
		HelpFunc: cli.BasicHelpFunc("tu"),
	}

	exitCode, err := cli.Run()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(exitCode)
}
