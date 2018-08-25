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

func readBDF() (map[rune]*bdf.Glyph, error) {
	_, current, _, _ := runtime.Caller(1)
	dir := filepath.Dir(current)

	fbdf, err := os.Open(filepath.Join(dir, "6x13.bdf"))
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
		r := rune(g.Encoding)
		if unicode.IsOgham(r) {
			// Ogham glyphs in misc-fixed are too condenced. Skip this.
			continue
		}
		if !unicode.IsEuropian(r) {
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

			if 0x2100 <= r && r <= 0x218f {
				// Letterlike Symbols
				// Number Forms
				continue
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
				continue
			}
			if 0x2800 <= r && r <= 0x28ff {
				// Braille Patterns
				continue
			}
			if 0x2900 <= r && r <= 0x2aff {
				// [Common]
				// Supplemental Arrows-B
				// Miscellaneous Mathematical Symbols-B
				// Supplemental Mathematical Operators
				continue
			}
			if 0x3000 <= r && r <= 0x303f {
				// CJK Symbols and Punctuation
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

var glyphs map[rune]*bdf.Glyph

func init() {
	var err error
	glyphs, err = readBDF()
	if err != nil {
		panic(err)
	}
}

func Glyph(r rune) (bdf.Glyph, bool) {
	g, ok := glyphs[r]
	if !ok {
		return bdf.Glyph{}, false
	}
	return *g, true
}
