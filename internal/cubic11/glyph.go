// Copyright 2024 Hajime Hoshi
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

package cubic11

import (
	"image"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	"github.com/hajimehoshi/bitmapfont/v3/internal/unicode"
)

var face font.Face

func readTTF() (font.Face, error) {
	_, current, _, _ := runtime.Caller(1)
	dir := filepath.Dir(current)

	ttfContent, err := os.ReadFile(filepath.Join(dir, "Cubic_11_1.200_R.ttf"))
	if err != nil {
		return nil, err
	}

	ttf, err := opentype.Parse(ttfContent)
	if err != nil {
		return nil, err
	}

	f, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, err
	}

	return f, nil
}

func init() {
	f, err := readTTF()
	if err != nil {
		panic(err)
	}
	face = f
}

func Glyph(r rune) (image.Image, bool) {
	if !unicode.IsCJKUnifiedIdeograph(r) {
		return nil, false
	}

	dr, img, _, _, ok := face.Glyph(fixed.Point26_6{}, r)
	if !ok {
		return nil, false
	}
	aimg := img.(*image.Alpha)
	aimg.Rect = aimg.Rect.Add(dr.Min).Add(image.Pt(-1, 11))

	return aimg, true
}
