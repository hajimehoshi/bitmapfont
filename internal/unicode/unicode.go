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

package unicode

func IsLatin(r rune) bool {
	if r <= 0x02ff {
		// Basic Latin
		// Latin-1 Supplement
		// Latin Extended-A
		// Latin Extended-B
		// IPA Extensions
		// Spacing Modifier Letters
		return true
	}
	if 0x1d00 <= r && r <= 0x1dbf {
		// Phonetic Extensions
		// Phonetic Extensions Supplement
		return true
	}
	if 0x1e00 <= r && r <= 0x1eff {
		// Latin Extended Additional
		return true
	}
	if 0x2070 <= r && r <= 0x209f {
		// Superscripts and Subscripts
		return true
	}
	if 0x2c60 <= r && r <= 0x2c7f {
		// Latin Extended-C
		return true
	}
	if 0xa720 <= r && r <= 0xa7ff {
		// Latin Extended-D
		return true
	}
	if 0xab30 <= r && r <= 0xab6f {
		// Latin Extended-E
		return true
	}
	if 0xfb00 <= r && r <= 0xfb06 {
		// Alphabetic Presentation Forms
		return true
	}
	return false
}

func IsGreek(r rune) bool {
	if 0x0370 <= r && r <= 0x03ff {
		// Greek and Coptic
		return true
	}
	if 0x1d00 <= r && r <= 0x1dbf {
		// Phonetic Extensions
		// Phonetic Extensions Supplement
		return true
	}
	if 0x1F00 <= r && r <= 0x1fff {
		// Greek Extended
		return true
	}
	if 0xab30 <= r && r <= 0xab6f {
		// Latin Extended-E
		return true
	}
	return false
}

func IsCyrillic(r rune) bool {
	if 0x0400 <= r && r <= 0x052f {
		// Cyrillic
		// Cyrillic Supplement
		return true
	}
	if 0x1c80 <= r && r <= 0x1C8f {
		// Cyrillic Extended-C
		return true
	}
	if 0x1d00 <= r && r <= 0x1d7f {
		// Phonetic Extensions
		return true
	}
	if 0x2de0 <= r && r <= 0x2dff {
		// Cyrillic Extended-A
		return true
	}
	if 0xa640 <= r && r <= 0xa69f {
		// Cyrillic Extended-B
		return true
	}
	if 0xFE20 <= r && r <= 0xfe2f {
		// Combining Half Marks
		return true
	}
	return false
}

func IsArmenian(r rune) bool {
	if 0x0530 <= r && r <= 0x058f {
		// Armenian
		return true
	}
	if 0xfb00 <= r && r <= 0xfb4f {
		// Alphabetic Presentation Forms
		return true
	}
	return false
}

func IsHebrew(r rune) bool {
	if 0x0590 <= r && r <= 0x05ff {
		// Hebrew
		return true
	}
	if 0xfb00 <= r && r <= 0xfb4f {
		// Alphabetic Presentation Forms
		return true
	}
	return false
}

func IsThai(r rune) bool {
	if 0x0e00 <= r && r <= 0x0e7f {
		// Thai
		return true
	}
	return false
}

func IsGeorgian(r rune) bool {
	if 0x10a0 <= r && r <= 0x10ff {
		// Georgian
		return true
	}
	if 0x1c90 <= r && r <= 0x1cbf {
		// Georgian Extended
		return true
	}
	if 0x2d00 <= r && r <= 0x2d2f {
		// Georgian Supplement
		return true
	}
	return false
}

func IsOgham(r rune) bool {
	if 0x1680 <= r && r <= 0x169f {
		// Ogham
		return true
	}
	return false
}

func IsRunic(r rune) bool {
	if 0x16A0 <= r && r <= 0x16FF {
		// Runic
		return true
	}
	return false
}

func IsGeneralPunctuation(r rune) bool {
	if 0x2000 <= r && r <= 0x206f {
		return true
	}
	return false
}

func IsSupplementalPunctuation(r rune) bool {
	if 0x2e00 <= r && r <= 0x2e7f {
		return true
	}
	return false
}

func IsCJKUnifiedIdeograph(r rune) bool {
	// CJK Unified Ideographs
	if 0x4E00 <= r && r <= 0x9FFF {
		return true
	}
	// CJK Unified Ideographs Extension A
	if 0x3400 <= r && r <= 0x4DBF {
		return true
	}
	return false
}
