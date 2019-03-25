// Copyright 2019 Hajime Hoshi
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

// +build example zh

package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"

	"github.com/pkg/browser"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/hajimehoshi/bitmapfont"
)

var (
	flagTest = flag.Bool("test", false, "test mode")
)

func run() error {
	const (
		ox = 16
		oy = 16
	)

	width := 640

	// https://www.unicode.org/udhr/
	// https://omniglot.com/udhr/
	text := `en: All human beings are born free and equal in dignity and rights.
zh-Hans: 人人生而自由,在尊严和权利上一律平等。他们赋有理性和良心,并应以兄弟关系的精神相对待。
zh-Hant: 人人生而自由，在尊嚴和權利上一律平等。他們賦有理性和良心，並應以兄弟關係的精神相對待。
`
	height := 16*len(strings.Split(strings.TrimSpace(text), "\n")) + 8
	if *flagTest {
		width = 12*256 + 16
		height = 16*256 + 16
		text = ""
		for i := 0; i < 256; i++ {
			for j := 0; j < 256; j++ {
				r := rune(i*256 + j)
				if r == '\n' {
					text += " "
					continue
				}
				text += string(r)
			}
			text += "\n"
		}
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.ZP, draw.Src)

	faces := []font.Face{
		bitmapfont.Gothic12r,
		bitmapfont.Gothic12r_ZhHans,
		bitmapfont.Gothic12r_ZhHant,
	}
	for i, l := range strings.Split(text, "\n") {
		if len(faces) <= i {
			break
		}
		f := faces[i]
		d := font.Drawer{
			Dst:  dst,
			Src:  image.NewUniform(color.Black),
			Face: f,
		}
		d.Dot.X = fixed.I(ox)
		d.Dot.Y = fixed.I(oy) + fixed.Int26_6(i)*f.Metrics().Height
		d.DrawString(l)
	}

	path := "example.png"
	if *flagTest {
		path = "example_test.png"
	}
	fout, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fout.Close()

	if err := png.Encode(fout, dst); err != nil {
		return err
	}

	if err := browser.OpenFile(path); err != nil {
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
