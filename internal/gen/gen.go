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
	"image/draw"
	"image/png"
	"os"

	"github.com/hajimehoshi/bitmapfont/internal/baekmuk"
)

var (
	flagInput  = flag.String("input", "", "input file")
	flagOutput = flag.String("output", "", "output file")
)

const (
	glyphWidth  = 12
	glyphHeight = 16
)

func addHangul(img draw.Image) error {
	for r := rune(0); r < 0xffff; r++ {
		g, ok := baekmuk.Glyph(r)
		if !ok {
			continue
		}
		dstX := (int(r)%256)*glyphWidth + g.X
		dstY := (int(r)/256)*glyphHeight + ((glyphHeight - g.Height) - 4 - g.Y)
		dstR := image.Rect(dstX, dstY, dstX+g.Width, dstY+g.Height)
		p := g.Bounds().Min
		draw.Draw(img, dstR, &g, p, draw.Over)
	}
	return nil
}

func run() error {
	fin, err := os.Open(*flagInput)
	if err != nil {
		return err
	}
	defer fin.Close()

	img, err := png.Decode(fin)
	if err != nil {
		return err
	}

	if err := addHangul(img.(*image.NRGBA)); err != nil {
		return err
	}

	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	rgba := make([]byte, w*h*4)
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			r, g, b, a := img.At(i, j).RGBA()
			rgba[4*(w*j+i)] = byte(r >> 8)
			rgba[4*(w*j+i)+1] = byte(g >> 8)
			rgba[4*(w*j+i)+2] = byte(b >> 8)
			rgba[4*(w*j+i)+3] = byte(a >> 8)
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

	if _, err := cw.Write(rgba); err != nil {
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