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

	"github.com/hajimehoshi/bitmapfont/v3"
	"golang.org/x/image/math/fixed"
)

func BenchmarkLazyFace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := bitmapfont.NewLazyFace("data/face_ja.bin")
		if _, _, _, _, ok := l.Glyph(fixed.P(0, 0), 'あ'); !ok {
			b.Fatal("Glyph failed")
		}
	}
}
