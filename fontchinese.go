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

//go:build chinese
// +build chinese

package bitmapfont

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/hajimehoshi/bitmapfont/v2/internal/bitmap"
)

func init() {
	const (
		imageWidth  = 12 * 256
		imageHeight = 16 * 256

		dotX = 4
		dotY = 12
	)

	s, err := gzip.NewReader(bytes.NewReader(compressedFontAlphaFaceChinese))
	if err != nil {
		panic(err)
	}
	defer s.Close()

	bits, err := ioutil.ReadAll(s)
	if err != nil {
		panic(err)
	}

	FaceChinese = bitmap.NewFace(bitmap.NewBinaryImage(bits, imageWidth, imageHeight), fixed.I(dotX), fixed.I(dotY), true)
}

var (
	// FaceChinese is a font.Face of the bitmap font (12px regular, prefer Chinese characters).
	FaceChinese font.Face
)
