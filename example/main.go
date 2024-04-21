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

package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/browser"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/language"

	"github.com/hajimehoshi/bitmapfont/v3"
)

var (
	flagTest     = flag.Bool("test", false, "test mode")
	flagTestSC   = flag.Bool("test-sc", false, "test mode (simplified Chinese)")
	flagTestTC   = flag.Bool("test-tc", false, "test mode (traditional Chinese)")
	flagEastAsia = flag.Bool("eastasia", false, "East Asia")
)

func isTest() bool {
	return *flagTest || *flagTestSC || *flagTestTC
}

func run() error {
	// https://www.ohchr.org/en/human-rights/universal-declaration/universal-declaration-human-rights/about-universal-declaration-human-rights-translation-project
	// https://omniglot.com/udhr/
	text := `en:      All human beings are born free and equal in dignity and rights.
en-Brai: ⠠⠁⠇⠇⠀⠓⠥⠍⠁⠝⠀⠃⠑⠬⠎⠀⠜⠑⠀⠃⠕⠗⠝⠀⠋⠗⠑⠑⠀⠯⠀⠑⠟⠥⠁⠇⠀⠔⠀⠙⠊⠛⠝⠰⠽⠀⠯⠀⠐⠗⠎⠲
ang:     Ealle fīras sind boren frēo ond geefenlican in ār ond riht.
ar:      يولد جميع الناس أحرارًا متساوين في الكرامة والحقوق.
de:      Alle Menschen sind frei und gleich an Würde und Rechten geboren.
el:      'Ολοι οι άνθρωποι γεννιούνται ελεύθεροι και ίσοι στην αξιοπρέπεια και τα δικαιώματα.
es:      Todos los seres humanos nacen libres e iguales en dignidad y derechos.
eo:      Ĉiuj homoj estas denaske liberaj kaj egalaj laŭ digno kaj rajtoj.
fr:      Tous les êtres humains naissent libres et égaux en dignité et en droits.
got:     ᚨᛚᛚᚨᛁ ᛗᚨᚾᚾᚨ ᚠᚱᛖᛁᚺᚨᛚᛋ ᛃᚨᚺ ᛋᚨᛗᚨᛚᛖᛁᚲᛟ ᛁᚾ ᚹᚨᛁᚱᚦᛁᛞᚨᛁ ᛃᚨᚺ ᚱᚨᛁᚺᛏᛖᛁᛋ ᚹᚨᚢᚱᚦᚨᚾᛋ.
he:      כל בני אדם נולדו בני חורין ושווים בערכם ובזכויותיהם.
hy:      Բոլոր մարդիկ ծնվում են ազատ ու հավասար՝ իրենց արժանապատվությամբ և իրավունքներով:
it:      Tutti gli esseri umani nascono liberi ed eguali in dignità e diritti.
ja:      すべての人間は、生れながらにして自由であり、かつ、尊厳と権利とについて平等である。
ka:      ყველა ადამიანი იბადება თავისუფალი და თანასწორი თავისი ღირსებითა და უფლებებით.
ko:      모든 인간은 태어날 때부터 자유로우며 그 존엄과 권리에 있어 동등하다.
mn:      Хүн бүр төрж мэндлэхэд эрх чөлөөтэй, адилхан нэр төртэй, ижил эрхтэй байдаг.
pl:      Wszyscy ludzie rodzą się wolni i równi pod względem swej godności i swych praw.
pt:      Todos os seres humanos nascem livres e iguais em dignidade e em direitos.
ru:      Все люди рождаются свободными и равными в своем достоинстве и правах.
sw:      Watu wote wamezaliwa huru, hadhi na haki zao ni sawa.
tr:      Bütün insanlar hür, haysiyet ve haklar bakımından eşit doğarlar.
uk:      Всі люди народжуються вільними і рівними у своїй гідності та правах.
vi:      Tất cả mọi người sinh ra đều được tự do và bình đẳng về nhân phẩm và quyền.
zh-Hans: 人人生而自由,在尊严和权利上一律平等。他们赋有理性和良心,并应以兄弟关系的精神相对待。
zh-Hant: 人人生而自由，拉尊嚴脫仔權利上一律平等。伊拉有理性脫仔良心，並應以兄弟關係個精神相對待。
`

	if *flagTest {
		text = ""
		for i := 0; i < 256; i++ {
			for j := 0; j < 256; j++ {
				r := rune(i*256 + j)
				if r == '\n' {
					text += " "
					continue
				}
				text += string(r)
			}
			text += "\n"
		}
	}

	path := "example.png"
	if isTest() {
		var suffix string
		if *flagTestSC {
			suffix += "_zh_hans"
		}
		if *flagTestTC {
			suffix += "_zh_hant"
		}
		if *flagEastAsia {
			suffix += "_ea"
		}
		path = "test" + suffix + ".png"
	}
	if err := outputImageFile(text, path); err != nil {
		return err
	}
	return nil
}

func defaultFace() font.Face {
	if *flagTestSC {
		if *flagEastAsia {
			return bitmapfont.FaceSCEA
		}
		return bitmapfont.FaceSC
	}
	if *flagTestTC {
		if *flagEastAsia {
			return bitmapfont.FaceTCEA
		}
		return bitmapfont.FaceTC
	}
	if *flagEastAsia {
		return bitmapfont.FaceEA
	}
	return bitmapfont.Face
}

func outputImageFile(text string, path string) error {
	const (
		offsetX = 8
		offsetY = 8
	)

	const (
		glyphWidth  = 12
		glyphHeight = 16
	)

	lines := strings.Split(strings.TrimSpace(text), "\n")
	width := 0
	for _, l := range lines {
		w := int(font.MeasureString(defaultFace(), l).Ceil())
		if width < w {
			width = w
		}
	}

	width += offsetX * 2
	height := glyphHeight*len(lines) + offsetY*2

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	if isTest() {
		gray := color.RGBA{0xcc, 0xcc, 0xcc, 0xff}
		for j := 0; j < 256; j++ {
			for i := 0; i < 256; i++ {
				if (i+j)%2 == 0 {
					continue
				}
				x := i*glyphWidth + offsetX
				y := j*glyphHeight + offsetY
				draw.Draw(dst, image.Rect(x, y, x+glyphWidth, y+glyphHeight), image.NewUniform(gray), image.Point{}, draw.Src)
			}
		}
	}

	d := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(color.Black),
		Face: defaultFace(),
		Dot:  fixed.Point26_6{X: fixed.I(offsetX), Y: defaultFace().Metrics().Ascent + fixed.I(offsetY)},
	}

	langRe := regexp.MustCompile(`^[a-zA-Z0-9-]+`)

	for _, l := range strings.Split(text, "\n") {
		langstr := langRe.FindString(l)
		if !isTest() && langstr != "" {
			lang, err := language.Parse(langstr)
			if err != nil {
				return err
			}
			l = bitmapfont.PresentationForms(l, bitmapfont.DirectionLeftToRight, lang)
		}
		f := defaultFace()
		if !isTest() {
			if strings.HasPrefix(langstr, "zh-Hans") {
				if *flagEastAsia {
					f = bitmapfont.FaceSCEA
				} else {
					f = bitmapfont.FaceSC
				}
			}
			if strings.HasPrefix(langstr, "zh-Hant") {
				if *flagEastAsia {
					f = bitmapfont.FaceTCEA
				} else {
					f = bitmapfont.FaceTC
				}
			}
		}
		d.Face = f
		d.Dot.X = fixed.I(offsetX)
		d.DrawString(l)
		d.Dot.Y += f.Metrics().Height
	}

	fout, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fout.Close()

	if err := png.Encode(fout, d.Dst); err != nil {
		return err
	}

	if err := browser.OpenFile(path); err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		panic(err)
	}
}
