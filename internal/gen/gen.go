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

// +build generate

package main

import (
	"compress/gzip"
	"flag"
	"image"
	"image/color"
	"image/draw"
	"os"

	"github.com/hajimehoshi/bitmapfont/internal/baekmuk"
	"github.com/hajimehoshi/bitmapfont/internal/bdf"
	"github.com/hajimehoshi/bitmapfont/internal/fixed"
	"github.com/hajimehoshi/bitmapfont/internal/mplus"
)

var (
	flagOutput = flag.String("output", "", "output file")
)

const (
	glyphWidth  = 12
	glyphHeight = 16
)

type fontType int

const (
	fontTypeNone fontType = iota
	fontTypeFixed
	fontTypeMPlus
	fontTypeBaekmuk
)

func getFontType(r rune) fontType {
	if 0x2500 <= r && r <= 0x257f {
		// Box Drawing
		// M+ defines a part of box drawing glyphs.
		// For consistency, use other font's glyphs instead.
		return fontTypeBaekmuk
	}
	if 0xff65 <= r && r <= 0xff9f {
		// Halfwidth Katakana
		return fontTypeMPlus
	}
	if _, ok := fixed.Glyph(r); ok {
		return fontTypeFixed
	}
	if _, ok := mplus.Glyph(r); ok {
		return fontTypeMPlus
	}
	if _, ok := baekmuk.Glyph(r); ok {
		return fontTypeBaekmuk
	}
	return fontTypeNone
}

func getGlyph(r rune) (bdf.Glyph, bool) {
	switch getFontType(r) {
	case fontTypeNone:
		return bdf.Glyph{}, false
	case fontTypeFixed:
		g, ok := fixed.Glyph(r)
		if ok {
			return g, true
		}
	case fontTypeMPlus:
		g, ok := mplus.Glyph(r)
		if ok {
			return g, true
		}
	case fontTypeBaekmuk:
		g, ok := baekmuk.Glyph(r)
		if ok {
			return g, true
		}
	default:
		panic("not reached")
	}
	return bdf.Glyph{}, false
}

func addGlyphs(img draw.Image) {
	for j := 0; j < 0x100; j++ {
		for i := 0; i < 0x100; i++ {
			r := rune(i + j*0x100)
			g, ok := getGlyph(r)
			if !ok {
				continue
			}

			dstX := i*glyphWidth + g.X
			dstY := j*glyphHeight + ((glyphHeight - g.Height) - 4 - g.Y)
			dstR := image.Rect(dstX, dstY, dstX+g.Width, dstY+g.Height)
			p := g.Bounds().Min
			draw.Draw(img, dstR, &g, p, draw.Over)
		}
	}
}

func run() error {
	img := image.NewAlpha(image.Rect(0, 0, glyphWidth*256, glyphHeight*256))
	addGlyphs(img)

	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	as := make([]byte, w*h/8)
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			a := img.At(i, j).(color.Alpha).A
			idx := w*j + i
			if a != 0 {
				as[idx/8] |= 1 << uint(7-idx%8)
			}
		}
	}

	fout, err := os.Create(*flagOutput)
	if err != nil {
		return err
	}
	defer fout.Close()

	cw, err := gzip.NewWriterLevel(fout, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer cw.Close()

	if _, err := cw.Write(as); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		panic(err)
	}
}
