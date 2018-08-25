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
			// TODO: Treat this as an error?
			continue
		}
		m[r] = g
	}

	return m, nil
}

func isRuneToDraw(r rune) bool {
	// KOREAN STANDARD SYMBOL
	if r == 0x327f {
		return true
	}
	// CIRCLE HANGLE
	if 0x3260 <= r && r <= 0x327e {
		return true
	}
	// PARENTHESIZED HANGUL
	if 0x3200 <= r && r <= 0x321f {
		return true
	}
	// HANGUL LETTER (Hangul Compatible Jamo)
	if 0x3130 <= r && r <= 0x318f {
		return true
	}
	// HANGUL SYLLABLE
	if 0xac00 <= r && r <= 0xd7af {
		return true
	}
	return false
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
	if !isRuneToDraw(r) {
		return bdf.Glyph{}, false
	}
	g, ok := glyphs[r]
	if !ok {
		return bdf.Glyph{}, false
	}
	return *g, true
}
