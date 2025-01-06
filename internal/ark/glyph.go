// Copyright 2024 Hajime Hoshi
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

package ark

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/hajimehoshi/bitmapfont/v3/internal/bdf"
	"github.com/hajimehoshi/bitmapfont/v3/internal/unicode"
)

// The current version is:
// https://github.com/TakWolf/ark-pixel-font/releases/tag/2025.01.06
// (ark-pixel-font-12px-monospaced-bdf-*.zip)

func readBDF(filename string) (map[rune]*bdf.Glyph, error) {
	_, current, _, _ := runtime.Caller(1)
	dir := filepath.Dir(current)

	f, err := os.Open(filepath.Join(dir, filename))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	glyphs, err := bdf.Parse(f)
	if err != nil {
		return nil, err
	}

	m := map[rune]*bdf.Glyph{}

	for _, g := range glyphs {
		g.ShiftY -= 1
		m[rune(g.Encoding)] = g
	}

	return m, nil
}

var (
	cnGlyphs map[rune]*bdf.Glyph
	trGlyphs map[rune]*bdf.Glyph
)

func init() {
	g, err := readBDF("ark-pixel-12px-monospaced-zh_cn.bdf")
	if err != nil {
		panic(err)
	}
	cnGlyphs = g

	g, err = readBDF("ark-pixel-12px-monospaced-zh_tr.bdf")
	if err != nil {
		panic(err)
	}
	trGlyphs = g
}

func Glyph(r rune, simplified bool) (*bdf.Glyph, bool) {
	if !unicode.IsCJKUnifiedIdeograph(r) {
		return nil, false
	}

	if simplified {
		if g, ok := cnGlyphs[r]; ok {
			return g, true
		}
		if g, ok := trGlyphs[r]; ok {
			return g, true
		}
	} else {
		if g, ok := trGlyphs[r]; ok {
			return g, true
		}
		if g, ok := cnGlyphs[r]; ok {
			return g, true
		}
	}

	return nil, false
}
