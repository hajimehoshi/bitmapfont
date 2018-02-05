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
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/browser"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/hajimehoshi/go-mplus-bitmap"
)

func writeImageToTempFile(img image.Image) (string, error) {
	f, err := ioutil.TempFile("", "mplus")
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return "", err
	}

	const suffix = ".png"
	n := f.Name()
	if err := os.Rename(n, n+suffix); err != nil {
		return "", err
	}
	return n + suffix, nil
}

const text = `Hello, World!

こんにちは世界!`

func run() error {
	const (
		ox = 16
		oy = 16
	)

	dst := image.NewRGBA(image.Rect(0, 0, 320, 240))
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.ZP, draw.Src)

	f := mplusbitmap.Gothic12r
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

	path, err := writeImageToTempFile(d.Dst)
	if err != nil {
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
