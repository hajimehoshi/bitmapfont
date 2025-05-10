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
	"embed"
	"image"
	"io"
	"sync"

	"github.com/pierrec/lz4/v4"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/hajimehoshi/bitmapfont/v4/internal/bitmap"
)

//go:embed data/*.bin
var data embed.FS

const (
	imageWidth  = 12 * 256
	imageHeight = 16 * 256

	dotX = 0
	dotY = 12
)

var _ font.Face = (*lazyFace)(nil)

type lazyFace struct {
	binFile  string
	ea       bool
	initOnce sync.Once
	face     font.Face
}

func newDelayedFace(binFile string, ea bool) *lazyFace {
	return &lazyFace{
		binFile: binFile,
		ea:      ea,
	}
}

func (f *lazyFace) ensureInitialization() {
	f.initOnce.Do(func() {
		binFile, err := data.Open(f.binFile)
		if err != nil {
			panic(err)
		}
		defer binFile.Close()

		s := lz4.NewReader(binFile)

		bits, err := io.ReadAll(s)
		if err != nil {
			panic(err)
		}

		f.face = bitmap.NewFace(bitmap.NewBinaryImage(bits, imageWidth, imageHeight), fixed.I(dotX), fixed.I(dotY), f.ea)
	})
}

func (f *lazyFace) Close() error {
	// Checking whether f.face is nil ro not is not concurrent-safe.
	// However, this is risky only when Close and other methods are called concurrently.
	// Such a situation is not realistic.
	if f.face == nil {
		return nil
	}
	if err := f.face.Close(); err != nil {
		return err
	}
	return nil
}

func (f *lazyFace) Glyph(dot fixed.Point26_6, r rune) (dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
	f.ensureInitialization()
	return f.face.Glyph(dot, r)
}

func (f *lazyFace) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	f.ensureInitialization()
	return f.face.GlyphBounds(r)
}

func (f *lazyFace) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	f.ensureInitialization()
	return f.face.GlyphAdvance(r)
}

func (f *lazyFace) Kern(r0, r1 rune) fixed.Int26_6 {
	f.ensureInitialization()
	return f.face.Kern(r0, r1)
}

func (f *lazyFace) Metrics() font.Metrics {
	f.ensureInitialization()
	return f.face.Metrics()
}
