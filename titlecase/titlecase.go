// tu
// Copyright (C) 2014 Karol 'Kenji Takahashi' Woźniak
// Original Perl version by: John Gruber http://daringfireball.net/ 10 May 2008
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

// Package titlecase provides a helper function for transforming an arbitrary
// text into a capitalized version, as described by NY Times Manual of Style.
//
// Additional hooks can be supplied to modify the standard behaviour,
// see Convert documentation for further details.
//
// Original Perl version by: John Gruber http://daringfireball.net/ 10 May 2008
// Python version by Stuart Colville http://muffinresearch.co.uk
package titlecase

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	SMALL = `a|an|and|as|at|but|by|en|for|if|in|of|on|or|the|to|v\.?|via|vs\.?`
	PUNCT = "!\"#$%&'‘()*+,\\-./:;?@[\\]_`{|}~"

	SMALL_WORDS   = regexp.MustCompile(fmt.Sprintf(`(?i)^(%s)$`, SMALL))
	INLINE_PERIOD = regexp.MustCompile(`(?i)[a-z][.][a-z]`)
	UC_ELSEWHERE  = regexp.MustCompile(fmt.Sprintf(`[%s]*?[a-zA-Z]+[A-Z]+?`, PUNCT))
	CAPFIRST      = regexp.MustCompile(fmt.Sprintf(`^[%s]*?([\p{L}])`, PUNCT))
	SMALL_FIRST   = regexp.MustCompile(fmt.Sprintf(`(?i)^([%s]*)(%s)\b`, PUNCT, SMALL))
	SMALL_LAST    = regexp.MustCompile(fmt.Sprintf(`(?i)\b(%s)[%s]?$`, SMALL, PUNCT))
	SUBPHRASE     = regexp.MustCompile(fmt.Sprintf(`([:.;?!][ ])(%s)`, SMALL))
	APOS_SECOND   = regexp.MustCompile(`(?i)^[dol]{1}['‘]{1}[a-z]+$`)
	ALL_CAPS      = regexp.MustCompile(fmt.Sprintf(`^[A-Z\s\d%s]+$`, PUNCT))
	UC_INITIALS   = regexp.MustCompile(`^(?:[A-Z]{1}\.{1}|[A-Z]{1}\.{1}[A-Z]{1})+$`)
	MAC_MC        = regexp.MustCompile(`^([Mm]c)(\w+)`)

	LINES = regexp.MustCompile(`[\r\n]+`)
	WORDS = regexp.MustCompile(`[\t ]`)
)

// PreHook defines Convert per hook function signature
type PreHook func(word string, all_caps bool) (string, bool)

// PostHook defines Convert post hook function signature
type PostHook func(word string, all_caps bool) string

// Convert changes input string to conform to the NY Times Manual of Style.
//
// If pre_hook and/or post_hook arguments are not nil, they will be run,
// respectively, before and after the standard transformations.
//
// Both hook functions should return a transformed version of the string.
// Additionally, pre_hook returns true if the returned string is final,
// or false if it should still be run through standard transformations.
func Convert(text string, pre_hook PreHook, post_hook PostHook) string {
	input := LINES.Split(text, -1)
	output := make([]string, len(input))

	for i, line := range input {
		all_caps := ALL_CAPS.MatchString(line)
		words := WORDS.Split(line, -1)
		tc_line := make([]string, len(words))

		for j, word := range words {
			if pre_hook != nil {
				stop := false
				word, stop = pre_hook(word, all_caps)
				if stop {
					tc_line[j] = word
					continue
				}
			}

			if all_caps {
				if UC_INITIALS.MatchString(word) {
					tc_line[j] = word
					continue
				}
				word = strings.ToLower(word)
			}

			if APOS_SECOND.MatchString(word) {
				tc_line[j] = APOS_SECOND.ReplaceAllStringFunc(word, strings.Title)
			} else if INLINE_PERIOD.MatchString(word) || UC_ELSEWHERE.MatchString(word) {
				tc_line[j] = word
			} else if SMALL_WORDS.MatchString(word) {
				tc_line[j] = strings.ToLower(word)
			} else if match := MAC_MC.FindStringSubmatch(word); match != nil {
				tc_line[j] = fmt.Sprintf(
					"%s%s", strings.Title(match[1]), strings.Title(match[2]),
				)
			} else {
				sep := "-"
				if strings.Contains(word, "/") && !strings.Contains(word, "//") {
					sep = "/"
				}
				sep_split := strings.Split(word, sep)
				cap_first := make([]string, len(sep_split))
				for k, ss := range sep_split {
					cap_first[k] = CAPFIRST.ReplaceAllStringFunc(ss, strings.ToTitle)
				}
				tc_line[j] = strings.Join(cap_first, sep)
			}

			if post_hook != nil {
				tc_line[j] = post_hook(tc_line[j], all_caps)
			}
		}

		result := strings.Join(tc_line, " ")
		result = SMALL_FIRST.ReplaceAllStringFunc(result, strings.Title)
		result = SMALL_LAST.ReplaceAllStringFunc(result, strings.Title)
		result = SUBPHRASE.ReplaceAllStringFunc(result, strings.Title)

		output[i] = result
	}

	return strings.Join(output, "\n")
}
