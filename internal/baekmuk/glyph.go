// Copyright 2018 Hajime Hoshi
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

package baekmuk

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hajimehoshi/bitmapfont/v4/internal/bdf"
	"github.com/hajimehoshi/bitmapfont/v4/internal/unicode"
	"github.com/hajimehoshi/bitmapfont/v4/internal/uniconv"
)

func readBDF(size int) (map[rune]*bdf.Glyph, error) {
	_, current, _, _ := runtime.Caller(1)
	dir := filepath.Dir(current)

	funi, err := os.Open(filepath.Join(dir, "KSX1001.TXT"))
	if err != nil {
		return nil, err
	}
	defer funi.Close()

	c, err := uniconv.Parse(funi, "  ")
	if err != nil {
		return nil, err
	}

	var filename string
	switch size {
	case 10:
		filename = "gulim10.bdf"
	case 12:
		filename = "gulim12.bdf"
	default:
		panic("not reached")
	}

	fbdf, err := os.Open(filepath.Join(dir, filename))
	if err != nil {
		return nil, err
	}
	defer fbdf.Close()

	m := map[rune]*bdf.Glyph{}

	glyphs, err := bdf.Parse(fbdf)
	if err != nil {
		return nil, err
	}
	for _, g := range glyphs {
		r, ok := c[g.Encoding]
		if !ok {
			if g.Encoding == 0x2266 {
				// EURO SIGN
				continue
			}
			if g.Encoding == 0x2267 {
				// REGISTERED SIGN
				continue
			}
			return nil, fmt.Errorf("baekmuk: invalid char code 0x%x", g.Encoding)
		}

		if needsShift(r, g) {
			g.ShiftX = -1
		}

		m[r] = g
	}

	if size == 10 {
		return m, nil
	}

	// Complement the glyphs for 12px.
	// See "Johab 8/4/4 Font Layout" in https://unifoundry.com/hangul/hangul-generation.html

	/*for r := 0xac00; r <= 0xd7a3; r++ {
		l, v, t := variationsFromHangulSyllable(rune(r))
		println(l, v, t)
	}*/

	for lIdx := range unicode.HangulLeadingConsonantCount {
		l := unicode.HangulLeadingConsonantCodePointBase + rune(lIdx)
		for lVar := range 8 {
			r := typicalGlyphForLeadingConsonantVariation(lVar, l)
			if _, ok := m[r]; !ok {
				panic(fmt.Sprintf("baekmuk: missing leading consonant variation: %c, l: %c lVar: %d", r, l, lVar))
			}
		}
	}

	/*for tIdx := range unicode.HangulTailingConsonantCount {
		for vIdx := range unicode.HangulVowelCount {
			for lIdx := range unicode.HangulLeadingConsonantCount {
				l := unicode.HangulLeadingConsonantCodePointBase + rune(lIdx)
				v := unicode.HangulVowelCodePointBase + rune(vIdx)
				var t rune
				if tIdx != 0 {
					t = unicode.HangulTailingConsonantCodePointBase + rune(tIdx)
				}
				s := unicode.ComposeHangulSyllable(l, v, t)
				fmt.Printf("%c\n", s)
			}
		}
	}*/

	return m, nil
}

func needsShift(r rune, g *bdf.Glyph) bool {
	// Basically glyphs needs to be shifted in X axis by 1px, but
	// there are some exceptions:

	// Box Drawing
	if 0x2500 <= r && r <= 0x257f {
		return false
	}

	// Check the left edge
	for i := 0; i < g.Height; i++ {
		if _, _, _, a := g.At(0, i).RGBA(); a != 0 {
			return false
		}
	}

	return true
}

var glyphs map[int]map[rune]*bdf.Glyph

func init() {
	glyphs10r, err := readBDF(10)
	if err != nil {
		panic(err)
	}
	glyphs12r, err := readBDF(12)
	if err != nil {
		panic(err)
	}

	glyphs = map[int]map[rune]*bdf.Glyph{
		10: glyphs10r,
		12: glyphs12r,
	}
}

func Glyph(r rune, size int) (*bdf.Glyph, bool) {
	g, ok := glyphs[size][r]
	if !ok {
		return nil, false
	}
	return g, true
}

func typicalGlyphForLeadingConsonantVariation(lVar int, l rune) rune {
	switch lVar {
	case 0:
		return unicode.ComposeHangulSyllable(l, 'ᅡ', 0)
	case 1:
		return unicode.ComposeHangulSyllable(l, 'ᅩ', 0)
	case 2:
		return unicode.ComposeHangulSyllable(l, 'ᅮ', 0)
	case 3:
		return unicode.ComposeHangulSyllable(l, 'ᅬ', 0)
	case 4:
		switch l {
		case 'ᄄ':
			return unicode.ComposeHangulSyllable(l, 'ᅰ', 0)
		case 'ᄈ':
			// No glyph for the variation 4 in KS X 1001. Use the variation 3 instead.
			return unicode.ComposeHangulSyllable(l, 'ᅬ', 0)
		default:
			return unicode.ComposeHangulSyllable(l, 'ᅯ', 0)
		}
	case 5:
		return unicode.ComposeHangulSyllable(l, 'ᅡ', 'ᆨ')
	case 6:
		return unicode.ComposeHangulSyllable(l, 'ᅩ', 'ᆨ')
	case 7:
		return unicode.ComposeHangulSyllable(l, 'ᅱ', 'ᆫ')
	default:
		panic(fmt.Sprintf("baekmuk: invalid leading consonant variation: %d", lVar))
	}
}

func typicalGlyphForVowelVariation(vVar int, v rune) rune {
	switch vVar {
	case 0:
		return unicode.ComposeHangulSyllable('ᄀ', v, 0)
	case 1:
		return unicode.ComposeHangulSyllable('ᄂ', v, 0)
	case 2:
		return unicode.ComposeHangulSyllable('ᄀ', v, 'ᆨ')
	case 3:
		return unicode.ComposeHangulSyllable('ᄂ', v, 'ᆨ')
	default:
		panic(fmt.Sprintf("baekmuk: invalid vowel variation: %d", vVar))
	}
}

func typicalGlyphForTailingConsonantVariation(tVar int, t rune) rune {
	switch tVar {
	case -1:
		return 0
	case 0:
		return unicode.ComposeHangulSyllable('ᄀ', 'ᅡ', t)
	case 1:
		return unicode.ComposeHangulSyllable('ᄀ', 'ᅥ', t)
	case 2:
		return unicode.ComposeHangulSyllable('ᄀ', 'ᅢ', t)
	case 3:
		return unicode.ComposeHangulSyllable('ᄀ', 'ᅩ', t)
	default:
		panic(fmt.Sprintf("baekmuk: invalid tailing consonant variation: %d", tVar))
	}
}

func variationsFromHangulSyllable(r rune) (int, int, int) {
	// See "Johab 8/4/4 Font Layout" in https://unifoundry.com/hangul/hangul-generation.html
	l, v, t := unicode.DecomposeHangulSyllable(r)

	if t == 0 {
		var lVar, vVar int
		switch v {
		case 'ᅡ', 'ᅢ', 'ᅣ', 'ᅤ', 'ᅥ', 'ᅦ', 'ᅧ', 'ᅨ', 'ᅵ':
			lVar = 0
		case 'ᅩ', 'ᅭ', 'ᅳ':
			lVar = 1
		case 'ᅮ', 'ᅲ':
			lVar = 2
		case 'ᅪ', 'ᅫ', 'ᅬ', 'ᅴ':
			lVar = 3
		case 'ᅯ', 'ᅰ', 'ᅱ':
			lVar = 4
		default:
			panic(fmt.Sprintf("baekmuk: invalid vowel: %c", v))
		}
		if l == 'ᄀ' || l == 'ᄏ' {
			vVar = 0
		} else {
			vVar = 1
		}
		return lVar, vVar, -1
	}

	var lVar, vVar, tVar int
	switch v {
	case 'ᅡ', 'ᅢ', 'ᅣ', 'ᅤ', 'ᅥ', 'ᅦ', 'ᅧ', 'ᅨ', 'ᅵ':
		lVar = 5
	case 'ᅩ', 'ᅭ', 'ᅮ', 'ᅲ', 'ᅳ':
		lVar = 6
	case 'ᅪ', 'ᅫ', 'ᅬ', 'ᅯ', 'ᅰ', 'ᅱ', 'ᅴ':
		lVar = 7
	default:
		panic(fmt.Sprintf("baekmuk: invalid vowel: %c", v))
	}
	if l == 'ᄀ' || l == 'ᄏ' {
		vVar = 2
	} else {
		vVar = 3
	}
	switch v {
	case 'ᅡ', 'ᅣ', 'ᅪ':
		tVar = 0
	case 'ᅥ', 'ᅧ', 'ᅬ', 'ᅯ', 'ᅱ', 'ᅴ', 'ᅵ':
		tVar = 1
	case 'ᅢ', 'ᅤ', 'ᅦ', 'ᅨ', 'ᅫ', 'ᅰ':
		tVar = 2
	case 'ᅩ', 'ᅭ', 'ᅮ', 'ᅲ', 'ᅳ':
		tVar = 3
	default:
		panic(fmt.Sprintf("baekmuk: invalid vowel: %c", v))
	}

	return lVar, vVar, tVar
}
