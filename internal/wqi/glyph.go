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

package wqi

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/hajimehoshi/bitmapfont/v3/internal/bdf"
)

func readBDF() (map[rune]*bdf.Glyph, error) {
	_, current, _, _ := runtime.Caller(1)
	dir := filepath.Dir(current)

	// Download wqy-bitmapsong-bdf-1.0.0-RC1.tar.gz from
	// https://sourceforge.net/projects/wqy/files/wqy-bitmapfont/1.0.0-RC1/
	// This file is not in the repository due to the license issue.
	f, err := os.Open(filepath.Join(dir, "wenquanyi_9pt.bdf"))
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

func isCJKUnifiedIdeograph(r rune) bool {
	// CJK Unified Ideographs
	if 0x4E00 <= r && r <= 0x9FFF {
		return true
	}
	// CJK Unified Ideographs Extension A
	if 0x3400 <= r && r <= 0x4DBF {
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

func Glyph(r rune) (*bdf.Glyph, bool) {
	if !isCJKUnifiedIdeograph(r) {
		return nil, false
	}

	g, ok := glyphs[r]
	if !ok {
		return nil, false
	}
	return g, true
}
