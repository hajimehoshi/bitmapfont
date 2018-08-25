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

// +build example

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
	height := 120
	text := `en: All human beings are born free and equal in dignity and rights.
es: Todos los seres humanos nacen libres e iguales en dignidad y derechos y,
fr: Tous les êtres humains naissent libres et égaux en dignité et en droits.
de: Alle Menschen sind frei und gleich an Würde und Rechten geboren.
pt: Todos os seres humanos nascem livres e iguais em dignidade e em direitos.
ja: すべての人間は、生れながらにして自由であり、かつ、尊厳と権利とについて平等である。
ko: 모든 인간은 태어날 때부터 자유로우며 그 존엄과 권리에 있어 동등하다.
`
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

	f := bitmapfont.Gothic12r
	d := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(color.Black),
		Face: f,
		Dot:  fixed.P(ox, oy),
	}

	for _, l := range strings.Split(text, "\n") {
		d.DrawString(l)
		d.Dot.X = fixed.I(ox)
		d.Dot.Y += f.Metrics().Height
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

	if err := png.Encode(fout, d.Dst); err != nil {
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
