// Copyright 2025 Hajime Hoshi
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

package unicode_test

import (
	"testing"

	"github.com/hajimehoshi/bitmapfont/v3/internal/unicode"
)

func TestComposeHangulSyllable(t *testing.T) {
	testCases := []struct {
		l    rune
		v    rune
		t    rune
		want rune
	}{
		{
			l:    'ᄀ',
			v:    'ᅡ',
			want: '가',
		},
		{
			l:    'ᄀ',
			v:    'ᅡ',
			t:    'ᆨ',
			want: '각',
		},
		{
			l:    'ᄒ',
			v:    'ᅵ',
			t:    'ᆼ',
			want: '힝',
		},
	}

	for _, tc := range testCases {
		got := unicode.ComposeHangulSyllable(tc.l, tc.v, tc.t)
		if got != tc.want {
			t.Errorf("ComposeHangulSyllable(%U, %U, %U) = %c; want %c", tc.l, tc.v, tc.t, got, tc.want)
		}
	}
}

func TestDecomposeHangulSyllable(t *testing.T) {
	testCases := []struct {
		syllable rune
		wantL    rune
		wantV    rune
		wantT    rune
	}{
		{
			syllable: '가',
			wantL:    'ᄀ',
			wantV:    'ᅡ',
			wantT:    0,
		},
		{
			syllable: '각',
			wantL:    'ᄀ',
			wantV:    'ᅡ',
			wantT:    'ᆨ',
		},
		{
			syllable: '힝',
			wantL:    'ᄒ',
			wantV:    'ᅵ',
			wantT:    'ᆼ',
		},
	}

	for _, tc := range testCases {
		gotL, gotV, gotT := unicode.DecomposeHangulSyllable(tc.syllable)
		if gotL != tc.wantL || gotV != tc.wantV || gotT != tc.wantT {
			t.Errorf("DecomposeHangulSyllable(%U) = %c, %c, %c; want %c, %c, %c", tc.syllable, gotL, gotV, gotT, tc.wantL, tc.wantV, tc.wantT)
		}
	}
}
