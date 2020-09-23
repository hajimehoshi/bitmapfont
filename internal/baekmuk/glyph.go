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

	"github.com/hajimehoshi/bitmapfont/v2/internal/bdf"
	"github.com/hajimehoshi/bitmapfont/v2/internal/uniconv"
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

func Glyph(r rune, size int) (bdf.Glyph, bool) {
	g, ok := glyphs[size][r]
	if !ok {
		return bdf.Glyph{}, false
	}
	return *g, true
}
