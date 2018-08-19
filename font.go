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

// Package bitmapfont offers a font.Face value of the bitmap font.
//
//   * U+AC00–D7AF: Bitmap Hangul glyphs by Hajime Hoshi (Public Domain)
//   * Others:      [M+ Bitmap Font](http://mplus-fonts.osdn.jp/mplus-bitmap-fonts/) (M+ Bitmap Fonts License)
package bitmapfont

import (
	"bytes"
	"compress/gzip"
	"image"
	"io/ioutil"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var imageData *image.RGBA

const (
	imageWidth  = 3072
	imageHeight = 4096
)

func init() {
	s, err := gzip.NewReader(bytes.NewReader(compressedMplusRGBA))
	if err != nil {
		panic(err)
	}
	defer s.Close()

	pix, err := ioutil.ReadAll(s)
	if err != nil {
		panic(err)
	}

	imageData = &image.RGBA{
		Pix:    pix,
		Stride: 4 * imageWidth,
		Rect:   image.Rect(0, 0, imageWidth, imageHeight),
	}
}

const (
	charHalfWidth = 6
	charFullWidth = 12
	charHeight    = 16
	charXNum      = 256
	charYNum      = 256

	dotX = 4
	dotY = 12
)

func runeWidth(r rune) int {
	if r < 0x100 {
		return charHalfWidth
	}
	return charFullWidth
}

// Gothic12r is a font.Face of M+ bitmap font (M+ gothic 12r).
var Gothic12r font.Face = &mplusFont{}

type mplusFont struct {
	scale int
}

func (m *mplusFont) Close() error {
	return nil
}

func (m *mplusFont) Glyph(dot fixed.Point26_6, r rune) (dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
	if r >= 0x10000 {
		return
	}

	// Use 'Wave Dash' glyph for 'Fullwidth Tilde'.
	if r == 0xff5e {
		r = 0x301c
	}

	rw := runeWidth(r)
	dx := dot.X.Floor() - dotX
	dy := dot.Y.Floor() - dotY
	dr = image.Rect(dx, dy, dx+rw, dy+charHeight)

	mx := (int(r) % charXNum) * charFullWidth
	my := (int(r) / charXNum) * charHeight
	mask = imageData.SubImage(image.Rect(mx, my, mx+rw, my+charHeight))
	maskp = image.Pt(mx, my)
	advance = fixed.I(runeWidth(r))
	ok = true
	return
}

func (m *mplusFont) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	if r >= 0x10000 {
		return
	}
	bounds = fixed.R(-dotX, -dotY, -dotX+runeWidth(r), -dotY+charHeight)
	advance = fixed.I(runeWidth(r))
	ok = true
	return
}

func (m *mplusFont) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	if r >= 0x10000 {
		return
	}
	advance = fixed.I(runeWidth(r))
	ok = true
	return
}

func (m *mplusFont) Kern(r0, r1 rune) fixed.Int26_6 {
	return 0

}
func (m *mplusFont) Metrics() font.Metrics {
	return font.Metrics{
		Height:  fixed.I(charHeight),
		Ascent:  fixed.I(dotY),
		Descent: fixed.I(charHeight - dotY),
	}
}
