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

// +build ignore

package main

import (
	"compress/gzip"
	"flag"
	"image"
	"image/draw"
	"image/png"
	"os"
)

var (
	flagInput       = flag.String("input", "", "input file")
	flagInputHangul = flag.String("inputhangul", "", "input Hangul file")
	flagOutput      = flag.String("output", "", "output file")
)

const (
	glyphWidth  = 12
	glyphHeight = 16
)

func addHangul(img draw.Image) error {
	finh, err := os.Open(*flagInputHangul)
	if err != nil {
		return err
	}
	defer finh.Close()

	imgh, err := png.Decode(finh)
	if err != nil {
		return err
	}

	b := imgh.Bounds()
	oy := 0xac00 / 256 * glyphHeight
	dst := image.Rect(0, oy, b.Dx(), oy+b.Dy())
	draw.Draw(img, dst, imgh, image.ZP, draw.Over)
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
