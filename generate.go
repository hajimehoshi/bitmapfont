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

package bitmapfont

// Package github.com/hajimehoshi/png2compressedrgba is required.
// Package github.com/hajimehoshi/file2byteslice is required.

//go:generate go run -tags=generate ./internal/gen -size 12 -output /tmp/compressedFontAlpha12r
//go:generate file2byteslice -input /tmp/compressedFontAlpha12r -output image12r.go -package bitmapfont -var compressedFontAlpha12r

//go:generate go run -tags=generate ./internal/gen -size 12 -eastasia -output /tmp/compressedFontAlpha12rEastAsia
//go:generate file2byteslice -input /tmp/compressedFontAlpha12rEastAsia -output image12rea.go -package bitmapfont -var compressedFontAlpha12rEastAsia

//go:generate go run -tags=generate ./internal/gen -size 10 -output /tmp/compressedFontAlpha10r
//go:generate file2byteslice -input /tmp/compressedFontAlpha10r -output image10r.go -package bitmapfont -var compressedFontAlpha10r

//go:generate go run -tags=generate ./internal/gen -size 10 -output /tmp/compressedFontAlpha10rEastAsia
//go:generate file2byteslice -input /tmp/compressedFontAlpha10rEastAsia -output image10rea.go -package bitmapfont -var compressedFontAlpha10rEastAsia

//go:generate gofmt -s -w .
