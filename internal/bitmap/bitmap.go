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

package bitmap

import (
	"fmt"
	"image"
	"image/color"
	"unicode"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/width"

	iunicode "github.com/hajimehoshi/bitmapfont/v3/internal/unicode"
)

type BinaryImage struct {
	bits   []byte
	width  int
	height int
	bounds image.Rectangle
}

func NewBinaryImage(bits []byte, width, height int) *BinaryImage {
	return &BinaryImage{
		bits:   bits,
		width:  width,
		height: height,
		bounds: image.Rect(0, 0, width, height),
	}
}

func (b *BinaryImage) At(i, j int) color.Color {
	if i < b.bounds.Min.X || j < b.bounds.Min.Y || i >= b.bounds.Max.X || j >= b.bounds.Max.Y {
		return color.Alpha{0}
	}
	idx := b.width*j + i
	if (b.bits[idx/8]>>uint(7-idx%8))&1 != 0 {
		return color.Alpha{0xff}
	}
	return color.Alpha{0}
}

func (b *BinaryImage) ColorModel() color.Model {
	return color.AlphaModel
}

func (b *BinaryImage) Bounds() image.Rectangle {
	return b.bounds
}

func (b *BinaryImage) SubImage(r image.Rectangle) image.Image {
	bounds := r.Intersect(b.bounds)
	if bounds.Empty() {
		return &BinaryImage{}
	}
	return &BinaryImage{
		bits:   b.bits,
		width:  b.width,
		height: b.height,
		bounds: bounds,
	}
}

const (
	charXNum = 256
	charYNum = 256
)

type Face struct {
	image        *BinaryImage
	dotX         fixed.Int26_6
	dotY         fixed.Int26_6
	eastAsiaWide bool
}

func NewFace(image *BinaryImage, dotX, dotY fixed.Int26_6, eastAsiaWide bool) *Face {
	return &Face{
		image:        image,
		dotX:         dotX,
		dotY:         dotY,
		eastAsiaWide: eastAsiaWide,
	}
}

func (f *Face) runeWidth(r rune) int {
	// For Latin glyphs, M+ doesn't work. Use the fixed font whatever the face is.
	if iunicode.IsLatin(r) {
		return f.charHalfWidth()
	}

	if _, ok := wideRunes[r]; ok {
		return f.charFullWidth()
	}

	switch k := width.LookupRune(r).Kind(); k {
	case width.Neutral:
		return f.charHalfWidth()
	case width.EastAsianAmbiguous:
		if f.eastAsiaWide {
			return f.charFullWidth()
		}
		return f.charHalfWidth()
	case width.EastAsianWide:
		return f.charFullWidth()
	case width.EastAsianNarrow:
		return f.charHalfWidth()
	case width.EastAsianFullwidth:
		return f.charFullWidth()
	case width.EastAsianHalfwidth:
		return f.charHalfWidth()
	default:
		panic(fmt.Sprintf("bitmap: unexpected kind: %d", k))
	}
}

func (f *Face) charFullWidth() int {
	return f.image.Bounds().Dx() / charXNum
}

func (f *Face) charHalfWidth() int {
	return f.charFullWidth() / 2
}

func (f *Face) charHeight() int {
	return f.image.Bounds().Dy() / charYNum
}

func (f *Face) Close() error {
	return nil
}

func (f *Face) Glyph(dot fixed.Point26_6, r rune) (dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
	if r >= 0x10000 {
		return
	}

	rw := f.runeWidth(r)
	dx := (dot.X - f.dotX).Floor()
	dy := (dot.Y - f.dotY).Floor()
	dr = image.Rect(dx, dy, dx+rw, dy+f.charHeight())

	mx := (int(r) % charXNum) * f.charFullWidth()
	my := (int(r) / charXNum) * f.charHeight()
	mask = f.image.SubImage(image.Rect(mx, my, mx+rw, my+f.charHeight()))
	maskp = image.Pt(mx, my)
	advance = fixed.I(f.runeWidth(r))
	ok = true
	return
}

func (f *Face) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	if r >= 0x10000 {
		return
	}
	bounds = fixed.Rectangle26_6{
		Min: fixed.Point26_6{X: -f.dotX, Y: -f.dotY},
		Max: fixed.Point26_6{X: -f.dotX + fixed.I(f.runeWidth(r)), Y: -f.dotY + fixed.I(f.charHeight())},
	}
	advance = fixed.I(f.runeWidth(r))
	ok = true
	return
}

func (f *Face) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	if r >= 0x10000 {
		return 0, false
	}
	return fixed.I(f.runeWidth(r)), true
}

func (f *Face) Kern(r0, r1 rune) fixed.Int26_6 {
	if unicode.Is(unicode.Mn, r1) {
		return -fixed.I(f.runeWidth(r1))
	}
	return 0

}

func (f *Face) Metrics() font.Metrics {
	return font.Metrics{
		Height:  fixed.I(f.charHeight()),
		Ascent:  f.dotY,
		Descent: fixed.I(f.charHeight()) - f.dotY,
	}
}
