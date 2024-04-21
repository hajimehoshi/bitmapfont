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

package bitmapfont

import (
	"compress/gzip"
	"io"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/hajimehoshi/bitmapfont/v3/internal/bitmap"
)

func init() {
	f, err := data.Open("data/face_zhhans.bin")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	bits, err := io.ReadAll(s)
	if err != nil {
		panic(err)
	}

	FaceSC = bitmap.NewFace(bitmap.NewBinaryImage(bits, imageWidth, imageHeight), fixed.I(dotX), fixed.I(dotY), false)
}

func init() {
	f, err := data.Open("data/face_zhhans_ea.bin")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	bits, err := io.ReadAll(s)
	if err != nil {
		panic(err)
	}
	FaceSCEA = bitmap.NewFace(bitmap.NewBinaryImage(bits, imageWidth, imageHeight), fixed.I(dotX), fixed.I(dotY), true)
}

var (
	// FaceSC is a font.Face of the bitmap font (12px regular, prefer simplified Chinese characters).
	FaceSC font.Face

	// FaceSCEA is a font.Face of the bitmap font (12px regular, prefer simplified Chinese characters and East Asia wide characters).
	FaceSCEA font.Face
)
