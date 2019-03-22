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

// Package bitmapfont offers a font.Face value of some bitmap fonts.
//
//   * [Baekmuk Gulim](http://kldp.net/baekmuk/) (Baekmuk License)
//   * [misc-fixed](https://www.cl.cam.ac.uk/~mgk25/ucs-fonts.html) (Public Domain)
//   * [M+ Bitmap Font](http://mplus-fonts.osdn.jp/mplus-bitmap-fonts/) (M+ Bitmap Fonts License)
package bitmapfont

import (
	"bytes"
	"compress/gzip"
	"image"
	"image/color"
	"io/ioutil"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/hajimehoshi/bitmapfont/internal/unicode"
)

const (
	imageWidth  = 3072
	imageHeight = 4096
)

type binaryImage struct {
	bits   []byte
	width  int
	height int
	bounds image.Rectangle
}

func (b *binaryImage) At(i, j int) color.Color {
	if i < b.bounds.Min.X || j < b.bounds.Min.Y || i >= b.bounds.Max.X || j >= b.bounds.Max.Y {
		return color.Alpha{0}
	}
	idx := b.width*j + i
	if (b.bits[idx/8]>>uint(7-idx%8))&1 != 0 {
		return color.Alpha{0xff}
	}
	return color.Alpha{0}
}

func (b *binaryImage) ColorModel() color.Model {
	return color.AlphaModel
}

func (b *binaryImage) Bounds() image.Rectangle {
	return b.bounds
}

func (b *binaryImage) SubImage(r image.Rectangle) image.Image {
	bounds := r.Intersect(b.bounds)
	if bounds.Empty() {
		return &binaryImage{}
	}
	return &binaryImage{
		bits:   b.bits,
		width:  b.width,
		height: b.height,
		bounds: bounds,
	}
}

func init() {
	s, err := gzip.NewReader(bytes.NewReader(compressedFontAlpha))
	if err != nil {
		panic(err)
	}
	defer s.Close()

	bits, err := ioutil.ReadAll(s)
	if err != nil {
		panic(err)
	}

	Gothic12r = &bitmapFont{
		imageData: &binaryImage{
			bits:   bits,
			width:  imageWidth,
			height: imageHeight,
			bounds: image.Rect(0, 0, imageWidth, imageHeight),
		},
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
	// TODO: This condition depends on the fact that Europian glyphs are from misc-fixed.
	// Refactor this.
	if unicode.IsEuropian(r) {
		return charHalfWidth
	}
	if 0xff61 <= r && r <= 0xffdc {
		return charHalfWidth
	}
	if 0xffe8 <= r && r <= 0xffee {
		return charHalfWidth
	}
	return charFullWidth
}

// Gothic12r is a font.Face of the bitmap font (12px regular).
var Gothic12r font.Face

type bitmapFont struct {
	scale     int
	imageData *binaryImage
}

func (m *bitmapFont) Close() error {
	return nil
}

func (m *bitmapFont) Glyph(dot fixed.Point26_6, r rune) (dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
	if r >= 0x10000 {
		return
	}

	rw := runeWidth(r)
	dx := dot.X.Floor() - dotX
	dy := dot.Y.Floor() - dotY
	dr = image.Rect(dx, dy, dx+rw, dy+charHeight)

	mx := (int(r) % charXNum) * charFullWidth
	my := (int(r) / charXNum) * charHeight
	mask = m.imageData.SubImage(image.Rect(mx, my, mx+rw, my+charHeight))
	maskp = image.Pt(mx, my)
	advance = fixed.I(runeWidth(r))
	ok = true
	return
}

func (m *bitmapFont) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	if r >= 0x10000 {
		return
	}
	bounds = fixed.R(-dotX, -dotY, -dotX+runeWidth(r), -dotY+charHeight)
	advance = fixed.I(runeWidth(r))
	ok = true
	return
}

func (m *bitmapFont) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	if r >= 0x10000 {
		return
	}
	advance = fixed.I(runeWidth(r))
	ok = true
	return
}

func (m *bitmapFont) Kern(r0, r1 rune) fixed.Int26_6 {
	return 0

}
func (m *bitmapFont) Metrics() font.Metrics {
	return font.Metrics{
		Height:  fixed.I(charHeight),
		Ascent:  fixed.I(dotY),
		Descent: fixed.I(charHeight - dotY),
	}
}
