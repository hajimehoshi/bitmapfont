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

// Package bitmapfont provides font.Face values with bitmap glyphs.
package bitmapfont

import (
	"golang.org/x/image/font"
)

func init() {
	Face = newDelayedFace("data/face_ja.bin")
	FaceEA = newDelayedFace("data/face_ja_ea.bin")
}

var (
	// Face is a font.Face of the bitmap font (12px regular).
	Face font.Face

	// FaceEA is a font.Face of the bitmap font (12px regular, prefer East Asian wide characters).
	FaceEA font.Face
)
