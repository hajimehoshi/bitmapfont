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

package bitmapfont

import (
	"image"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func init() {
	FaceTC = &tcFace{face: newDelayedFace("data/face_zhhant.bin", false)}
	FaceTCEA = &tcFace{face: newDelayedFace("data/face_zhhant_ea.bin", true)}
}

var _ font.Face = (*tcFace)(nil)

type tcFace struct {
	face font.Face
}

func (t *tcFace) Close() error {
	return t.face.Close()
}

func (t *tcFace) Glyph(dot fixed.Point26_6, r rune) (dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
	dr, mask, maskp, advance, ok = t.face.Glyph(dot, r)
	if r == '、' || r == '，' || r == '。' || r == '．' {
		dr = dr.Add(image.Pt(3, -3))
	}
	return dr, mask, maskp, advance, ok
}

func (t *tcFace) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	return t.face.GlyphBounds(r)
}

func (t *tcFace) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	return t.face.GlyphAdvance(r)
}

func (t *tcFace) Kern(r0, r1 rune) fixed.Int26_6 {
	return t.face.Kern(r0, r1)
}

func (t *tcFace) Metrics() font.Metrics {
	return t.face.Metrics()
}

var (
	// FaceTC is a font.Face of the bitmap font (12px regular, prefer traditional Chinese characters).
	FaceTC font.Face

	// FaceTCEA is a font.Face of the bitmap font (12px regular, prefer traditional Chinese characters and East Asia wide characters).
	FaceTCEA font.Face
)
