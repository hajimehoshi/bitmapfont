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

func run() error {
	const (
		ox = 16
		oy = 16
	)

	dst := image.NewRGBA(image.Rect(0, 0, 640, 120))
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.ZP, draw.Src)

	f := bitmapfont.Gothic12r
	d := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(color.Black),
		Face: f,
		Dot:  fixed.P(ox, oy),
	}

	text := `All human beings are born free and equal in dignity and rights.

すべての人間は、生れながらにして自由であり、かつ、尊厳と権利とについて平等である。

모든 인간은 태어날 때부터 자유로우며 그 존엄과 권리에 있어 동등하다.
`
	for _, l := range strings.Split(text, "\n") {
		d.DrawString(l)
		d.Dot.X = fixed.I(ox)
		d.Dot.Y += f.Metrics().Height
	}

	path := "example.png"
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
	if err := run(); err != nil {
		panic(err)
	}
}
