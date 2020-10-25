// Copyright 2020 Hajime Hoshi
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

package arabic

import (
	"image"
	"image/color"
	_ "image/png"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hajimehoshi/bitmapfont/v2/internal/fixed"
)

const (
	glyphFullWidth = 12
	glyphHalfWidth = 6
	glyphHeight    = 12
)

var (
	images = map[rune]image.Image{}
)

func subImage(img image.Image, x, y int) image.Image {
	const yoffset = 2

	var wide bool

	result := image.NewAlpha(image.Rect(0, 0, glyphFullWidth, glyphHeight+yoffset))
	for j := 0; j < glyphHeight; j++ {
		for i := 0; i < glyphFullWidth; i++ {
			r, _, _, _ := img.At(x+i, y+j).RGBA()
			var c color.Alpha
			if r == 0 {
				if i >= glyphHalfWidth {
					wide = true
				}
				c = color.Alpha{0xff}
			}
			result.SetAlpha(i, j+yoffset, c)
		}
	}
	if !wide {
		return result.SubImage(image.Rect(0, 0, glyphHalfWidth, glyphHeight+yoffset))
	}
	return result
}

func init() {
	_, current, _, _ := runtime.Caller(0)
	dir := filepath.Dir(current)

	f, err := os.Open(filepath.Join(dir, "0600-06ff.png"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	r := rune(0x0600)
	for j := 0; j < 16; j++ {
		for i := 0; i < 16; i++ {
			x := i * glyphFullWidth
			y := j * glyphHeight
			images[r] = subImage(img, x, y)
			r++
		}
	}
}

func init() {
	_, current, _, _ := runtime.Caller(0)
	dir := filepath.Dir(current)

	f, err := os.Open(filepath.Join(dir, "fe70-feff.png"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	r := rune(0xfe70)
	for j := 0; j < 16; j++ {
		for i := 0; i < 16; i++ {
			x := i * glyphFullWidth
			y := j * glyphHeight
			images[r] = subImage(img, x, y)
			r++
		}
	}
}

func Glyph(r rune) (image.Image, bool) {
	const (
		arabicComma        = '\u060c'
		arabicSemicolon    = '\u061b'
		arabicQuestionMark = '\u061f'
	)

	switch r {
	case arabicComma:
		if g, ok := fixed.Glyph(',', 12); ok {
			return rotatedImage{
				img:    &g,
				shiftY: 7,
			}, true
		}
	case arabicSemicolon:
		if g, ok := fixed.Glyph(';', 12); ok {
			return rotatedImage{
				img:    &g,
				shiftY: 2,
			}, true
		}
	case arabicQuestionMark:
		if g, ok := fixed.Glyph('?', 12); ok {
			return mirroredImage{&g}, true
		}
	default:
		if img, ok := images[r]; ok {
			return img, true
		}
	}
	return nil, false
}
