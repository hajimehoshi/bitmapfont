# bitmapfont (v3)

[![PkgGoDev](https://pkg.go.dev/badge/github.com/hajimehoshi/bitmapfont/v3)](https://pkg.go.dev/github.com/hajimehoshi/bitmapfont/v3)

Package bitmapfont offers a font.Face value of some bitmap fonts.

## 12px glyphs

```
var Face font.Face
var FaceEA font.Face
var FaceSC font.Face
var FaceSCEA font.Face
var FaceTC font.Face
var FaceTCEA font.Face
```

![Example](example.png)

The `EA` version includes wide glyphs for the characters that have East Asian ambiguous widths (e.g., `※`, `…`, `α`).

The `SC` version prefers simplified Chinese characters.

The `TC` version prefers traditional Chinese characters.

The only real difference between `SC` and `TC` is the position of punctuation marks.

## Sources

   * [Baekmuk Gulim](https://kldp.net/baekmuk/) (Baekmuk License)
   * [Cubic 11](https://github.com/ACh-K/Cubic-11) (OFL-1.1)
   * [misc-fixed](https://www.cl.cam.ac.uk/~mgk25/ucs-fonts.html) (Public Domain)
   * [M+ Bitmap Font](https://mplus-fonts.osdn.jp/mplus-bitmap-fonts/) (M+ Bitmap Fonts License)
   * Arabic glyphs by [@MansourSorosoro](https://twitter.com/MansourSorosoro) (Eternal Dream Arabization) (OFL-1.1)

There is one font face with glyph size 6x13 for halfwidth, and 12x13 for fullwidth so far.

## Baekmuk License

```
Copyright (c) 1986-2002 Kim Jeong-Hwan
All rights reserved.

Permission to use, copy, modify and distribute this font is
hereby granted, provided that both the copyright notice and
this permission notice appear in all copies of the font,
derivative works or modified versions, and that the following
acknowledgement appear in supporting documentation:
    Baekmuk Batang, Baekmuk Dotum, Baekmuk Gulim, and
    Baekmuk Headline are registered trademarks owned by
    Kim Jeong-Hwan.
```

## Cubic 11 License

```
[Cubic 11]
These fonts are free software.
Unlimited permission is granted to use, copy, and distribute them, with or without modification, either commercially or noncommercially.
THESE FONTS ARE PROVIDED "AS IS" WITHOUT WARRANTY.
此字型是免費的。
無論您是否進行對本字型進行商業或非商業性修改，均可無限制地使用，複製和分發它們。
本字型的衍生品之授權必須與此字型相同，且不作任何擔保。
[JF Dot M+H 12]
Copyright(c) 2005 M+ FONTS PROJECT
[M+ BITMAP FONTS]
Copyright (C) 2002-2004 COZ
These fonts are free software.
Unlimited permission is granted to use, copy, and distribute it, with or without modification, either commercially and noncommercially.
THESE FONTS ARE PROVIDED "AS IS" WITHOUT WARRANTY.
これらのフォントはフリー（自由な）ソフトウエアです。
あらゆる改変の有無に関わらず、また商業的な利用であっても、自由にご利用、複製、再配布することができますが、全て無保証とさせていただきます。

This Font Software is licensed under the SIL Open Font License, Version 1.1.
This license is copied below, and is also available with a FAQ at:
https://scripts.sil.org/OFL
```

## M+ Bitmap Font License

```
-
M+ BITMAP FONTS            Copyright 2002-2005  COZ <coz@users.sourceforge.jp>
-

LICENSE




These fonts are free softwares.
Unlimited permission is granted to use, copy, and distribute it, with
or without modification, either commercially and noncommercially.
THESE FONTS ARE PROVIDED "AS IS" WITHOUT WARRANTY.
```
