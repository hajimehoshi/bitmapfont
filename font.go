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

// Package bitmapfont offers a font.Face value of some bitmap fonts.
//
//   * [Baekmuk Gulim](https://kldp.net/baekmuk/) (Baekmuk License)
//   * [misc-fixed](https://www.cl.cam.ac.uk/~mgk25/ucs-fonts.html) (Public Domain)
//   * [M+ Bitmap Font](https://mplus-fonts.osdn.jp/mplus-bitmap-fonts/) (M+ Bitmap Fonts License)
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

	s, err := gzip.NewReader(bytes.NewReader(compressedFontAlphaFace))
	if err != nil {
		panic(err)
	}
	defer s.Close()

	bits, err := ioutil.ReadAll(s)
	if err != nil {
		panic(err)
	}

	Face = bitmap.NewFace(bitmap.NewBinaryImage(bits, imageWidth, imageHeight), fixed.I(dotX), fixed.I(dotY), false)
}

func init() {
	const (
		imageWidth  = 12 * 256
		imageHeight = 16 * 256

		dotX = 4
		dotY = 12
	)

	s, err := gzip.NewReader(bytes.NewReader(compressedFontAlphaFaceEA))
	if err != nil {
		panic(err)
	}
	defer s.Close()

	bits, err := ioutil.ReadAll(s)
	if err != nil {
		panic(err)
	}

	FaceEA = bitmap.NewFace(bitmap.NewBinaryImage(bits, imageWidth, imageHeight), fixed.I(dotX), fixed.I(dotY), true)
}

var (
	// Face is a font.Face of the bitmap font (12px regular).
	Face font.Face

	// FaceEA is a font.Face of the bitmap font (12px regular, prefer East Asian wide characters).
	FaceEA font.Face
)
