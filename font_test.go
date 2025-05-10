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

package bitmapfont_test

import (
	"testing"

	"github.com/hajimehoshi/bitmapfont/v4"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func BenchmarkLazyFace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := bitmapfont.NewLazyFace("data/face_ja.bin", false)
		if _, _, _, _, ok := l.Glyph(fixed.P(0, 0), 'あ'); !ok {
			b.Fatal("Glyph failed")
		}
	}
}

func TestWidth(t *testing.T) {
	testCaeses := []struct {
		str string
		w   fixed.Int26_6
	}{
		{
			str: "",
			w:   0,
		},
		{
			str: "a",
			w:   fixed.I(6),
		},
		{
			str: "ä",
			w:   fixed.I(6),
		},
		{
			str: "あ",
			w:   fixed.I(12),
		},
		{
			str: "亜",
			w:   fixed.I(12),
		},
		{
			str: "ｱ",
			w:   fixed.I(6),
		},
		{
			str: "ｯ",
			w:   fixed.I(6),
		},
	}
	for _, tc := range testCaeses {
		advance := font.MeasureString(bitmapfont.Face, tc.str)
		if got, want := advance, tc.w; got != want {
			t.Errorf("width for %q: got: %v, want: %v", tc.str, got, want)
		}
	}
}
