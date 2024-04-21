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

package mplus

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hajimehoshi/bitmapfont/v3/internal/bdf"
	"github.com/hajimehoshi/bitmapfont/v3/internal/uniconv"
)

func jisToShiftJIS(index int) int {
	upper := byte(index >> 8)
	lower := byte(index)
	upper -= 0x21
	if upper&0x1 == 0 {
		lower += 0x1f
		if lower >= 0x7f {
			lower += 0x01
		}
	} else {
		lower += 0x7e
	}

	upper >>= 1
	if upper <= 0x1e {
		upper += 0x81
	} else {
		upper += 0xc1
	}

	return (int(upper) << 8) | int(lower)
}

var (
	cp932ToRune    = map[int]rune{}
	jisx0201ToRune = map[int]rune{}
)

func init() {
	_, current, _, _ := runtime.Caller(0)
	dir := filepath.Dir(current)

	f, err := os.Open(filepath.Join(dir, "CP932.TXT"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	c, err := uniconv.Parse(f, "\t")
	if err != nil {
		panic(err)
	}

	cp932ToRune = c
}

func init() {
	_, current, _, _ := runtime.Caller(0)
	dir := filepath.Dir(current)

	f, err := os.Open(filepath.Join(dir, "JIS0201.TXT"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	c, err := uniconv.Parse(f, "\t")
	if err != nil {
		panic(err)
	}

	jisx0201ToRune = c
}

func readBDF(size int) (map[rune]*bdf.Glyph, error) {
	_, current, _, _ := runtime.Caller(1)
	dir := filepath.Dir(current)

	// Latin alphabets
	prefix := "f"
	postfix := ""
	if size == 10 {
		prefix = "j"
		postfix = "-iso-W4"
	}
	fe, err := os.Open(filepath.Join(dir, fmt.Sprintf("mplus_%s%dr%s.bdf", prefix, size, postfix)))
	if err != nil {
		return nil, err
	}
	defer fe.Close()

	glyphse, err := bdf.Parse(fe)
	if err != nil {
		return nil, err
	}

	// Hankaku-kana glyphs
	prefix = "f"
	if size == 10 {
		prefix = "h"
	}
	fej, err := os.Open(filepath.Join(dir, fmt.Sprintf("mplus_%s%dr_jisx0201.bdf", prefix, size)))
	if err != nil {
		return nil, err
	}
	defer fej.Close()

	glyphsej, err := bdf.Parse(fej)
	if err != nil {
		return nil, err
	}

	// Other characters (e.g., Kanji)
	fj, err := os.Open(filepath.Join(dir, fmt.Sprintf("mplus_j%dr.bdf", size)))
	if err != nil {
		return nil, err
	}
	defer fj.Close()

	glyphsj, err := bdf.Parse(fj)
	if err != nil {
		return nil, err
	}

	m := map[rune]*bdf.Glyph{}

	for _, g := range glyphse {
		m[rune(g.Encoding)] = g
	}

	for _, g := range glyphsej {
		r, ok := jisx0201ToRune[g.Encoding]
		if !ok {
			if g.Encoding < 0x20 {
				// Control chars
				continue
			}
			if g.Encoding == 0x7f {
				// DELETE
				continue
			}
			if g.Encoding == 0xa0 {
				// NO-BREAK SPACE
				continue
			}
			return nil, fmt.Errorf("mplus: invalid char code 0x%x as JIS X 0201", g.Encoding)
		}
		if 0xff60 <= r && r <= 0xffdf {
			m[r] = g
		}
	}

	for _, g := range glyphsj {
		r := rune(0)
		if g.Encoding == 0x2474 {
			// HIRAGANA LETTER VU
			r = 0x3094
		} else {
			s := jisToShiftJIS(g.Encoding)
			ok := false
			r, ok = cp932ToRune[s]
			if !ok {
				return nil, fmt.Errorf("mplus: invalid char code 0x%x (Shift_JIS: 0x%x)", g.Encoding, s)
			}
		}

		if _, ok := m[r]; ok {
			// Prefer f12r/j10r-iso for Latin glyphs.
			continue
		}

		if size == 12 && !isValidGlyph(r, g) {
			return nil, fmt.Errorf("mplus: invalid glyph for rune 0x%x", r)
		}

		if size == 10 {
			g.ShiftY = 1
		}

		m[r] = g
	}

	const (
		uniWaveDash       = 0x301c
		uniFullwidthTilde = 0xff5e
	)

	if _, ok := m[uniFullwidthTilde]; !ok {
		return nil, fmt.Errorf("mplus: FULLWIDTH TILDE (0x%x) not found", uniFullwidthTilde)
	}
	m[uniWaveDash] = m[uniFullwidthTilde]

	// https://unicode.org/Public/MAPPINGS/VENDORS/MICSFT/WINDOWS/CP932.TXT
	// According to CP932.TXT, a fullwidth dash is mapped to U+2015, not U+2014. This is a known issue.
	// As their glyphs are very similar, use the same glyphs as U+2015 for U+2014.
	const (
		uniEmDash        = 0x2014
		uniHorizontalBar = 0x2015
	)
	m[uniEmDash] = m[uniHorizontalBar]

	return m, nil
}

func isValidGlyph(r rune, g *bdf.Glyph) bool {
	// Box Drawing
	if 0x2500 <= r && r <= 0x257f {
		return true
	}

	// Check the edges
	left := false
	right := false
	for i := 0; i < g.Height; i++ {
		if _, _, _, a := g.At(0, i).RGBA(); a != 0 {
			left = true
			break
		}
	}
	for i := 0; i < g.Height; i++ {
		if _, _, _, a := g.At(g.Width-1, i).RGBA(); a != 0 {
			right = true
			break
		}
	}
	if !left && right {
		return false
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
