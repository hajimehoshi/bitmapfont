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
	"bufio"
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/hajimehoshi/bitmapfont/bdf"
)

var (
	flagOutput = flag.String("output", "", "output file")
	flagTest   = flag.Bool("test", false, "test output")
)

func conv() (map[int]rune, error) {
	_, current, _, _ := runtime.Caller(1)
	dir := filepath.Dir(current)

	f, err := os.Open(filepath.Join(dir, "KSX1001.TXT"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := map[int]rune{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if idx := strings.Index(line, "#"); idx != -1 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		tokens := strings.Split(line, "  ")
		if len(tokens) != 2 {
			continue
		}
		ksx, err := strconv.ParseInt(tokens[0], 0, 32)
		if err != nil {
			return nil, err
		}
		uni, err := strconv.ParseInt(tokens[1], 0, 32)
		if err != nil {
			return nil, err
		}
		m[int(ksx)] = rune(uni)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return m, nil
}

type glyph struct {
	bits [][]byte
	ox   int
	oy   int
}

func readBDF() (map[rune]glyph, error) {
	c, err := conv()
	if err != nil {
		return nil, err
	}

	_, current, _, _ := runtime.Caller(1)
	dir := filepath.Dir(current)

	f, err := os.Open(filepath.Join(dir, "gulim12.bdf"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := map[rune]glyph{}

	glyphs, err := bdf.Parse(f)
	if err != nil {
		return nil, err
	}
	for _, g := range glyphs {
		r, ok := c[g.Encoding]
		if !ok {
			continue
		}
		ox := int(g.X)
		oy := int(g.Y) + (int(g.Height) - 12)
		m[r] = glyph{
			bits: g.Bitmap,
			ox:   ox,
			oy:   oy,
		}
	}

	return m, nil
}

func isRuneToDraw(r rune) bool {
	// KOREAN STANDARD SYMBOL
	if r == 0x327f {
		return true
	}
	// CIRCLE HANGLE
	if 0x3260 <= r && r <= 0x327e {
		return true
	}
	// PARENTHESIZED HANGUL
	if 0x3200 <= r && r <= 0x321f {
		return true
	}
	// HANGUL LETTER (Hangul Compatible Jamo)
	if 0x3130 <= r && r <= 0x318f {
		return true
	}
	// HANGUL SYLLABLE
	if 0xac00 <= r && r <= 0xd7af {
		return true
	}
	return false
}

func run() error {
	b, err := readBDF()
	if err != nil {
		return err
	}

	const (
		num  = 65536
		numX = 256
		srcW = 12
		srcH = 12
		dstW = 12
		dstH = 16
	)

	offsetY := 0
	if *flagTest {
		offsetY = 3 * dstH
	}
	result := image.NewRGBA(image.Rect(0, 0, numX*dstW, ((num-1)/numX+1)*dstH+offsetY))

	runeToPos := func(r rune) image.Point {
		return image.Pt(int(r%numX)*dstW, int(r/numX)*dstH+offsetY)
	}

	for r, g := range b {
		if !isRuneToDraw(r) {
			continue
		}
		pos := runeToPos(r)
		for j := 0; j < len(g.bits); j++ {
			w := srcW
			if w < len(g.bits[j])*8 {
				w = len(g.bits[j]) * 8
			}
			for i := 0; i < w; i++ {
				bits := g.bits[j][i/8]
				if (bits>>uint(7-i%8))&1 != 0 {
					result.Set(pos.X+i+g.ox, pos.Y+j-g.oy, color.White)
				}
			}
		}
	}

	// Draw test text
	if *flagTest {
		const testText = "수학에서, 편미분 방정식 은 여러 개의 독립 변수로 구성된 함수와 그 함수의 편미분으로 연관된 방정식이다. 각각의 변수들의 상관관계를 고려하지 않고 변화량을 보고 싶을 때 이용할 수 있으며, 상미분방정식에 비해 응용범위가 훨씬 크다. 소리나 열의 전파 과정, 전자기학, 유체역학, 양자역학 등 수많은 역학계에 관련된 예가 많다."
		x := 0
		for _, r := range testText {
			if r == ' ' {
				x += dstW / 2
				continue
			}
			if !isRuneToDraw(r) {
				continue
			}
			pos := runeToPos(r)
			draw.Draw(result, image.Rect(x, 0, x+dstW, dstH), result, pos, draw.Over)
			x += dstW
		}
	}

	// Replace transparent part with black
	if *flagTest {
		r := result.Bounds()
		for i := r.Min.X; i < r.Max.X; i++ {
			for j := r.Min.Y; j < r.Max.Y; j++ {
				c := result.At(i, j).(color.RGBA)
				if c.A == 0 {
					result.Set(i, j, color.Black)
				}
			}
		}
	}

	fout, err := os.Create(*flagOutput)
	if err != nil {
		return err
	}
	defer fout.Close()

	if err := png.Encode(fout, result); err != nil {
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
