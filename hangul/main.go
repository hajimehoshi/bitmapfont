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
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

const (
	initialG rune = 0x1100 + iota
	initialGG
	initialN
	initialD
	initialDD
	initialR
	initialM
	initialB
	initialBB
	initialS
	initialSS
	initialIeung
	initialJ
	initialJJ
	initialC
	initialK
	initialT
	initialP
	initialH
)

const (
	middleA rune = 0x1161 + iota
	middleAE
	middleYA
	middleYAE
	middleEO
	middleE
	middleYEO
	middleYE
	middleO
	middleWA
	middleWAE
	middleOE
	middleYO
	middleU
	middleWEO
	middleWE
	middleWI
	middleYU
	middleEU
	middleYI
	middleI
)

const (
	finalG rune = 0x11A8 + iota
	finalGG
	finalGS
	finalN
	finalNJ
	finalNH
	finalD
	finalL
	finalLG
	finalLM
	finalLB
	finalLS
	finalLT
	finalLP
	finalLH
	finalM
	finalB
	finalBS
	finalS
	finalSS
	finalNG
	finalJ
	finalC
	finalK
	finalT
	finalP
	finalH
)

type Source struct {
	Image       image.Image
	GlyphWidth  int
	GlyphHeight int
}

func srcGroups(initial, middle, final rune) (int, int, int) {
	// Reference: http://mytears.org/resources/doc/Hangul/HANGUL.TXT
	if final == -1 {
		gInitial := 0
		gMiddle := 0

		switch middle {
		case middleA, middleAE, middleYA, middleYAE, middleEO, middleE, middleYEO, middleYE, middleI:
			gInitial = 0
		case middleO, middleYO, middleEU:
			gInitial = 1
		case middleU, middleYU:
			gInitial = 2
		case middleWA, middleWAE, middleOE, middleYI:
			gInitial = 3
		case middleWEO, middleWE, middleWI:
			gInitial = 4
		default:
			panic("not reached")
		}

		switch initial {
		case initialG, initialK:
			gMiddle = 0
		default:
			gMiddle = 1
		}

		return gInitial, gMiddle, 0
	}

	gInitial := 0
	gMiddle := 0
	gFinal := -1

	switch middle {
	case middleA, middleAE, middleYA, middleYAE, middleEO, middleE, middleYEO, middleYE, middleI:
		gInitial = 5
	case middleO, middleYO, middleU, middleYU, middleEU:
		gInitial = 6
	case middleWA, middleWAE, middleOE, middleWEO, middleWE, middleWI, middleYI:
		gInitial = 7
	default:
		panic("not reached")
	}
	switch initial {
	case initialG, initialK:
		gMiddle = 2
	default:
		gMiddle = 3
	}

	switch middle {
	case middleA, middleYA, middleWA:
		gFinal = 0
	case middleEO, middleYEO, middleOE, middleWEO, middleWI, middleYI, middleI:
		gFinal = 1
	case middleAE, middleYAE, middleE, middleYE, middleWAE, middleWE:
		gFinal = 2
	case middleO, middleYO, middleU, middleYU, middleEU:
		gFinal = 3
	default:
		panic("not reached")
	}

	return gInitial, gMiddle, gFinal
}

func (src *Source) Gen() image.Image {
	const (
		num  = int((initialH - initialG + 1) * (middleI - middleA + 1) * (finalH - finalG + 1 + 1))
		numX = 256
		dstW = 12
		dstH = 16
	)

	offsetY := 0
	if *flagTest {
		offsetY = 3 * dstH
	}
	result := image.NewRGBA(image.Rect(0, 0, numX*dstW, ((num-1)/numX+1)*dstH+offsetY))

	idx := 0
	finals := []rune{-1}
	for i := finalG; i <= finalH; i++ {
		finals = append(finals, i)
	}
	for i := initialG; i <= initialH; i++ {
		for j := middleA; j <= middleI; j++ {
			for _, k := range finals {
				gInitial, gMiddle, gFinal := srcGroups(i, j, k)
				w := src.GlyphWidth
				h := src.GlyphHeight

				dstX := (idx%numX)*dstW + (dstW-w)/2
				dstY := (idx/numX)*dstH + (dstH-h)/2 + offsetY
				dst := image.Rect(dstX, dstY, dstX+w, dstY+h)

				posInitial := image.Pt(int(i-initialG)*w, int(gInitial)*h)
				posMiddle := image.Pt(int(j-middleA)*w, int(gMiddle+8)*h)
				var posFinal image.Point
				if k == -1 {
					posFinal = image.Pt(0, int(gFinal+12)*h)
				} else {
					posFinal = image.Pt(int(k-finalG+1)*w, int(gFinal+12)*h)
				}
				draw.Draw(result, dst, src.Image, posInitial, draw.Src)
				draw.Draw(result, dst, src.Image, posMiddle, draw.Over)
				draw.Draw(result, dst, src.Image, posFinal, draw.Over)

				idx++
			}
		}
	}

	if *flagTest {
		const testText = "수학에서, 편미분 방정식 은 여러 개의 독립 변수로 구성된 함수와 그 함수의 편미분으로 연관된 방정식이다. 각각의 변수들의 상관관계를 고려하지 않고 변화량을 보고 싶을 때 이용할 수 있으며, 상미분방정식에 비해 응용범위가 훨씬 크다. 소리나 열의 전파 과정, 전자기학, 유체역학, 양자역학 등 수많은 역학계에 관련된 예가 많다."
		x := 0
		for _, r := range testText {
			if r == ' ' {
				x += dstW / 2
				continue
			}
			if 0xac00 <= r && r <= 0xd7a3 {
				pos := image.Pt(int(r%numX)*dstW, int((r-0xac00)/numX)*dstH+offsetY)
				draw.Draw(result, image.Rect(x, 0, x+dstW, dstH), result, pos, draw.Over)
				x += dstW
				continue
			}
		}
	}

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

	return result
}

var (
	flagInput  = flag.String("input", "", "input file")
	flagOutput = flag.String("output", "", "output file")
	flagTest   = flag.Bool("test", false, "test output")
)

func run() error {
	fin, err := os.Open(*flagInput)
	if err != nil {
		return err
	}
	defer fin.Close()

	in, err := png.Decode(fin)
	if err != nil {
		return err
	}
	src := &Source{
		Image:       in,
		GlyphWidth:  12,
		GlyphHeight: 12,
	}

	fout, err := os.Create(*flagOutput)
	if err != nil {
		return err
	}
	defer fout.Close()

	if err := png.Encode(fout, src.Gen()); err != nil {
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
