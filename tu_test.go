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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ParsePatternTests = []struct {
	input    string
	expected []*PatternPiece
}{
	{"%n1", []*PatternPiece{
		&PatternPiece{Name: "n1"},
	}},
	{"%{n%_1b}", []*PatternPiece{
		&PatternPiece{Name: "n%_1b"},
	}},
	{"_-_", []*PatternPiece{
		&PatternPiece{Sep: "_-_"},
	}},
	{"%n1%n2", []*PatternPiece{
		&PatternPiece{Name: "n1"},
		&PatternPiece{Name: "n2"},
	}},
	{"%{n%_1b}%{n%_2c}", []*PatternPiece{
		&PatternPiece{Name: "n%_1b"},
		&PatternPiece{Name: "n%_2c"},
	}},
	{"%n1_%n2", []*PatternPiece{
		&PatternPiece{Sep: "_", Name: "n1"},
		&PatternPiece{Name: "n2"},
	}},
	{"%{n%_1b}_%{n%_2c}", []*PatternPiece{
		&PatternPiece{Sep: "_", Name: "n%_1b"},
		&PatternPiece{Name: "n%_2c"},
	}},
	{"%n1_%{n%_2c}", []*PatternPiece{
		&PatternPiece{Sep: "_", Name: "n1"},
		&PatternPiece{Name: "n%_2c"},
	}},
	{"%{n%_1b}_%n2", []*PatternPiece{
		&PatternPiece{Sep: "_", Name: "n%_1b"},
		&PatternPiece{Name: "n2"},
	}},
	{"_%n2", []*PatternPiece{
		&PatternPiece{Sep: "_"},
		&PatternPiece{Name: "n2"},
	}},
	{"_%{n%_2c}", []*PatternPiece{
		&PatternPiece{Sep: "_"},
		&PatternPiece{Name: "n%_2c"},
	}},
	{"%n1_", []*PatternPiece{
		&PatternPiece{Sep: "_", Name: "n1"},
	}},
	{"%{n%_1b}_", []*PatternPiece{
		&PatternPiece{Sep: "_", Name: "n%_1b"},
	}},
	{"%{artist}_-_%album - %{tracknumber}@%title.flac", []*PatternPiece{
		&PatternPiece{Sep: "_-_", Name: "artist"},
		&PatternPiece{Sep: " - ", Name: "album"},
		&PatternPiece{Sep: "@", Name: "tracknumber"},
		&PatternPiece{Sep: ".flac", Name: "title"},
	}},
}

func TestParsePattern(t *testing.T) {
	cmd := ParseCommand{}

	for i, tt := range ParsePatternTests {
		actual := cmd.ParsePattern(tt.input)

		assert.Equal(
			t, tt.expected, actual,
			fmt.Sprintf("%d: %q => %q != %q", i, tt.input, actual, tt.expected),
		)
	}
}
