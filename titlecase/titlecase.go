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
	rSmall = `a|an|and|as|at|but|by|en|for|if|in|of|on|or|the|to|v\.?|via|vs\.?`
	rPunct = "!\"#$%&'‘()*+,\\-./:;?@[\\]_`{|}~"

	rSmallWords   = regexp.MustCompile(fmt.Sprintf(`(?i)^(%s)$`, rSmall))
	rInlinePeriod = regexp.MustCompile(`(?i)[a-z][.][a-z]`)
	rUCElsewhere  = regexp.MustCompile(fmt.Sprintf(`[%s]*?[a-zA-Z]+[A-Z]+?`, rPunct))
	rCapFirst     = regexp.MustCompile(fmt.Sprintf(`^[%s]*?([\p{L}])`, rPunct))
	rSmallFirst   = regexp.MustCompile(fmt.Sprintf(`(?i)^([%s]*)(%s)\b`, rPunct, rSmall))
	rSmallLast    = regexp.MustCompile(fmt.Sprintf(`(?i)\b(%s)[%s]?$`, rSmall, rPunct))
	rSubPhrase    = regexp.MustCompile(fmt.Sprintf(`([:.;?!][ ])(%s)`, rSmall))
	rAposSecond   = regexp.MustCompile(`(?i)^[dol]{1}['‘]{1}[a-z]+$`)
	rAllCaps      = regexp.MustCompile(fmt.Sprintf(`^[A-Z\s\d%s]+$`, rPunct))
	rUCInitials   = regexp.MustCompile(`^(?:[A-Z]{1}\.{1}|[A-Z]{1}\.{1}[A-Z]{1})+$`)

	rLines = regexp.MustCompile(`[\r\n]+`)
	rWords = regexp.MustCompile(`[\t ]`)
)

// PreHook defines Convert per hook function signature
type PreHook func(word string, allCaps bool) (string, bool)

// PostHook defines Convert post hook function signature
type PostHook func(word string, allCaps bool) string

// Convert changes input string to conform to the NY Times Manual of Style.
//
// If pre_hook and/or post_hook arguments are not nil, they will be run,
// respectively, before and after the standard transformations.
//
// Both hook functions should return a transformed version of the string.
// Additionally, pre_hook returns true if the returned string is final,
// or false if it should still be run through standard transformations.
func Convert(text string, preHook PreHook, postHook PostHook) string {
	input := rLines.Split(text, -1)
	output := make([]string, len(input))

	for i, line := range input {
		allCaps := rAllCaps.MatchString(line)
		words := rWords.Split(line, -1)
		tcLine := make([]string, len(words))

		for j, word := range words {
			if preHook != nil {
				stop := false
				word, stop = preHook(word, allCaps)
				if stop {
					tcLine[j] = word
					continue
				}
			}

			if allCaps {
				if rUCInitials.MatchString(word) {
					tcLine[j] = word
					continue
				}
				word = strings.ToLower(word)
			}

			if rAposSecond.MatchString(word) {
				tcLine[j] = rAposSecond.ReplaceAllStringFunc(word, strings.Title)
			} else if rInlinePeriod.MatchString(word) || rUCElsewhere.MatchString(word) {
				tcLine[j] = word
			} else if rSmallWords.MatchString(word) {
				tcLine[j] = strings.ToLower(word)
			} else {
				sep := "-"
				if strings.Contains(word, "/") && !strings.Contains(word, "//") {
					sep = "/"
				}
				sepSplit := strings.Split(word, sep)
				capFirst := make([]string, len(sepSplit))
				for k, ss := range sepSplit {
					capFirst[k] = rCapFirst.ReplaceAllStringFunc(ss, strings.ToTitle)
				}
				tcLine[j] = strings.Join(capFirst, sep)
			}

			if postHook != nil {
				tcLine[j] = postHook(tcLine[j], allCaps)
			}
		}

		result := strings.Join(tcLine, " ")
		result = rSmallFirst.ReplaceAllStringFunc(result, strings.Title)
		result = rSmallLast.ReplaceAllStringFunc(result, strings.Title)
		result = rSubPhrase.ReplaceAllStringFunc(result, strings.Title)

		output[i] = result
	}

	return strings.Join(output, "\n")
}
