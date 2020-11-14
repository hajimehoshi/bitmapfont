// Copyright 2020 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bitmapfont

import (
	"sort"
	"unicode"

	"golang.org/x/text/language"
	"golang.org/x/text/unicode/bidi"
)

type arabicLetterPresentationForms struct {
	isolated rune
	initial  rune
	medial   rune
	final    rune
}

// TODO: Implement the table for the other languages like Kurdish.
var arabicLetterTable = map[rune]arabicLetterPresentationForms{
	// ARABIC LETTER HAMZA
	0x0621: {isolated: 0xFE80, initial: 0, medial: 0, final: 0},
	// ARABIC LETTER ALEF WITH MADDA ABOVE
	0x0622: {isolated: 0xFE81, initial: 0, medial: 0, final: 0xFE82},
	// ARABIC LETTER ALEF WITH HAMZA ABOVE
	0x0623: {isolated: 0xFE83, initial: 0, medial: 0, final: 0xFE84},
	// ARABIC LETTER WAW WITH HAMZA ABOVE
	0x0624: {isolated: 0xFE85, initial: 0, medial: 0, final: 0xFE86},
	// ARABIC LETTER ALEF WITH HAMZA BELOW
	0x0625: {isolated: 0xFE87, initial: 0, medial: 0, final: 0xFE88},
	// ARABIC LETTER YEH WITH HAMZA ABOVE
	0x0626: {isolated: 0xFE89, initial: 0xFE8B, medial: 0xFE8C, final: 0xFE8A},
	// ARABIC LETTER ALEF
	0x0627: {isolated: 0xFE8D, initial: 0, medial: 0, final: 0xFE8E},
	// ARABIC LETTER BEH
	0x0628: {isolated: 0xFE8F, initial: 0xFE91, medial: 0xFE92, final: 0xFE90},
	// ARABIC LETTER TEH MARBUTA
	0x0629: {isolated: 0xFE93, initial: 0, medial: 0, final: 0xFE94},
	// ARABIC LETTER TEH
	0x062A: {isolated: 0xFE95, initial: 0xFE97, medial: 0xFE98, final: 0xFE96},
	// ARABIC LETTER THEH
	0x062B: {isolated: 0xFE99, initial: 0xFE9B, medial: 0xFE9C, final: 0xFE9A},
	// ARABIC LETTER JEEM
	0x062C: {isolated: 0xFE9D, initial: 0xFE9F, medial: 0xFEA0, final: 0xFE9E},
	// ARABIC LETTER HAH
	0x062D: {isolated: 0xFEA1, initial: 0xFEA3, medial: 0xFEA4, final: 0xFEA2},
	// ARABIC LETTER KHAH
	0x062E: {isolated: 0xFEA5, initial: 0xFEA7, medial: 0xFEA8, final: 0xFEA6},
	// ARABIC LETTER DAL
	0x062F: {isolated: 0xFEA9, initial: 0, medial: 0, final: 0xFEAA},
	// ARABIC LETTER THAL
	0x0630: {isolated: 0xFEAB, initial: 0, medial: 0, final: 0xFEAC},
	// ARABIC LETTER REH
	0x0631: {isolated: 0xFEAD, initial: 0, medial: 0, final: 0xFEAE},
	// ARABIC LETTER ZAIN
	0x0632: {isolated: 0xFEAF, initial: 0, medial: 0, final: 0xFEB0},
	// ARABIC LETTER SEEN
	0x0633: {isolated: 0xFEB1, initial: 0xFEB3, medial: 0xFEB4, final: 0xFEB2},
	// ARABIC LETTER SHEEN
	0x0634: {isolated: 0xFEB5, initial: 0xFEB7, medial: 0xFEB8, final: 0xFEB6},
	// ARABIC LETTER SAD
	0x0635: {isolated: 0xFEB9, initial: 0xFEBB, medial: 0xFEBC, final: 0xFEBA},
	// ARABIC LETTER DAD
	0x0636: {isolated: 0xFEBD, initial: 0xFEBF, medial: 0xFEC0, final: 0xFEBE},
	// ARABIC LETTER TAH
	0x0637: {isolated: 0xFEC1, initial: 0xFEC3, medial: 0xFEC4, final: 0xFEC2},
	// ARABIC LETTER ZAH
	0x0638: {isolated: 0xFEC5, initial: 0xFEC7, medial: 0xFEC8, final: 0xFEC6},
	// ARABIC LETTER AIN
	0x0639: {isolated: 0xFEC9, initial: 0xFECB, medial: 0xFECC, final: 0xFECA},
	// ARABIC LETTER GHAIN
	0x063A: {isolated: 0xFECD, initial: 0xFECF, medial: 0xFED0, final: 0xFECE},
	// ARABIC TATWEEL
	0x0640: {isolated: 0x0640, initial: 0x0640, medial: 0x0640, final: 0x0640},
	// ARABIC LETTER FEH
	0x0641: {isolated: 0xFED1, initial: 0xFED3, medial: 0xFED4, final: 0xFED2},
	// ARABIC LETTER QAF
	0x0642: {isolated: 0xFED5, initial: 0xFED7, medial: 0xFED8, final: 0xFED6},
	// ARABIC LETTER KAF
	0x0643: {isolated: 0xFED9, initial: 0xFEDB, medial: 0xFEDC, final: 0xFEDA},
	// ARABIC LETTER LAM
	0x0644: {isolated: 0xFEDD, initial: 0xFEDF, medial: 0xFEE0, final: 0xFEDE},
	// ARABIC LETTER MEEM
	0x0645: {isolated: 0xFEE1, initial: 0xFEE3, medial: 0xFEE4, final: 0xFEE2},
	// ARABIC LETTER NOON
	0x0646: {isolated: 0xFEE5, initial: 0xFEE7, medial: 0xFEE8, final: 0xFEE6},
	// ARABIC LETTER HEH
	0x0647: {isolated: 0xFEE9, initial: 0xFEEB, medial: 0xFEEC, final: 0xFEEA},
	// ARABIC LETTER WAW
	0x0648: {isolated: 0xFEED, initial: 0, medial: 0, final: 0xFEEE},
	// ARABIC LETTER ALEF MAKSURA
	0x0649: {isolated: 0xFEEF, initial: 0, medial: 0, final: 0xFEF0},
	// ARABIC LETTER YEH
	0x064A: {isolated: 0xFEF1, initial: 0xFEF3, medial: 0xFEF4, final: 0xFEF2},
	// ARABIC LETTER ALEF WASLA
	0x0671: {isolated: 0xFB50, initial: 0, medial: 0, final: 0xFB51},
	// ARABIC LETTER U WITH HAMZA ABOVE
	0x0677: {isolated: 0xFBDD, initial: 0, medial: 0, final: 0},
	// ARABIC LETTER TTEH
	0x0679: {isolated: 0xFB66, initial: 0xFB68, medial: 0xFB69, final: 0xFB67},
	// ARABIC LETTER TTEHEH
	0x067A: {isolated: 0xFB5E, initial: 0xFB60, medial: 0xFB61, final: 0xFB5F},
	// ARABIC LETTER BEEH
	0x067B: {isolated: 0xFB52, initial: 0xFB54, medial: 0xFB55, final: 0xFB53},
	// ARABIC LETTER PEH
	0x067E: {isolated: 0xFB56, initial: 0xFB58, medial: 0xFB59, final: 0xFB57},
	// ARABIC LETTER TEHEH
	0x067F: {isolated: 0xFB62, initial: 0xFB64, medial: 0xFB65, final: 0xFB63},
	// ARABIC LETTER BEHEH
	0x0680: {isolated: 0xFB5A, initial: 0xFB5C, medial: 0xFB5D, final: 0xFB5B},
	// ARABIC LETTER NYEH
	0x0683: {isolated: 0xFB76, initial: 0xFB78, medial: 0xFB79, final: 0xFB77},
	// ARABIC LETTER DYEH
	0x0684: {isolated: 0xFB72, initial: 0xFB74, medial: 0xFB75, final: 0xFB73},
	// ARABIC LETTER TCHEH
	0x0686: {isolated: 0xFB7A, initial: 0xFB7C, medial: 0xFB7D, final: 0xFB7B},
	// ARABIC LETTER TCHEHEH
	0x0687: {isolated: 0xFB7E, initial: 0xFB80, medial: 0xFB81, final: 0xFB7F},
	// ARABIC LETTER DDAL
	0x0688: {isolated: 0xFB88, initial: 0, medial: 0, final: 0xFB89},
	// ARABIC LETTER DAHAL
	0x068C: {isolated: 0xFB84, initial: 0, medial: 0, final: 0xFB85},
	// ARABIC LETTER DDAHAL
	0x068D: {isolated: 0xFB82, initial: 0, medial: 0, final: 0xFB83},
	// ARABIC LETTER DUL
	0x068E: {isolated: 0xFB86, initial: 0, medial: 0, final: 0xFB87},
	// ARABIC LETTER RREH
	0x0691: {isolated: 0xFB8C, initial: 0, medial: 0, final: 0xFB8D},
	// ARABIC LETTER JEH
	0x0698: {isolated: 0xFB8A, initial: 0, medial: 0, final: 0xFB8B},
	// ARABIC LETTER VEH
	0x06A4: {isolated: 0xFB6A, initial: 0xFB6C, medial: 0xFB6D, final: 0xFB6B},
	// ARABIC LETTER PEHEH
	0x06A6: {isolated: 0xFB6E, initial: 0xFB70, medial: 0xFB71, final: 0xFB6F},
	// ARABIC LETTER KEHEH
	0x06A9: {isolated: 0xFB8E, initial: 0xFB90, medial: 0xFB91, final: 0xFB8F},
	// ARABIC LETTER NG
	0x06AD: {isolated: 0xFBD3, initial: 0xFBD5, medial: 0xFBD6, final: 0xFBD4},
	// ARABIC LETTER GAF
	0x06AF: {isolated: 0xFB92, initial: 0xFB94, medial: 0xFB95, final: 0xFB93},
	// ARABIC LETTER NGOEH
	0x06B1: {isolated: 0xFB9A, initial: 0xFB9C, medial: 0xFB9D, final: 0xFB9B},
	// ARABIC LETTER GUEH
	0x06B3: {isolated: 0xFB96, initial: 0xFB98, medial: 0xFB99, final: 0xFB97},
	// ARABIC LETTER NOON GHUNNA
	0x06BA: {isolated: 0xFB9E, initial: 0, medial: 0, final: 0xFB9F},
	// ARABIC LETTER RNOON
	0x06BB: {isolated: 0xFBA0, initial: 0xFBA2, medial: 0xFBA3, final: 0xFBA1},
	// ARABIC LETTER HEH DOACHASHMEE
	0x06BE: {isolated: 0xFBAA, initial: 0xFBAC, medial: 0xFBAD, final: 0xFBAB},
	// ARABIC LETTER HEH WITH YEH ABOVE
	0x06C0: {isolated: 0xFBA4, initial: 0, medial: 0, final: 0xFBA5},
	// ARABIC LETTER HEH GOAL
	0x06C1: {isolated: 0xFBA6, initial: 0xFBA8, medial: 0xFBA9, final: 0xFBA7},
	// ARABIC LETTER KIRGHIZ OE
	0x06C5: {isolated: 0xFBE0, initial: 0, medial: 0, final: 0xFBE1},
	// ARABIC LETTER OE
	0x06C6: {isolated: 0xFBD9, initial: 0, medial: 0, final: 0xFBDA},
	// ARABIC LETTER U
	0x06C7: {isolated: 0xFBD7, initial: 0, medial: 0, final: 0xFBD8},
	// ARABIC LETTER YU
	0x06C8: {isolated: 0xFBDB, initial: 0, medial: 0, final: 0xFBDC},
	// ARABIC LETTER KIRGHIZ YU
	0x06C9: {isolated: 0xFBE2, initial: 0, medial: 0, final: 0xFBE3},
	// ARABIC LETTER VE
	0x06CB: {isolated: 0xFBDE, initial: 0, medial: 0, final: 0xFBDF},
	// ARABIC LETTER FARSI YEH
	0x06CC: {isolated: 0xFBFC, initial: 0xFBFE, medial: 0xFBFF, final: 0xFBFD},
	// ARABIC LETTER E
	0x06D0: {isolated: 0xFBE4, initial: 0xFBE6, medial: 0xFBE7, final: 0xFBE5},
	// ARABIC LETTER YEH BARREE
	0x06D2: {isolated: 0xFBAE, initial: 0, medial: 0, final: 0xFBAF},
	// ARABIC LETTER YEH BARREE WITH HAMZA ABOVE
	0x06D3: {isolated: 0xFBB0, initial: 0, medial: 0, final: 0xFBB1},
	// ZERO WIDTH JOINER
	0x200D: {isolated: 0x200D, initial: 0x200D, medial: 0x200D, final: 0x200D},
}

type Direction int

const (
	DirectionLeftToRight Direction = iota
	DirectionRightToLeft
)

type arabicForm int

const (
	arabicFormNeutral arabicForm = iota
	arabicFormIsolated
	arabicFormInitial
	arabicFormMedial
	arabicFormFinal
)

type runeWithForm struct {
	r    rune
	form arabicForm
}

// PresentationForms returns runes as presentation forms in order to render it easily.
//
// PresentationForms mainly converts RTL texts into LTR glyphs for presentation.
// The result can be passed to e.g., golang.org/x/image.Drawer's DrawString.
// PresentationForms should work with texts whose directions are mixed, but the strict bidirectional algorithm [1]
// is not implemented yet.
//
// lang represents a language that is a hint to compose the representation forms.
// lang is not used in the implementation yet, but might be used in the future.
//
// [1] https://unicode.org/reports/tr9/
func PresentationForms(input string, defaultDirection Direction, lang language.Tag) string {
	canConnectBefore := func(r rune) bool {
		f, ok := arabicLetterTable[r]
		if !ok {
			return false
		}
		return f.final != 0 || f.medial != 0
	}
	canConnectAfter := func(r rune) bool {
		f, ok := arabicLetterTable[r]
		if !ok {
			return false
		}
		return f.initial != 0 || f.medial != 0
	}
	canConnectBeforeAndAfter := func(r rune) bool {
		f, ok := arabicLetterTable[r]
		if !ok {
			return false
		}
		return f.medial != 0
	}

	// TODO: Treat ZWS correctly

	var runeWithForms []runeWithForm
	for _, r := range input {
		if _, ok := arabicLetterTable[r]; !ok {
			runeWithForms = append(runeWithForms, runeWithForm{r: r})
			continue
		}

		var prev runeWithForm
		prevIdx := len(runeWithForms) - 1
		for ; prevIdx >= 0; prevIdx-- {
			prev = runeWithForms[prevIdx]
			// Nonspacing mark should not be involved in the connection algorithm.
			if unicode.Is(unicode.Mn, prev.r) {
				continue
			}
			break
		}

		if prevIdx == -1 || prev.form == arabicFormNeutral {
			runeWithForms = append(runeWithForms, runeWithForm{r: r, form: arabicFormIsolated})
			continue
		}
		if !canConnectBefore(r) {
			runeWithForms = append(runeWithForms, runeWithForm{r: r, form: arabicFormIsolated})
			continue
		}
		if !canConnectAfter(prev.r) {
			runeWithForms = append(runeWithForms, runeWithForm{r: r, form: arabicFormIsolated})
			continue
		}
		if prev.form == arabicFormFinal && !canConnectBeforeAndAfter(prev.r) {
			runeWithForms = append(runeWithForms, runeWithForm{r: r, form: arabicFormIsolated})
			continue
		}
		if prev.form == arabicFormIsolated {
			runeWithForms[prevIdx].form = arabicFormInitial
			runeWithForms = append(runeWithForms, runeWithForm{r: r, form: arabicFormFinal})
			continue
		}
		runeWithForms[prevIdx].form = arabicFormMedial
		runeWithForms = append(runeWithForms, runeWithForm{r: r, form: arabicFormFinal})
	}

	var runes []rune
	for i := 0; i < len(runeWithForms); i++ {
		if i < len(runeWithForms)-1 {
			if r, ok := processLigature(runeWithForms[i], runeWithForms[i+1]); ok {
				i++
				runes = append(runes, r)
				continue
			}
		}

		rf := runeWithForms[i]
		var r rune
		switch rf.form {
		case arabicFormNeutral:
			r = rf.r
		case arabicFormIsolated:
			r = arabicLetterTable[rf.r].isolated
		case arabicFormInitial:
			r = arabicLetterTable[rf.r].initial
		case arabicFormMedial:
			r = arabicLetterTable[rf.r].medial
		case arabicFormFinal:
			r = arabicLetterTable[rf.r].final
		}
		runes = append(runes, r)
	}

	// TODO: Implement the strict bidi algorithm.
	// https://unicode.org/reports/tr9/

	type runesWithDir struct {
		idx  int
		rs   []rune
		dir0 Direction
		dir1 Direction
	}
	var runesWithDirs []*runesWithDir
	for _, r := range runes {
		p, _ := bidi.LookupRune(r)
		var dir0, dir1 Direction
		var prev *runesWithDir
		if len(runesWithDirs) > 0 {
			prev = runesWithDirs[len(runesWithDirs)-1]
		}
		if prev != nil {
			dir0 = prev.dir0
		} else {
			dir0 = defaultDirection
		}
		switch p.Class() {
		case bidi.L:
			dir0 = DirectionLeftToRight
			dir1 = dir0
		case bidi.R, bidi.AL:
			dir0 = DirectionRightToLeft
			dir1 = dir0
		case bidi.EN, bidi.ES, bidi.ET, bidi.AN, bidi.CS, bidi.BN:
			// Weak
			dir1 = DirectionLeftToRight
		case bidi.B, bidi.S, bidi.WS, bidi.ON:
			// Neutral
			dir1 = dir0
		case bidi.NSM:
			// Non-spacing mark is weak, but treat this as nutral.
			// TODO: Is this correct?
			dir1 = dir0
		default:
			// TODO: Implement control characters
			dir1 = dir0
		}

		if prev != nil && prev.dir0 == dir0 && prev.dir1 == dir1 {
			prev.rs = append(prev.rs, r)
			continue
		}
		runesWithDirs = append(runesWithDirs, &runesWithDir{
			idx:  len(runesWithDirs),
			rs:   []rune{r},
			dir0: dir0,
			dir1: dir1,
		})
	}

	if defaultDirection == DirectionRightToLeft {
		sort.Slice(runesWithDirs, func(i, j int) bool {
			return runesWithDirs[i].idx > runesWithDirs[j].idx
		})
	}

	var result []rune
	for _, rd := range runesWithDirs {
		switch rd.dir1 {
		case DirectionLeftToRight:
			for _, r := range rd.rs {
				result = append(result, r)
			}
		case DirectionRightToLeft:
			var marks []rune
			for i := range rd.rs {
				r := rd.rs[len(rd.rs)-i-1]

				// Place the mark character in the logically correct position.
				// Accumulate marks until the current character is not a mark.
				if unicode.Is(unicode.Mn, r) {
					marks = append(marks, r)
					continue
				}

				result = append(result, r)
				for i := range marks {
					result = append(result, marks[len(marks)-i-1])
				}
				marks = nil
			}
			for i := range marks {
				result = append(result, marks[len(marks)-i-1])
			}
		}
	}

	return string(result)
}

// processLigature returns a ligature for the runes r1 and r2 when possible.
// processLigature processes only part of Arabic ligatures for this package's glyphs.
func processLigature(r1, r2 runeWithForm) (rune, bool) {
	const (
		arabicLetterLam                = 0x0644
		arabicLetterAlefWithMaddaAbove = 0x0622
		arabicLetterAlefWithHamzaAbove = 0x0623
		arabicLetterAlefWithHamzaBelow = 0x0625
		arabicLetterAlef               = 0x0627
	)

	if r1.r != arabicLetterLam {
		return 0, false
	}
	switch r2.r {
	case arabicLetterAlefWithMaddaAbove:
		switch r1.form {
		case arabicFormInitial:
			return 0xFEF5, true
		case arabicFormMedial:
			return 0xFEF6, true
		}
	case arabicLetterAlefWithHamzaAbove:
		switch r1.form {
		case arabicFormInitial:
			return 0xFEF7, true
		case arabicFormMedial:
			return 0xFEF8, true
		}
	case arabicLetterAlefWithHamzaBelow:
		switch r1.form {
		case arabicFormInitial:
			return 0xFEF9, true
		case arabicFormMedial:
			return 0xFEFA, true
		}
	case arabicLetterAlef:
		switch r1.form {
		case arabicFormInitial:
			return 0xFEFB, true
		case arabicFormMedial:
			return 0xFEFC, true
		}
	}
	return 0, false
}
