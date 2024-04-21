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

//go:generate go run -C=_gen . -widths -output ./../internal/bitmap/widths.go

//go:generate go run -C=_gen . -lang="ja" -output ./../data/face_ja.bin
//go:generate go run -C=_gen . -lang="ja" -eastasia -output ./../data/face_ja_ea.bin
//go:generate go run -C=_gen . -lang="zh-Hans" -output ./../data/face_zhhans.bin
//go:generate go run -C=_gen . -lang="zh-Hans" -eastasia -output ./../data/face_zhhans_ea.bin
//go:generate go run -C=_gen . -lang="zh-Hant" -output ./../data/face_zhhant.bin
//go:generate go run -C=_gen . -lang="zh-Hant" -eastasia -output ./../data/face_zhhant_ea.bin

//go:generate gofmt -s -w .
