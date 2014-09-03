// tu
// Copyright (C) 2014 Karol 'Kenji Takahashi' Woźniak
// Python version by Stuart Colville http://muffinresearch.co.uk
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

package titlecase

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ConvertTests = []struct {
	input    string
	expected string
}{
	{"word/word",
		"Word/Word"},
	{"dance with me/let’s face the music and dance",
		"Dance With Me/Let’s Face the Music and Dance"},
	{"34th 3rd 2nd",
		"34th 3rd 2nd"},
	{"Q&A with steve jobs: 'that's what happens in technology'",
		"Q&A With Steve Jobs: 'That's What Happens in Technology'"},
	{"What is AT&T's problem?",
		"What Is AT&T's Problem?"},
	{"Apple deal with AT&T falls through",
		"Apple Deal With AT&T Falls Through"},
	{"this v that",
		"This v That"},
	{"this v. that",
		"This v. That"},
	{"this vs that",
		"This vs That"},
	{"this vs. that",
		"This vs. That"},
	{"The SEC's Apple probe: what you need to know",
		"The SEC's Apple Probe: What You Need to Know"},
	{"'by the Way, small word at the start but within quotes.'",
		"'By the Way, Small Word at the Start but Within Quotes.'"},
	{"Small word at end is nothing to be afraid of",
		"Small Word at End Is Nothing to Be Afraid Of"},
	{"Starting Sub-Phrase With a Small Word: a Trick, Perhaps?",
		"Starting Sub-Phrase With a Small Word: A Trick, Perhaps?"},
	{"Sub-Phrase With a Small Word in Quotes: 'a Trick, Perhaps?'",
		"Sub-Phrase With a Small Word in Quotes: 'A Trick, Perhaps?'"},
	{`sub-phrase with a small word in quotes: "a trick, perhaps?"`,
		`Sub-Phrase With a Small Word in Quotes: "A Trick, Perhaps?"`},
	{`"Nothing to Be Afraid of?"`,
		`"Nothing to Be Afraid Of?"`},
	{`"Nothing to be Afraid Of?"`,
		`"Nothing to Be Afraid Of?"`},
	{`a thing`,
		`A Thing`},
	{"2lmc Spool: 'gruber on OmniFocus and vapo(u)rware'",
		"2lmc Spool: 'Gruber on OmniFocus and Vapo(u)rware'"},
	{`this is just an example.com`,
		`This Is Just an example.com`},
	{`this is something listed on del.icio.us`,
		`This Is Something Listed on del.icio.us`},
	{`iTunes should be unmolested`,
		`iTunes Should Be Unmolested`},
	{`reading between the lines of steve jobs’s ‘thoughts on music’`,
		`Reading Between the Lines of Steve Jobs’s ‘Thoughts on Music’`},
	{`seriously, ‘repair permissions’ is voodoo`,
		`Seriously, ‘Repair Permissions’ Is Voodoo`},
	{`generalissimo francisco franco: still dead; kieren McCarthy: still a jackass`,
		`Generalissimo Francisco Franco: Still Dead; Kieren McCarthy: Still a Jackass`},
	{"O'Reilly should be untouched",
		"O'Reilly Should Be Untouched"},
	{"my name is o'reilly",
		"My Name Is O'Reilly"},
	{"WASHINGTON, D.C. SHOULD BE FIXED BUT MIGHT BE A PROBLEM",
		"Washington, D.C. Should Be Fixed but Might Be a Problem"},
	{"THIS IS ALL CAPS AND SHOULD BE ADDRESSED",
		"This Is All Caps and Should Be Addressed"},
	{"Mr McTavish went to MacDonalds",
		"Mr McTavish Went to MacDonalds"},
	{"this shouldn't\nget mangled",
		"This Shouldn't\nGet Mangled"},
	{"this is http://foo.com",
		"This Is http://foo.com"},
	{"mac mc MAC MC machine",
		"Mac Mc MAC MC Machine"},
	{"FOO BAR 5TH ST",
		"Foo Bar 5th St"},
	{"foo bar 5th st",
		"Foo Bar 5th St"},
	{"śledź",
		"Śledź"},
}

func TestConvert(t *testing.T) {
	for i, tt := range ConvertTests {
		actual := Convert(tt.input, nil, nil)

		assert.Equal(t, tt.expected, actual, fmt.Sprintf("%d", i))
	}
}

var ConvertPreHookTests = []struct {
	input    string
	expected string
}{
	{"TEST", "mock"},
	{"test", "Mock"},
}

func TestConvert_PreHook(t *testing.T) {
	hook := func(word string, all_caps bool) (string, bool) {
		return "mock", all_caps
	}
	for i, tt := range ConvertPreHookTests {
		actual := Convert(tt.input, hook, nil)

		assert.Equal(t, tt.expected, actual, fmt.Sprintf("%d", i))
	}
}

func TestConvert_PostHook(t *testing.T) {
	hook := func(word string, all_caps bool) string {
		return "mock"
	}

	actual := Convert("test", nil, hook)

	assert.Equal(t, "mock", actual)
}
