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

package fixed

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hajimehoshi/bitmapfont/internal/bdf"
	"github.com/hajimehoshi/bitmapfont/internal/unicode"
)

func readBDF(size int) (map[rune]*bdf.Glyph, error) {
	_, current, _, _ := runtime.Caller(1)
	dir := filepath.Dir(current)

	var filename string
	switch size {
	case 10:
		filename = "5x8.bdf"
	case 12:
		filename = "6x13.bdf"
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
	isOK := func(r rune) bool {
		if unicode.IsEuropian(r) {
			return true
		}
		if unicode.IsGeneralPunctuation(r) {
			return true
		}
		if unicode.IsSupplementalPunctuation(r) {
			return true
		}
		if 0x2100 <= r && r <= 0x218f {
			// Letterlike Symbols
			// Number Forms
			return true
		}
		if 0x2190 <= r && r <= 0x27ff {
			// [Common]
			// Arrows
			// Mathematical Operators
			// Miscellaneous Technical
			// Control Pictures
			// Optical Character Recognition
			// Enclosed Alphanumerics
			// Box Drawing
			// Block Elements
			// Geometric Shapes
			// Miscellaneous Symbols
			// Dingbats
			// Miscellaneous Mathematical Symbols-A
			// Supplemental Arrows-A
			return true
		}
		if 0x2800 <= r && r <= 0x28ff {
			// Braille Patterns
			return true
		}
		if 0x2900 <= r && r <= 0x2aff {
			// [Common]
			// Supplemental Arrows-B
			// Miscellaneous Mathematical Symbols-B
			// Supplemental Mathematical Operators
			return true
		}
		if 0x3000 <= r && r <= 0x303f {
			// CJK Symbols and Punctuation
			return true
		}
		return false
	}

	for _, g := range glyphs {
		r := rune(g.Encoding)
		if !isOK(r) {
			if unicode.IsOgham(r) {
				// Ogham glyphs in misc-fixed are too condenced. Skip this.
				continue
			}
			if unicode.IsHebrew(r) {
				continue
			}
			if unicode.IsThai(r) {
				continue
			}
			if 0x1100 <= r && r <= 0x11ff {
				// Hangul Jamo
				continue
			}
			if 0xe000 <= r && r <= 0xf8ff {
				// [Unknown]
				// Private Use Area
				continue
			}
			if 0xff00 <= r && r <= 0xffff {
				// Halfwidth and Fullwidth Forms
				// Specials
				continue
			}
			return nil, fmt.Errorf("fixed: unexpected char: 0x%x (%s)", r, string(r))
		}
		m[r] = g
	}

	return m, nil
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

func Glyph(r rune, size int) (bdf.Glyph, bool) {
	g, ok := glyphs[size][r]
	if !ok {
		return bdf.Glyph{}, false
	}
	return *g, true
}
