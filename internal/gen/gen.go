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

	"golang.org/x/text/width"

	"github.com/hajimehoshi/bitmapfont/v2/internal/baekmuk"
	"github.com/hajimehoshi/bitmapfont/v2/internal/bdf"
	"github.com/hajimehoshi/bitmapfont/v2/internal/fixed"
	"github.com/hajimehoshi/bitmapfont/v2/internal/mplus"
)

var (
	flagOutput   = flag.String("output", "", "output file")
	flagSize     = flag.Int("size", 12, "font size in pixels")
	flagEastAsia = flag.Bool("eastasia", false, "prefer east Asia punctuations")
)

func glyphSize(size int) (width, height int) {
	switch size {
	case 10:
		return 10, 12
	case 12:
		return 12, 16
	default:
		panic("not reached")
	}
}

type fontType int

const (
	fontTypeNone fontType = iota
	fontTypeFixed
	fontTypeMPlus
	fontTypeBaekmuk
)

func getFontType(r rune, size int) fontType {
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

	if width.LookupRune(r).Kind() == width.EastAsianAmbiguous {
		if *flagEastAsia {
			return fontTypeMPlus
		}
		return fontTypeFixed
	}

	if _, ok := fixed.Glyph(r, size); ok {
		return fontTypeFixed
	}
	if _, ok := mplus.Glyph(r, size); ok {
		return fontTypeMPlus
	}
	if _, ok := baekmuk.Glyph(r, size); ok {
		return fontTypeBaekmuk
	}
	return fontTypeNone
}

func getGlyph(r rune, size int) (bdf.Glyph, bool) {
	switch getFontType(r, size) {
	case fontTypeNone:
		return bdf.Glyph{}, false
	case fontTypeFixed:
		g, ok := fixed.Glyph(r, size)
		if ok {
			return g, true
		}
	case fontTypeMPlus:
		g, ok := mplus.Glyph(r, size)
		if ok {
			return g, true
		}
	case fontTypeBaekmuk:
		g, ok := baekmuk.Glyph(r, size)
		if ok {
			return g, true
		}
	default:
		panic("not reached")
	}
	return bdf.Glyph{}, false
}

func addGlyphs(img draw.Image, size int) {
	gw, gh := glyphSize(size)
	for j := 0; j < 0x100; j++ {
		for i := 0; i < 0x100; i++ {
			r := rune(i + j*0x100)
			g, ok := getGlyph(r, size)
			if !ok {
				continue
			}

			dstX := i*gw + g.X
			dstY := j*gh + ((gh - g.Height) - 4 - g.Y)
			dstR := image.Rect(dstX, dstY, dstX+g.Width, dstY+g.Height)
			p := g.Bounds().Min
			draw.Draw(img, dstR, &g, p, draw.Over)
		}
	}
}

func run() error {
	s := *flagSize

	gw, gh := glyphSize(s)
	img := image.NewAlpha(image.Rect(0, 0, gw*256, gh*256))
	addGlyphs(img, s)

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
