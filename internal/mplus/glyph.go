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

	"github.com/hajimehoshi/bitmapfont/internal/bdf"
	"github.com/hajimehoshi/bitmapfont/internal/uniconv"
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

func readBDF() (map[rune]*bdf.Glyph, error) {
	_, current, _, _ := runtime.Caller(1)
	dir := filepath.Dir(current)

	fe, err := os.Open(filepath.Join(dir, "mplus_f12r.bdf"))
	if err != nil {
		return nil, err
	}
	defer fe.Close()

	fej, err := os.Open(filepath.Join(dir, "mplus_f12r_jisx0201.bdf"))
	if err != nil {
		return nil, err
	}
	defer fej.Close()

	fj, err := os.Open(filepath.Join(dir, "mplus_j12r.bdf"))
	if err != nil {
		return nil, err
	}
	defer fj.Close()

	m := map[rune]*bdf.Glyph{}

	glyphse, err := bdf.Parse(fe)
	if err != nil {
		return nil, err
	}

	glyphsej, err := bdf.Parse(fej)
	if err != nil {
		return nil, err
	}

	glyphsj, err := bdf.Parse(fj)
	if err != nil {
		return nil, err
	}

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
		if _, ok := m[r]; ok {
			continue
		}
		m[r] = g
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
			// Prefer f12r for Latin glyphs.
			continue
		}
		if 0x2500 <= r && r <= 0x257f {
			// Box Drawing
			// M+ defines a part of box drawing glyphs.
			// For consistency, use other font's glyphs instead.
			continue
		}

		if !isValidGlyph(r, g) {
			return nil, fmt.Errorf("mplus: invalid glyph for rune 0x%x", r)
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
