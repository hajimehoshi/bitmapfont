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

	"github.com/hajimehoshi/bitmapfont/internal/bdf"
	"github.com/hajimehoshi/bitmapfont/internal/uniconv"
)

func readBDF() (map[rune]*bdf.Glyph, error) {
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

	fbdf, err := os.Open(filepath.Join(dir, "gulim12.bdf"))
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
