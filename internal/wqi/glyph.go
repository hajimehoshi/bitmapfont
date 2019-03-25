// Copyright 2019 Hajime Hoshi
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

	"github.com/hajimehoshi/bitmapfont/internal/bdf"
)

func Includes(r rune) bool {
	if 0x4e00 <= r && r <= 0x9fff {
		return true
	}
	return false
}

func readBDF() (map[rune]*bdf.Glyph, error) {
	_, current, _, _ := runtime.Caller(1)
	dir := filepath.Dir(current)

	fbdf, err := os.Open(filepath.Join(dir, "wenquanyi_9pt.bdf"))
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
		if Includes(r) {
			m[r] = g
		}
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
