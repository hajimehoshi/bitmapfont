// Copyright 2020 Hajime Hoshi
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

package arabic

import (
	"image"
	"image/color"
)

type mirroredImage struct {
	img image.Image
}

func (m mirroredImage) ColorModel() color.Model {
	return m.img.ColorModel()
}

func (m mirroredImage) Bounds() image.Rectangle {
	return m.img.Bounds()
}

func (m mirroredImage) At(x, y int) color.Color {
	b := m.Bounds()
	x -= b.Min.X
	x = b.Dx() - x
	x += b.Min.X
	return m.img.At(x, y)
}

type rotatedImage struct {
	img    image.Image
	shiftX int
	shiftY int
}

func (r rotatedImage) ColorModel() color.Model {
	return r.img.ColorModel()
}

func (r rotatedImage) Bounds() image.Rectangle {
	return r.img.Bounds()
}

func (r rotatedImage) At(x, y int) color.Color {
	b := r.Bounds()
	x -= b.Min.X
	y -= b.Min.Y
	x = b.Dx() - x
	y = b.Dy() - y
	x += b.Min.X
	y += b.Min.Y
	x += r.shiftX
	y += r.shiftY
	return r.img.At(x, y)
}
