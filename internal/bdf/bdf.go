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

package bdf

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"io"
	"strconv"
	"strings"
)

type Glyph struct {
	Encoding int
	Width    int
	Height   int
	X        int
	Y        int
	Bitmap   [][]byte

	ShiftX int
	ShiftY int
}

func (g *Glyph) ColorModel() color.Model {
	return color.AlphaModel
}

func (g *Glyph) Bounds() image.Rectangle {
	return image.Rect(0, 0, g.Width, g.Height)
}

func (g *Glyph) At(x, y int) color.Color {
	x -= g.ShiftX
	y -= g.ShiftY
	if x < 0 || y < 0 || x >= g.Width || y >= g.Height {
		return color.Alpha{}
	}
	bits := g.Bitmap[y][x/8]
	if (bits>>uint(7-x%8))&1 != 0 {
		return color.Alpha{0xff}
	}
	return color.Alpha{}
}

func Parse(f io.Reader) ([]*Glyph, error) {
	glyphs := []*Glyph{}
	var current *Glyph

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "STARTCHAR ") {
			if current != nil {
				panic("not reach")
			}
			current = &Glyph{}
			continue
		}
		if strings.HasPrefix(line, "ENCODING ") {
			idx, err := strconv.ParseInt(line[len("ENCODING "):], 10, 32)
			if err != nil {
				return nil, err
			}
			current.Encoding = int(idx)
			continue
		}
		if strings.HasPrefix(line, "BBX ") {
			tokens := strings.Split(line, " ")
			w, err := strconv.ParseInt(tokens[1], 10, 32)
			if err != nil {
				return nil, err
			}
			h, err := strconv.ParseInt(tokens[2], 10, 32)
			if err != nil {
				return nil, err
			}
			x, err := strconv.ParseInt(tokens[3], 10, 32)
			if err != nil {
				return nil, err
			}
			y, err := strconv.ParseInt(tokens[4], 10, 32)
			if err != nil {
				return nil, err
			}
			current.Width = int(w)
			current.Height = int(h)
			current.X = int(x)
			current.Y = int(y)
		}
		if strings.HasPrefix(line, "BITMAP") {
			current.Bitmap = [][]byte{}
			continue
		}
		if strings.HasPrefix(line, "ENDCHAR") {
			glyphs = append(glyphs, current)
			current = nil
			continue
		}
		if current == nil {
			continue
		}
		if current.Bitmap == nil {
			continue
		}
		if len(line)%2 != 0 {
			return nil, fmt.Errorf("bdf: len(line) must be even")
		}
		bits := []byte{}
		for ; len(line) > 0; line = line[2:] {
			b, err := strconv.ParseInt(line[:2], 16, 32)
			if err != nil {
				return nil, err
			}
			bits = append(bits, byte(b))
		}
		current.Bitmap = append(current.Bitmap, bits)
	}
	return glyphs, nil
}
